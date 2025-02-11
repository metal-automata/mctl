package get

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	mctl "github.com/metal-automata/mctl/cmd"
	"github.com/metal-automata/mctl/internal/app"
)

type getFirmwareFlags struct {
	id string
}

var (
	flagsDefinedGetFirmware *getFirmwareFlags
)

// Get firmware info
var getFirmware = &cobra.Command{
	Use:   "firmware",
	Short: "Get information for given firmware identifier",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		client, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		fwID, err := uuid.Parse(flagsDefinedGetFirmware.id)
		if err != nil {
			log.Fatal(err)
		}

		firmware, _, err := client.GetServerComponentFirmware(ctx, fwID)
		if err != nil {
			log.Fatal("fleetdb API client returned error: ", err)
		}

		mctl.PrintResults(output, firmware)
		os.Exit(0)
	},
}

func init() {
	flagsDefinedGetFirmware = &getFirmwareFlags{}

	mctl.AddFirmwareIDFlag(getFirmware, &flagsDefinedGetFirmware.id)
	mctl.RequireFlag(getFirmware, mctl.FirmwareIDFlag)
}
