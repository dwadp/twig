package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemStore_GetSave(t *testing.T) {
	store := NewInMemStore()
	savedCfg := &Config{
		PHP: []*PHP{
			{
				Version:        "8.1",
				ExecutablePath: "./php",
				IsDefault:      false,
			},
		},
		Composer: Composer{
			ExecutablePath: "./composer",
		},
		name:  "twig-test",
		dir:   "home",
		store: store,
	}

	err := store.Save(configKeyStore, savedCfg)

	assert.NoError(t, err)
	assert.True(t, store.Has(configKeyStore))

	cfg := Config{}
	assert.NoError(t, store.Get(configKeyStore, &cfg))

	assert.Len(t, cfg.PHP, len(savedCfg.PHP))
	assert.Equal(t, savedCfg.Composer, cfg.Composer)

	store.Flush()
	assert.False(t, store.Has(configKeyStore))
}

func TestInMemStore_Exists(t *testing.T) {
	store := NewInMemStore()

	assert.False(t, store.Has("not_exist"))

	err := store.Save(configKeyStore, &Config{
		PHP: []*PHP{
			{
				Version:        "8.1",
				ExecutablePath: "./php",
				IsDefault:      false,
			},
		},
		Composer: Composer{
			ExecutablePath: "./composer",
		},
		name:  "twig-test",
		dir:   "home",
		store: store,
	})

	assert.NoError(t, err)
	assert.True(t, store.Has(configKeyStore))

	store.Flush()
	assert.False(t, store.Has(configKeyStore))
}
