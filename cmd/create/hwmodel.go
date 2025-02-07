package create

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	fleetdbapi "github.com/metal-automata/fleetdb/pkg/api/v1"
	mctl "github.com/metal-automata/mctl/cmd"
	"github.com/metal-automata/mctl/internal/app"
)

type hwModelCreateParams struct {
	vendor string
	model  string
}

var (
	hwModelCreateFlags *hwModelCreateParams
)

var hwModelCreate = &cobra.Command{
	Use:   "hardware-model",
	Short: "Create hardware model",
	Run: func(cmd *cobra.Command, _ []string) {
		createHwModel(cmd.Context())
	},
}

func createHwModel(ctx context.Context) {
	hwmodel := fleetdbapi.HardwareModel{
		Name:               hwModelCreateFlags.model,
		HardwareVendorName: hwModelCreateFlags.vendor,
	}

	appObj := mctl.MustCreateApp(ctx)

	client, err := app.NewFleetDBAPIClient(ctx, appObj.Config.FleetDBAPI, appObj.Reauth)
	if err != nil {
		log.Fatal(err)

	}

	resp, errCreate := client.CreateHardwareModel(ctx, &hwmodel)
	if errCreate != nil {
		log.Fatal(errCreate)
	}

	log.Printf("hardware model created: %s", resp.Slug)
}

func init() {
	hwModelCreateFlags = &hwModelCreateParams{}

	mctl.AddHardwareVendorFlag(hwModelCreate, &hwModelCreateFlags.vendor)
	mctl.AddHardwareModelFlag(hwModelCreate, &hwModelCreateFlags.model)

	mctl.RequireFlag(hwModelCreate, mctl.HardwareVendorFlag)
	mctl.RequireFlag(hwModelCreate, mctl.HardwareModelFlag)
}
