package main

import (
	"flag"
	"os"

	"github.com/dwadp/twig/cmd"
	"github.com/dwadp/twig/cmd/composer"
	"github.com/dwadp/twig/cmd/php"
)

func main() {
	setupCmd := flag.NewFlagSet("setup", flag.ExitOnError)
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)

	if len(os.Args) > 1 {
		args := []string{}
		if len(os.Args) > 2 {
			args = os.Args[2:]
		}

		switch os.Args[1] {
		case "php":
			os.Exit(php.Exec(args))
		case "composer":
			os.Exit(composer.Exec(args))
		case "setup":
			setupCmd.Parse(args)
			os.Exit(cmd.ExecSetup(setupCmd))
		case "set":
			setCmd.Parse(args)
			os.Exit(cmd.ExecSet(setCmd))
		}
	}

	os.Exit(0)
}
