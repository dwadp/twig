package twig

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dwadp/twig/pkg/config"
	"github.com/dwadp/twig/pkg/schema"
)

var execCommand = exec.Command

func runPhpExec() ([]byte, error) {
	cmd := execCommand("php", "-v")
	return cmd.CombinedOutput()
}

func isRootProjectDir(cwd string) bool {
	_, err := os.Stat(filepath.Join(cwd, "composer.json"))

	if err != nil && errors.Is(err, os.ErrNotExist) {
		vendorDir, err := os.Stat(filepath.Join(cwd, "vendor"))

		if err != nil && errors.Is(err, os.ErrNotExist) {
			return false
		}

		return vendorDir.IsDir()
	}

	return err == nil
}

func preparePHP(cfg *config.Config) (*config.PHP, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if isRootProjectDir(wd) {
		file, err := os.Open(filepath.Join(wd, "composer.json"))
		if err != nil {
			return nil, err
		}

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		versions, err := schema.GetVersion(b)
		if err != nil {
			return nil, err
		}

		php, err := cfg.GetPreferredPHPVersion(versions)
		if err != nil {
			return nil, err
		}

		return php, nil
	}

	return cfg.GetDefaultPHPVersion()
}

func RunPHP(cfg *config.Config, args []string) error {
	php, err := preparePHP(cfg)
	if err != nil {
		return err
	}

	cmd := exec.Command(php.ExecutablePath, args...)

	if err := runInTTY(cmd); err != nil {
		return err
	}

	return nil
}

func RunComposer(cfg *config.Config, args []string) error {
	php, err := preparePHP(cfg)
	if err != nil {
		return err
	}

	args = append([]string{cfg.Composer.ExecutablePath}, args...)

	cmd := exec.Command(php.ExecutablePath, args...)

	if err := runInTTY(cmd); err != nil {
		return err
	}

	return nil
}
