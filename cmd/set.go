package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/dwadp/twig/config"
)

func ExecSet(cmd *flag.FlagSet) int {
	version := cmd.Arg(0)

	cfg, err := config.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if err := cfg.SetDefault(version); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	fmt.Fprintf(os.Stdout, "Default PHP version successfully set to %s\n", version)

	return 0
}
