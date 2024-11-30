package twig

import (
	"os"

	"github.com/dwadp/twig/pkg/twig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(composerCmd)
}

// composerCmd represents the composer command
var composerCmd = &cobra.Command{
	Use:                "composer",
	Short:              "Run a Composer command",
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true, // Accepts any given command line flags and let PHP handle them
	PreRun:             preRun,
	Run: func(cmd *cobra.Command, args []string) {
		if err := twig.RunComposer(cfg, args); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}
