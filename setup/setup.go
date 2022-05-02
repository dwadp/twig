package setup

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"
)

var (
	out io.Writer
)

// Setup will create configuration file inside your home directory (.twig/config.yml)
// after the configuration file created, you should fill it with your php installation paths
func Setup(o io.Writer, stub string) error {
	out = o
	basepath, err := getBasepath()

	if err != nil {
		return err
	}

	filename := path.Join(basepath, "config.yml")

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return createCfgFile(basepath, filename, stub)
		}
		return err
	}

	fmt.Fprintln(out, "Configuration file already exist!")

	return nil
}

func createCfgFile(basepath, filename, stub string) error {
	if err := os.Mkdir(basepath, os.ModeDir); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("error creating directory at [%s]; %v\n", basepath, err)
		}
	}

	file, err := os.Create(filename)

	defer file.Close()
	if err != nil {
		return fmt.Errorf("could not create configuration file [%s]; %v\n", filename, err)
	}

	if _, err := file.WriteString(stub); err != nil {
		return fmt.Errorf("could not create configuration file [%s]; %v\n", filename, err)
	}

	fmt.Fprintf(out, "Configuration file successfully created in [%s]. You're ready to go!\n", filename)

	return nil
}

func getBasepath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get base path; %v\n", err)
	}
	return path.Join(home, ".twig"), nil
}
