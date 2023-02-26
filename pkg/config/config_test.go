package config

import (
	"github.com/Masterminds/semver/v3"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path"
	"testing"
)

func TestConfig_Path(t *testing.T) {
	home, err := homedir.Dir()

	assert.NoError(t, err)

	tests := []struct {
		name  string
		want  string
		value string
	}{
		{
			name: "test default name",
			want: path.Join(home, ".twig"),
		},
		{
			name:  "test custom name",
			want:  path.Join(home, ".custom"),
			value: "custom",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := NewConfig(test.value)

			assert.NoError(t, err)
			assert.Equal(t, test.want, cfg.BasePath())
		})
	}
}

func TestConfig_NewWithOption(t *testing.T) {
	cfg, err := NewConfig("twig-test", WithStub("the stub"))

	assert.NoError(t, err)
	assert.Equal(t, "the stub", cfg.stub)
}

func TestConfig_Init(t *testing.T) {
	cfg := createTestConfig(t)

	// Remove the configuration files & folder
	defer func() {
		if err := CleanupTestFiles("twig-test"); err != nil {
			assert.NoError(t, err)
		}
	}()

	if _, err := os.Stat(cfg.FilePath()); os.IsNotExist(err) {
		t.Fatal("expect configuration file was created but it doesn't")
	}

	file, err := os.Open(cfg.FilePath())

	assert.NoError(t, err, "failed to open the config file")

	defer func() {
		assert.NoError(t, file.Close(), "error when closing the file")
	}()

	buff, err := io.ReadAll(file)

	assert.NoError(t, err, "failed to read the config file")
	assert.Equal(t, cfgStub, string(buff))
}

func TestConfig_ReadingConfiguration(t *testing.T) {
	cfg, err := NewConfig("twig-test", WithStub(`php:
  - version: 7.2
    executable_path: /your/path/to/php7.2
    default: true
  - version: 7.4
    executable_path: /your/path/to/php7.4
  - version: 8.0
    executable_path: /your/path/to/php8.0
  - version: 8.1
    executable_path: /your/path/to/php8.1
composer:
  executable_path: /your/composer.phar/path`))

	assert.NoError(t, err)
	assert.NoError(t, cfg.Init())

	// Remove the configuration files & folder
	defer func() {
		assert.NoError(t, CleanupTestFiles("twig-test"))
	}()

	assert.NoError(t, err, "error during config initialization")
	assert.NoError(t, cfg.Read())

	assert.NotEmpty(t, cfg.Composer.ExecutablePath, "expected composer executable path to exists but it doesn't")
	assert.Lenf(t, cfg.PHP, 4, "len of PHP should be 4 but got (%d)", len(cfg.PHP))

	orders := []string{"8.1", "8.0", "7.4", "7.2"}

	for k, php := range cfg.PHP {
		version, err := semver.NewVersion(orders[k])
		assert.NoError(t, err)
		assert.Truef(t, version.Equal(php.Semver), "want (%s) got (%s)", version, php.Semver)
	}
}

func TestConfig_WithCacheStore(t *testing.T) {
	cfg, err := NewConfig("twig-test", WithStore(NewInMemStore()))

	// Remove the configuration files & folder
	defer func() {
		assert.NoError(t, CleanupTestFiles("twig-test"))
	}()

	assert.NoError(t, err)
	assert.NoError(t, cfg.Init())
	assert.NoError(t, cfg.Read())

	assert.True(t, cfg.store.Has(configKeyStore))
	assert.NotNil(t, cfg.PHP)
	assert.NotNil(t, cfg.Composer)
	assert.NotEmpty(t, cfg.name)
	assert.NotEmpty(t, cfg.dir)
}

func createTestConfig(t *testing.T) *Config {
	cfg, err := NewConfig("twig-test")

	assert.NoError(t, err)
	assert.NoError(t, cfg.Init())
	assert.NoError(t, cfg.Read())

	return cfg
}
