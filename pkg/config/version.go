package config

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"sort"
)

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

func (c *Config) createSortSemver() error {
	for k, php := range c.PHP {
		version, err := semver.NewVersion(php.Version)
		if err != nil {
			return err
		}
		c.PHP[k].Semver = version
	}

	sort.Slice(c.PHP, func(i, j int) bool {
		return c.PHP[i].Semver.GreaterThan(c.PHP[j].Semver)
	})

	return nil
}
