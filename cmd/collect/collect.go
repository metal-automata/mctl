package collect

import (
	"github.com/spf13/cobra"

	"github.com/metal-automata/mctl/cmd"
)

var collect = &cobra.Command{
	Use:   "collect",
	Short: "Collect current server firmware status and bios configuration",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(collect)
	collect.AddCommand(collectInventoryCmd)
	collect.AddCommand(inventoryStatus)
}
