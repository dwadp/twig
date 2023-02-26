package twig

import "github.com/spf13/cobra"

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the configuration file",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cfg.Init(); err != nil {
			cmd.PrintErrf("Error while initializing configuration: %v\n", err)
		}
		cmd.Printf("The configuration file successfully created on %q\n", cfg.FilePath())
	},
}
