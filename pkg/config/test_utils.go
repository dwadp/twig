package config

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

func CleanupTestFiles(name string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(path.Join(home, "."+name)); err != nil {
		return err
	}

	return nil
}
