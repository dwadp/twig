package php

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dwadp/twig/config"
	"github.com/dwadp/twig/proxy"
)

func Exec(args []string) int {
	cfg, err := config.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	version, err := proxy.GetVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if len(os.Args) > 1 {
		if ver := os.Args[1]; cfg.IsVersionAvailable(ver) {
			version = ver
			args = os.Args[2:]
		}
	}

	c, err := cfg.GetPHP(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not find configuration for PHP version %s; %v\n", version, err)
		return 1
	}

	return run(c, args)
}

func run(cfg *config.Item, args []string) int {
	cmd := exec.Command(cfg.Path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run php %v\n", err)
		return 1
	}

	return 0
}
