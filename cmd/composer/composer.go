package composer

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

	return run(cfg, args)
}

func run(cfg *config.Config, args []string) int {
	version, err := proxy.GetVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	php, err := cfg.GetPHP(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not find configuration for PHP version %s; %v\n", version, err)
		return 1
	}

	composer := cfg.Composer
	options := []string{composer.Path}
	options = append(options, args...)

	cmd := exec.Command(php.Path, options...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run composer %v\n", err)
		return 1
	}

	return 0
}
