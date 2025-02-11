package edit

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	fleetdbapi "github.com/metal-automata/fleetdb/pkg/api/v1"
	"github.com/spf13/cobra"

	mctl "github.com/metal-automata/mctl/cmd"
	"github.com/metal-automata/mctl/internal/app"
	"github.com/metal-automata/mctl/pkg/model"
)

var (
	editFWSetFlags mctl.FirmwareSetFlags
)

var editFirmwareSet = &cobra.Command{
	Use:   "firmware-set",
	Short: "Edit a firmware set",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(editFWSetFlags.ID)
		if err != nil {
			log.Fatalf("invalid set-id: %s, error: %s", editFWSetFlags.ID, err.Error())
		}

		payload := fleetdbapi.ComponentFirmwareSetRequest{
			ID:                     id,
			ComponentFirmwareUUIDs: []string{},
		}

		var attrs *fleetdbapi.Attributes
		var payloadUpdated bool

		if len(editFWSetFlags.Labels) > 0 {
			attrs, err = mctl.AttributeFromLabels(model.AttributeNSFirmwareSetLabels, editFWSetFlags.Labels)
			if err != nil {
				log.Fatal(err)
			}

			payload.Attributes = []fleetdbapi.Attributes{*attrs}
			payloadUpdated = true

		}

		if len(editFWSetFlags.AddFirmwareUUIDs) > 0 {
			for _, id := range editFWSetFlags.AddFirmwareUUIDs {
				_, err = uuid.Parse(id)
				if err != nil {
					log.Fatal(err)
				}

				payload.ComponentFirmwareUUIDs = append(payload.ComponentFirmwareUUIDs, id)
				payloadUpdated = true
			}
		}

		if editFWSetFlags.FirmwareSetName != "" {
			payload.Name = editFWSetFlags.FirmwareSetName
			payloadUpdated = true
		}

		if payloadUpdated {
			_, err = client.UpdateComponentFirmwareSetRequest(cmd.Context(), id, payload)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("firmware set updated: " + id.String())
		}

		if len(editFWSetFlags.RemoveFirmwareUUIDs) > 0 {
			for _, id := range editFWSetFlags.RemoveFirmwareUUIDs {
				_, err = uuid.Parse(id)
				if err != nil {
					log.Fatal(err)
				}

				payload.ComponentFirmwareUUIDs = append(payload.ComponentFirmwareUUIDs, id)
			}

			_, err = client.RemoveServerComponentFirmwareSetFirmware(cmd.Context(), id, payload)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("firmware set uuids removed: " + id.String())
		}
	},
}

func init() {
	mctl.AddFirmwareSetFlag(editFirmwareSet, &editFWSetFlags.ID)
	mctl.AddNameFlag(editFirmwareSet, &editFWSetFlags.FirmwareSetName, "New name of the firmware set")
	mctl.AddLabelsFlag(editFirmwareSet, &editFWSetFlags.Labels,
		"Labels to assign to the firmware set - 'vendor=foo,model=bar'")
	mctl.AddFirmwareAddIDsFlag(editFirmwareSet, &editFWSetFlags.AddFirmwareUUIDs)
	mctl.AddFirmwareRemoveIDsFlag(editFirmwareSet, &editFWSetFlags.RemoveFirmwareUUIDs)

	mctl.RequireFlag(editFirmwareSet, mctl.FirmwareSetFlag)
}
