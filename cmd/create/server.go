package create

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	fleetdbapi "github.com/metal-automata/fleetdb/pkg/api/v1"
	mctl "github.com/metal-automata/mctl/cmd"
	"github.com/metal-automata/mctl/internal/app"
)

type serverCreateParams struct {
	serverID    string
	facility    string
	bmcIP       string
	bmcUsername string
	bmcPassword string
	bmcMac      string
	vendor      string
	model       string
	fromFile    string
}

var (
	serverCreateFlags *serverCreateParams
)

var serverCreate = &cobra.Command{
	Use:   "server",
	Short: "Create server record with BMC information",
	Run: func(cmd *cobra.Command, _ []string) {
		if serverCreateFlags.fromFile == "" {
			mctl.RequireFlag(cmd, mctl.BMCAddressFlag)
			mctl.RequireFlag(cmd, mctl.BMCUsernameFlag)
			mctl.RequireFlag(cmd, mctl.BMCPasswordFlag)
			mctl.RequireFlag(cmd, mctl.FacilityFlag)
			mctl.RequireFlag(cmd, mctl.HardwareVendorFlag)
		}

		createServer(cmd.Context())
	},
}

func createServer(ctx context.Context) {
	var servers []*fleetdbapi.Server

	// load server data from file
	if serverCreateFlags.fromFile != "" {
		sbytes, err := os.ReadFile(serverCreateFlags.fromFile)
		if err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(sbytes, &servers); err != nil {
			log.Fatal(err)
		}
	} else {

		// load server data from CLI input
		serverUUID := uuid.MustParse(serverCreateFlags.serverID)
		servers = append(servers, &fleetdbapi.Server{
			UUID:         serverUUID,
			FacilityCode: serverCreateFlags.facility,
			Vendor:       serverCreateFlags.vendor,
			Model:        serverCreateFlags.model,
			BMC: &fleetdbapi.ServerBMC{
				ServerID:           serverUUID,
				IPAddress:          serverCreateFlags.bmcIP,
				Username:           serverCreateFlags.bmcUsername,
				Password:           serverCreateFlags.bmcPassword,
				MacAddress:         serverCreateFlags.bmcMac,
				HardwareVendorName: serverCreateFlags.vendor,
				HardwareModelName:  serverCreateFlags.model,
			},
		})
	}

	createServerRecords(ctx, servers)
}

func createServerRecords(ctx context.Context, servers []*fleetdbapi.Server) {
	appObj := mctl.MustCreateApp(ctx)

	client, err := app.NewFleetDBAPIClient(ctx, appObj.Config.FleetDBAPI, appObj.Reauth)
	if err != nil {
		log.Fatal(err)

	}

	for idx := range servers {
		// The serverID is a prerequisite for BMCs to be added in the same tx
		// once BMCs can be added without a server identifier
		// the client does not have to generate the serverID
		if servers[idx].UUID == uuid.Nil {
			servers[idx].UUID = uuid.New()
			servers[idx].BMC.ServerID = servers[idx].UUID
		}

		if servers[idx].BMC.MacAddress == "" {
			servers[idx].BMC.MacAddress = "00:00:00:00:00:00"
		} else {
			_, err := net.ParseMAC(servers[idx].BMC.MacAddress)
			if err != nil {
				log.Fatal(err)
			}
		}

		id, _, err := client.Create(ctx, *servers[idx])
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("server created: %s", id)
	}

}

func init() {
	serverCreateFlags = &serverCreateParams{}

	mctl.AddFacilityFlag(serverCreate, &serverCreateFlags.facility)
	mctl.AddServerFlag(serverCreate, &serverCreateFlags.serverID)
	mctl.AddBMCAddressFlag(serverCreate, &serverCreateFlags.bmcIP)
	mctl.AddBMCMacAddressFlag(serverCreate, &serverCreateFlags.bmcMac)
	mctl.AddBMCUsernameFlag(serverCreate, &serverCreateFlags.bmcUsername)
	mctl.AddBMCPasswordFlag(serverCreate, &serverCreateFlags.bmcPassword)
	mctl.AddHardwareVendorFlag(serverCreate, &serverCreateFlags.vendor)
	mctl.AddHardwareModelFlag(serverCreate, &serverCreateFlags.model)
	mctl.AddFromFileFlag(serverCreate, &serverCreateFlags.fromFile, "JSON file with server data")

}
