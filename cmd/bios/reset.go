package bios

import (
	"log"

	mctl "github.com/metal-automata/mctl/cmd"
	rctypes "github.com/metal-automata/rivets/condition"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset BIOS settings to default values",
	Run: func(cmd *cobra.Command, _ []string) {
		err := CreateBiosControlCondition(cmd.Context(), rctypes.ResetConfig)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	mctl.AddServerFlag(resetCmd, &biosFlags.serverID)

	mctl.RequireFlag(resetCmd, mctl.ServerFlag)

	biosCmd.AddCommand(resetCmd)
}
