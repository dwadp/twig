package twig

import "github.com/spf13/cobra"

var configReloadCmd = &cobra.Command{
	Use: "reload",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.PrintErrln("Not implemented yet")
	},
}
