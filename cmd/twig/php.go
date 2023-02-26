package twig

import (
	"github.com/dwadp/twig/pkg/twig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(phpCmd)
}

var phpCmd = &cobra.Command{
	Use:    "php",
	Short:  "Run a PHP command",
	Args:   cobra.ArbitraryArgs,
	PreRun: preRun,
	Run: func(cmd *cobra.Command, args []string) {
		if err := twig.RunPHP(cfg, args); err != nil {
			cmd.PrintErrln(err)
		}
	},
}
