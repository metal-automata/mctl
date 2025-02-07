package create

import (
	"github.com/metal-automata/mctl/cmd"

	"github.com/spf13/cobra"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "Create resources",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(create)
	create.AddCommand(hwVendorCreate)
	create.AddCommand(hwModelCreate)
	create.AddCommand(createFirmware)
	create.AddCommand(createFirmwareSet)
	create.AddCommand(serverCreate)
}
