package get

import (
	"github.com/spf13/cobra"

	"github.com/metal-automata/mctl/cmd"
)

var (
	output string
)

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get resource",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(cmdGet)
	cmdGet.AddCommand(getServer)
	cmdGet.AddCommand(getCondition)
	cmdGet.AddCommand(getFirmware)
	cmdGet.AddCommand(getFirmwareSet)
	cmdGet.AddCommand(getBiosConfig)
	cmdGet.AddCommand(getBomInfoByMacAddress)

	cmd.AddOutputFlag(cmdGet, &output)
}
