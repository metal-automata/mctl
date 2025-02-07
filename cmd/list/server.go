package list

import (
	"log"
	"os"

	fleetdbapi "github.com/metal-automata/fleetdb/pkg/api/v1"
	mctl "github.com/metal-automata/mctl/cmd"
	"github.com/metal-automata/mctl/internal/app"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type listServerFlags struct {
	records   bool
	bmcerrors bool
	creds     bool
	table     bool
	vendor    string
	model     string
	serial    string
	facility  string
	limit     int
	page      int
}

var (
	flagsListServer *listServerFlags
)

// List
var cmdListServer = &cobra.Command{
	Use:   "server",
	Short: "List servers",
	Run: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()
		theApp := mctl.MustCreateApp(ctx)

		if flagsListServer.limit > fleetdbapi.MaxPaginationSize {
			log.Printf("Notice: Limit was set above max, setting limit to %d. If you want to list more than %d servers, please use '--page` to index individual pages", fleetdbapi.MaxPaginationSize, fleetdbapi.MaxPaginationSize)
			flagsListServer.limit = fleetdbapi.MaxPaginationSize
		}

		client, err := app.NewFleetDBAPIClient(ctx, theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		lsp := &fleetdbapi.ServerQueryParams{
			FilterParams: &fleetdbapi.FilterParams{
				Target: &fleetdbapi.Server{},
				Filters: []fleetdbapi.Filter{
					{
						Attribute:          "facility_code",
						ComparisonOperator: fleetdbapi.ComparisonOpEqual,
						Value:              flagsListServer.facility,
					},
				},
			},
			PaginationParams: &fleetdbapi.
				PaginationParams{
				Limit:   flagsListServer.limit,
				Page:    flagsListServer.page,
				Preload: false,
				OrderBy: "",
			},
		}

		servers, res, err := client.ListServers(ctx, lsp)
		if err != nil {
			log.Fatal(err)
		}

		if flagsListServer.records {
			d := struct {
				CurrentPage      int
				Limit            int
				TotalPages       int
				TotalRecordCount int64
				Link             string
			}{
				res.Page,
				res.PageCount,
				res.TotalPages,
				res.TotalRecordCount,
				res.Links.Self.Href,
			}

			printJSON(d)

			os.Exit(0)
		}

		if len(servers) == 0 {
			log.Println("no servers matched filters")
			os.Exit(0)
		}

		printJSON(servers)
	},
}

func serversTable(servers []*fleetdbapi.Server, fl *listServerFlags) {
	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"UUID", "Name", "Vendor", "Model", "Serial", "BMCAddr"}

	if fl.creds {
		headers = append(headers, []string{"BMCUser", "BMCPass"}...)
	}

	table.SetHeader(headers)
	for _, server := range servers {
		row := []string{
			server.UUID.String(),
			server.Name,
			server.Vendor,
			server.Model,
			server.Serial,
			server.BMC.IPAddress,
		}

		if fl.creds {
			row = append(row, []string{server.BMC.Username, server.BMC.Password}...)
		}

		table.Append(row)
	}

	table.Render()
}

func init() {
	flagsListServer = &listServerFlags{}

	mctl.AddWithRecordsFlag(cmdListServer, &flagsListServer.records)
	mctl.AddVendorFlag(cmdListServer, &flagsListServer.vendor)
	mctl.AddModelFlag(cmdListServer, &flagsListServer.model)
	mctl.AddFacilityFlag(cmdListServer, &flagsListServer.facility)
	mctl.AddPageFlag(cmdListServer, &flagsListServer.page)
	mctl.AddPageLimitFlag(cmdListServer, &flagsListServer.limit)
	mctl.AddWithCredsFlag(cmdListServer, &flagsListServer.creds)
	mctl.AddPrintTableFlag(cmdListServer, &flagsListServer.table)
	mctl.AddServerSerialFlag(cmdListServer, &flagsListServer.serial)
}
