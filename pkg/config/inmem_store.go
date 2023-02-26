package config

import (
	"errors"
	"fmt"
	"github.com/bluele/gcache"
	"gopkg.in/yaml.v3"
)

type InMemStore struct {
	store gcache.Cache
}

func NewInMemStore() *InMemStore {
	gc := gcache.New(20).
		LRU().
		Build()

	return &InMemStore{store: gc}
}

func (im InMemStore) Has(key string) bool {
	return im.store.Has(key)
}

func (im InMemStore) Save(key string, cfg *Config) error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return im.store.Set(key, b)
}

func (im InMemStore) Get(key string, cfg *Config) error {
	data, err := im.store.Get(key)

	if err != nil {
		if !errors.Is(err, gcache.KeyNotFoundError) {
			return err
		}
	}

	if data == nil {
		return nil
	}

	b, ok := data.([]byte)

	if !ok {
		return fmt.Errorf("invalid data")
	}

	if err := yaml.Unmarshal(b, cfg); err != nil {
		return err
	}

	return nil
}

func (im InMemStore) Flush() {
	im.store.Purge()
}
