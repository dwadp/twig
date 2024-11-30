package config

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
)

// GetPreferredPHPVersion will determine which PHP version to obtain from the system based on the composer.json
// It will automatically use the newer version if there's a definition like (7 || 8)
func (c *Config) GetPreferredPHPVersion(constraint string) (*PHP, error) {
	constraints, err := semver.NewConstraint(constraint)
	if err != nil {
		return nil, err
	}

	for _, php := range c.PHP {
		if constraints.Check(php.Semver) {
			return php, nil
		}
	}

	return nil, fmt.Errorf("could not find any match for contstraints [%s]\n", constraints)
}

func (c *Config) GetDefaultPHPVersion() (*PHP, error) {
	for _, php := range c.PHP {
		if php.IsDefault {
			return php, nil
		}
	}
	return nil, fmt.Errorf("could not find any default PHP version")
}

func (c *Config) createSortSemver() error {
	for k, php := range c.PHP {
		version, err := semver.NewVersion(php.Version)
		if err != nil {
			return err
		}
		c.PHP[k].Semver = version
	}

	sort.Slice(c.PHP, func(i, j int) bool {
		return c.PHP[i].Semver.LessThan(c.PHP[j].Semver)
	})

	return nil
}
