package config

type Store interface {
	Has(key string) bool
	Save(key string, cfg *Config) error
	Get(key string, cfg *Config) error
	Flush()
}
