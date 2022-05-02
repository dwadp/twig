package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoDefaultVersion  = errors.New("no default version was set")
	ErrInvalidPHPVersion = errors.New("invalid PHP version")
)

type Config struct {
	versions []string
	Php      map[string]Item `yaml:"php"`
	Composer Item            `yaml:"composer"`
}

type Item struct {
	Version string `yaml:"-"`
	Path    string `yaml:"path"`
	Default bool   `yaml:"default,omitempty"`
}

func New() (*Config, error) {
	cfg := Config{}
	if err := cfg.read(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c Config) GetPHP(version string) (*Item, error) {
	item, ok := c.Php[version]
	if !ok {
		return c.GetDefault()
	}
	return &item, nil
}

func (c Config) GetDefault() (*Item, error) {
	for k, item := range c.Php {
		if item.Default {
			item.Version = k
			return &item, nil
		}
	}
	return nil, ErrNoDefaultVersion
}

func (c *Config) SetDefault(version string) error {
	def, err := c.GetDefault()
	if err != nil {
		if !errors.Is(err, ErrNoDefaultVersion) {
			return fmt.Errorf("could not retrieve default PHP version %v\n", err)
		}
	}

	if !c.IsVersionAvailable(version) {
		return ErrInvalidPHPVersion
	}

	if def != nil {
		if def.Version == version {
			return nil
		}
	}

	items := []Item{}

	for _, item := range c.PHP() {
		if item.Version == version {
			item.Default = true
		} else if def != nil {
			if item.Version == def.Version {
				item.Default = false
			}
		}
		items = append(items, item)
	}

	c.SetPHP(items)

	return c.write()
}

func (c *Config) SetPHP(items []Item) {
	php := make(map[string]Item)
	for _, item := range items {
		php[item.Version] = item
	}
	c.Php = php
}

func (c Config) IsVersionAvailable(version string) bool {
	_, ok := c.Php[version]
	return ok
}

func (c Config) PHP() []Item {
	items := []Item{}
	for version, item := range c.Php {
		item.Version = version
		items = append(items, item)
	}
	return items
}

func (c *Config) write() error {
	conf, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("could not parse configuration file %v\n", err)
	}

	filename, err := c.filename()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, conf, 0)
}

func (c Config) filename() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find user home directory [%v]\n", err)
	}
	return path.Join(home, ".twig", "config.yml"), nil
}

func (c *Config) read() error {
	filename, err := c.filename()
	if err != nil {
		return err
	}

	file, err := os.Open(filename)
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read configuration file [%s]; %v\n", filename, err)
	}

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("configuration file was not found, please run `setup` command first\n")
		}
		return fmt.Errorf("error reading config file\n")
	}

	if err := yaml.Unmarshal(content, &c); err != nil {
		return fmt.Errorf("error parsing twig config file\n")
	}

	for version := range c.Php {
		c.versions = append(c.versions, version)
	}

	return nil
}
