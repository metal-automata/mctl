package deleteresource

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	fleetdbapi "github.com/metal-automata/fleetdb/pkg/api/v1"
	"github.com/spf13/cobra"

	mctl "github.com/metal-automata/mctl/cmd"
	"github.com/metal-automata/mctl/internal/app"
)

type deleteFirmwareFlags struct {
	// firmware UUID
	id string
}

var (
	flagsDefinedDeleteFirmware *deleteFirmwareFlags
)
var deleteFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "Delete a firmware object",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(flagsDefinedDeleteFirmware.id)
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.DeleteServerComponentFirmware(cmd.Context(), fleetdbapi.ComponentFirmwareVersion{UUID: id})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("firmware deleted: " + id.String())
	},
}

func init() {
	flagsDefinedDeleteFirmware = &deleteFirmwareFlags{}

	mctl.AddFirmwareIDFlag(deleteFirmware, &flagsDefinedDeleteFirmware.id)
	mctl.RequireFlag(deleteFirmware, mctl.FirmwareIDFlag)
}
