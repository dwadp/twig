package twig

import (
    "github.com/dwadp/twig/pkg/config"
    "github.com/dwadp/twig/pkg/schema"
    "io"
    "os"
    "os/exec"
)

var execCommand = exec.Command

func runPhpExec(container string) ([]byte, error) {
	cmd := execCommand("php", "-v")
	return cmd.CombinedOutput()
}

func RunPHP(cfg *config.Config, args []string) error {
	file, err := os.Open("./composer.json")
	if err != nil {
		return err
	}

	b, err := io.ReadAll(file)

	versions, err := schema.GetVersion(b)
	if err != nil {
		return err
	}

	php, err := cfg.GetPreferredPHPVersion(versions)
	if err != nil {
		return err
	}

	cmd := exec.Command(php.ExecutablePath, args...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func RunComposer() {}
