package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetPreferredVersion(t *testing.T) {
	cfg := createTestConfig(t)

	tests := []struct {
		name        string
		constraints string
		want        string
	}{
		{
			name:        "should get the PHP 8.1",
			constraints: "^7.4 || ^8.0",
			want:        "7.4",
		},
		{
			name:        "should get the PHP 7.4",
			constraints: "^7.2 || ^7.4",
			want:        "7.2",
		},
		{
			name:        "should get the PHP 7.2",
			constraints: "<=7.2",
			want:        "7.2",
		},
		{
			name:        "should get the PHP 8.0",
			constraints: "^8.0 || ^8.1",
			want:        "8.0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			php, err := cfg.GetPreferredPHPVersion(test.constraints)

			assert.NoError(t, err)
			assert.NotNil(t, php)
			assert.Equal(t, test.want, php.Version)
		})
	}

	t.Cleanup(func() {
		os.RemoveAll(cfg.BasePath())
	})
}
