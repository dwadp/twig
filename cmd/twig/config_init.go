package twig

import (
	"os"

	"github.com/spf13/cobra"
)

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the configuration file",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cfg.Init(); err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		cmd.Printf("The configuration file successfully created on %q\n", cfg.FilePath())
	},
}
