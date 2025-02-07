package create

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	fleetdbapi "github.com/metal-automata/fleetdb/pkg/api/v1"
	mctl "github.com/metal-automata/mctl/cmd"
	"github.com/metal-automata/mctl/internal/app"
)

type hwVendorCreateParams struct {
	vendor string
}

var (
	hwVendorCreateFlags *hwVendorCreateParams
)

var hwVendorCreate = &cobra.Command{
	Use:   "hardware-vendor",
	Short: "Create hardware vendor",
	Run: func(cmd *cobra.Command, _ []string) {
		createHwVendor(cmd.Context())
	},
}

func createHwVendor(ctx context.Context) {
	hwVendor := fleetdbapi.HardwareVendor{
		Name: hwVendorCreateFlags.vendor,
	}

	appObj := mctl.MustCreateApp(ctx)

	client, err := app.NewFleetDBAPIClient(ctx, appObj.Config.FleetDBAPI, appObj.Reauth)
	if err != nil {
		log.Fatal(err)

	}

	resp, errCreate := client.CreateHardwareVendor(ctx, &hwVendor)
	if errCreate != nil {
		log.Fatal(errCreate)
	}

	log.Printf("hardware model created: %s", resp.Slug)
}

func init() {
	hwVendorCreateFlags = &hwVendorCreateParams{}
	mctl.AddHardwareVendorFlag(hwVendorCreate, &hwVendorCreateFlags.vendor)
	mctl.RequireFlag(hwVendorCreate, mctl.HardwareVendorFlag)
}
