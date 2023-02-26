package twig

import (
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configSetDefaultCmd, configInitCmd, configReloadCmd)

	configCmd.Flags().BoolP("locate", "l", true, "Print where the configuration file is located.")
}

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Manage the twig configuration",
	Aliases: []string{"c"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		locate, err := cmd.Flags().GetBool("locate")
		if err != nil {
			cmd.PrintErrf("Failed to parse the \"locate\" flag: %v\n", err)
		}
		if locate {
			cmd.Printf("The twig configuration file is located in %q\n", cfg.FilePath())
		}
	},
}
