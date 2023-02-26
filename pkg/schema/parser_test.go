package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name string
		want string
		data []byte
		err  error
	}{
		{
			name: "get version from require",
			want: "^8.0",
			data: []byte(`{
				"require": {
					"php": "^8.0"
				}
			}`),
		},
		{
			name: "use config platform as fallback",
			want: "^8.1",
			data: []byte(`{
				"require": {
					"package": "1.0.0"
				},
				"config": {
					"platform": {
						"php": "^8.1"
					}
				}
			}`),
		},
		{
			name: "composer.json with multiple version constraints",
			want: "^7.4||^8.0",
			data: []byte(`{
				"require": {
					"php": "^7.4|^8.0"
				}
			}`),
		},
		{
			name: "no version could be found",
			err:  ErrNoVersionDefined,
			data: []byte(`{
				"require": {
					"package": "1.0.0"
				}
			}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			version, err := GetVersion(test.data)

			assert.ErrorIs(t, err, test.err)
			assert.Equalf(t, test.want, version, "want (%s) got (%s)", test.want, version)
		})
	}
}
