package get

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	fleetdbapi "github.com/metal-automata/fleetdb/pkg/api/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	mctl "github.com/metal-automata/mctl/cmd"
	"github.com/metal-automata/mctl/internal/app"
)

type getServerFlags struct {
	// server UUID
	id             string
	component      string
	listComponents bool
	biosconfig     bool
	table          bool
	creds          bool
}

var (
	cmdArgs    *getServerFlags
	cmdTimeout = 2 * time.Minute
)

var getServer = &cobra.Command{
	Use:   "server",
	Short: "Get server information",
	Run: func(cmd *cobra.Command, _ []string) {
		ctx, cancel := context.WithTimeout(cmd.Context(), cmdTimeout)
		defer cancel()

		theApp := mctl.MustCreateApp(ctx)

		client, err := app.NewFleetDBAPIClient(ctx, theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.Parse(cmdArgs.id)
		if err != nil {
			log.Fatal(err)
		}

		withComponents := cmdArgs.listComponents || cmdArgs.component != ""
		params := fleetdbapi.ServerQueryParams{
			IncludeComponents: withComponents,
			IncludeBMC:        cmdArgs.creds,
		}

		if withComponents {
			params.ComponentParams = &fleetdbapi.ServerComponentGetParams{
				InstalledFirmware: true,
			}
		}

		server, _, err := client.GetServer(ctx, id, &params)
		if err != nil {
			log.Fatal(err)
		}

		if cmdArgs.table {
			switch {
			case cmdArgs.listComponents:
				renderComponentListTable(server.Components)
			default:
				renderServerTable(server, cmdArgs.creds)
			}

			os.Exit(0)
		}

		if cmdArgs.component != "" {
			printComponent(server.Components, cmdArgs.component)
			os.Exit(0)
		}

		if cmdArgs.biosconfig {
			fmt.Println("server bios config listing to be implemented")
			os.Exit(0)
		}

		mctl.PrintResults(output, server)
	},
}

func printComponent(components []*fleetdbapi.ServerComponent, slug string) {
	got := []*fleetdbapi.ServerComponent{}

	for _, c := range components {
		c := c
		if strings.EqualFold(slug, c.Name) {
			got = append(got, c)
		}
	}

	mctl.PrintResults(output, got)
}

func renderServerTable(server *fleetdbapi.Server, withCreds bool) {
	tableServer := tablewriter.NewWriter(os.Stdout)
	tableServer.Append([]string{"ID", server.UUID.String()})
	tableServer.Append([]string{"Name", server.Name})
	tableServer.Append([]string{"Model", server.Model})
	tableServer.Append([]string{"Vendor", server.Vendor})
	tableServer.Append([]string{"Serial", server.Serial})
	tableServer.Append([]string{"BMCAddr", server.BMC.IPAddress})
	if withCreds {
		tableServer.Append([]string{"BMCUser", server.BMC.Username})
		tableServer.Append([]string{"BMCPassword", server.BMC.Password})
	}
	tableServer.Append([]string{"Facility", server.FacilityCode})
	tableServer.Append([]string{"Reported", humanize.Time(server.UpdatedAt)})

	tableServer.Render()
}

func renderComponentListTable(components []*fleetdbapi.ServerComponent) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Component", "Vendor", "Model", "Serial", "FW", "Status", "Reported"})
	for _, c := range components {
		vendor := "-"
		model := "-"
		serial := "-"
		installed := "-"
		status := "-"

		if c.InstalledFirmware != nil && c.InstalledFirmware.Version != "" {
			installed = c.InstalledFirmware.Version
		}

		if c.Status != nil {
			if c.Status.Health != "" {
				status = c.Status.Health
			} else if c.Status.State != "" {
				status = c.Status.State
			}
		}

		if c.Vendor != "" {
			vendor = c.Vendor
		}

		if c.Model != "" {
			model = c.Model
		}

		if c.Serial != "" {
			serial = c.Serial
		}

		table.Append([]string{c.Name, vendor, model, serial, installed, status, humanize.Time(c.UpdatedAt)})
	}

	table.Render()
}

func init() {
	cmdArgs = &getServerFlags{}

	mctl.AddServerFlag(getServer, &cmdArgs.id)
	mctl.AddSlugFlag(getServer, &cmdArgs.component, "list component on server by slug (drive/nic/cpu..)")
	mctl.AddWithCredsFlag(getServer, &cmdArgs.creds)
	mctl.AddPrintTableFlag(getServer, &cmdArgs.table)
	mctl.AddBIOSConfigFlag(getServer, &cmdArgs.biosconfig)
	mctl.AddListComponentsFlag(getServer, &cmdArgs.listComponents)

	mctl.RequireFlag(getServer, mctl.ServerFlag)
}
