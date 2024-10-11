package bios

import (
	"log"

	mctl "github.com/metal-automata/mctl/cmd"
	rctypes "github.com/metal-automata/rivets/condition"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set BIOS settings from github config file url",
	Run: func(cmd *cobra.Command, _ []string) {
		err := CreateBiosControlCondition(cmd.Context(), rctypes.SetConfig)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	mctl.AddServerFlag(setCmd, &biosFlags.serverID)
	mctl.AddBIOSConfigURLFlag(setCmd, &biosFlags.biosConfigURL)

	mctl.RequireFlag(setCmd, mctl.ServerFlag)
	mctl.RequireFlag(setCmd, mctl.BIOSConfigURLFlag)

	biosCmd.AddCommand(setCmd)
}
