package cmd

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/dwadp/twig/setup"
)

var (
	//go:embed stub/config.stub
	stub string
)

func ExecSetup(cmd *flag.FlagSet) int {
	if err := setup.Setup(os.Stdout, stub); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run setup command: %v\n", err)
		return 1
	}
	return 0
}
