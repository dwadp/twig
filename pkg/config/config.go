package config

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/Masterminds/semver/v3"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

//go:embed config.stub
var cfgStub string

const (
	configKeyStore       = "twig.config"
	defaultFileMode      = 0664
	readOnlyFileFlag     = os.O_RDONLY
	readWriteFileFlag    = os.O_RDWR | os.O_TRUNC
	readOrCreateFileFlag = os.O_CREATE | os.O_RDWR
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
		return nil, fmt.Errorf("config: unable to get home user directory: %w\n", err)
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

// TODO: Add WithBaseDir option to modify the base configuration directory

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
	// If there is a configuration stored in the cache, we don't want to manually parse the file
	// to obtain which version is needed in order the program to run
	if c.store != nil {
		if err := c.store.Get(configKeyStore, c); err != nil {
			return err
		}
	}

	cfgFile, err := os.OpenFile(c.FilePath(), readOnlyFileFlag, defaultFileMode)
	if err != nil {
		return fmt.Errorf("config: failed to open file: %w\n", err)
	}

	defer func() {
		err = cfgFile.Close()
	}()

	buf, err := io.ReadAll(cfgFile)
	if err != nil {
		return fmt.Errorf("config: error reading the file: %w\n", err)
	}

	if err := yaml.Unmarshal(buf, c); err != nil {
		return fmt.Errorf("config: error parsing the file: %w\n", err)
	}

	err = c.createSortSemver()

	if c.store != nil {
		if err := c.store.Save(configKeyStore, c); err != nil {
			return err
		}
	}

	return
}

// IsVersionExist will check if the given version exists in the configuration file
func (c *Config) IsVersionExist(version string) bool {
	for _, php := range c.PHP {
		if php.Version == version {
			return true
		}
	}

	return false
}

// SetDefault will set the given version as the default version
func (c *Config) SetDefault(version string) error {
	if err := c.Read(); err != nil {
		return err
	}

	if !c.IsVersionExist(version) {
		return fmt.Errorf("config: version %s does not exist\n", version)
	}

	for _, php := range c.PHP {
		if php.Version == version {
			php.IsDefault = true
		} else {
			php.IsDefault = false
		}
	}

	return c.writeFile()
}

func (c *Config) writeFile() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("config: error marshalling config: %w\n", err)
	}

	file, err := os.OpenFile(c.FilePath(), readWriteFileFlag, defaultFileMode)
	if err != nil {
		return fmt.Errorf("config: error opening the file: %w\n", err)
	}
	defer file.Close()

	if _, err := file.WriteString(string(b)); err != nil {
		return fmt.Errorf("config: error writing to the file: %w\n", err)
	}

	if err := file.Sync(); err != nil {
		return fmt.Errorf("config: error saving the file: %w\n", err)
	}

	return nil
}

// createDir Creates the base directory for the configuration file
func (c *Config) createDir() (err error) {
	_, err = os.Stat(c.BasePath())

	if os.IsNotExist(err) {
		err = os.MkdirAll(c.BasePath(), os.ModePerm)
	}

	return
}

// createFile will create base configuration file based on the stub file defined in config.stub
func (c *Config) createFile() error {
	_, err := os.Stat(c.FilePath())

	if err == nil {
		return fmt.Errorf("config: file already exist\n")
	}

	file, err := os.OpenFile(c.FilePath(), readOrCreateFileFlag, defaultFileMode)
	if err != nil {
		return fmt.Errorf("config: failed to create file: %w\n", err)
	}

	defer file.Close()

	stub := cfgStub
	if c.stub != "" {
		stub = c.stub
	}

	if _, err = file.WriteString(stub); err != nil {
		return fmt.Errorf("config: failed to write stub file: %w\n", err)
	}

	return nil
}
