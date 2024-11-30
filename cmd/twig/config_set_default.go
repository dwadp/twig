package twig

import (
	"os"

	"github.com/spf13/cobra"
)

var configSetDefaultCmd = &cobra.Command{
	Use:  "set-default",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cfg.SetDefault(args[0]); err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		cmd.Printf("default PHP version set to %s\n", args[0])
	},
}
