package install

import (
	"github.com/spf13/cobra"

	"github.com/metal-automata/mctl/cmd"
)

var install = &cobra.Command{
	Use:   "install",
	Short: "Install actions",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(install)

	install.AddCommand(installFirmwareSet)
	install.AddCommand(installStatus)
}
