package twig

import (
	"github.com/spf13/cobra"
)

var configSetDefaultCmd = &cobra.Command{
	Use: "set-default",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.PrintErrln("Not implemented yet")
	},
}
