package config

import (
	_ "embed"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
)

//go:embed config_stub.stub
var cfgStub string

const (
	configKeyStore = "twig.config"
)

type PHP struct {
	Version        string          `yaml:"version"`
	ExecutablePath string          `yaml:"executable_path"`
	IsDefault      bool            `yaml:"default,omitempty"`
	Semver         *semver.Version `yaml:"-"`
}

type Composer struct {
	ExecutablePath string `yaml:"executable_path"`
}

type Config struct {
	PHP      []*PHP   `yaml:"php"`
	Composer Composer `yaml:"composer"`
	name     string   `yaml:"-"`
	dir      string   `yaml:"-"`
	stub     string   `yaml:"-"`
	store    Store    `yaml:"-"`
}

type Option func(*Config)

func NewConfig(name string, options ...Option) (*Config, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("[CONFIG] unable to get user home directory: %w\n", err)
	}

	if name == "" {
		name = "twig"
	}

	cfg := &Config{
		name: name,
		dir:  dir,
	}

	for _, opt := range options {
		opt(cfg)
	}

	return cfg, nil
}

func WithStub(stub string) Option {
	return func(c *Config) {
		c.stub = stub
	}
}

func WithStore(store Store) Option {
	return func(c *Config) {
		c.store = store
	}
}

func (c *Config) Init() error {
	if err := c.createDir(); err != nil {
		return err
	}
	return c.createFile()
}

func (c *Config) BasePath() string {
	return path.Join(c.dir, "."+c.name)
}

func (c *Config) FilePath() string {
	return path.Join(c.BasePath(), "config.yml")
}

func (c *Config) Read() (err error) {
	if c.store != nil {
		if err := c.store.Get(configKeyStore, c); err != nil {
			return err
		}
	}

	cfgFile, err := os.Open(c.FilePath())
	defer func() {
		err = cfgFile.Close()
	}()

	if err != nil {
		return fmt.Errorf("[CONFIG] error opening config file: %w\n", err)
	}

	buff, err := io.ReadAll(cfgFile)
	if err != nil {
		return fmt.Errorf("[CONFIG] error reading config file: %w\n", err)
	}

	if err := yaml.Unmarshal(buff, c); err != nil {
		return fmt.Errorf("[CONFIG] error parsing config file: %w\n", err)
	}

	err = c.createSortSemver()

	// TODO save the versions to the caching layer
	if c.store != nil {
		if err := c.store.Save(configKeyStore, c); err != nil {
			return err
		}
	}

	return
}

func (c *Config) createDir() (err error) {
	_, err = os.Stat(c.BasePath())

	if os.IsNotExist(err) {
		err = os.MkdirAll(c.BasePath(), os.ModePerm)
	}

	return
}

func (c *Config) createFile() (err error) {
	var file *os.File

	if _, err := os.Stat(c.FilePath()); os.IsNotExist(err) {
		file, err = os.Create(c.FilePath())
	}

	if file != nil {
		defer func() {
			err = file.Close()
		}()

		stub := cfgStub

		if c.stub != "" {
			stub = c.stub
		}

		if _, err = file.WriteString(stub); err != nil {
			return fmt.Errorf("[CONFIG] error creating file: %w\n", err)
		}
	}

	return
}
