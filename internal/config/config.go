package config

import (
	"errors"
)

var ErrMissingStringConfig = errors.New("missing string configuration")

type Config struct {
	MeilisearchHost   string
	MeilisearchAPIKey string
	MeiliIndex        string
	Debug             bool
}

// New creates a new Config using the provided EnvProvider.
func New(prov EnvProvider) (*Config, error) {
	conf := &Config{
		MeilisearchHost:   prov.GetString("MEILISEARCH_HOST"),
		MeilisearchAPIKey: prov.GetString("MEILISEARCH_API_KEY"),
		MeiliIndex:        prov.GetString("MEILI_INDEX"),
		Debug:             prov.GetBool("DEBUG"),
	}

	err := conf.ensureData()

	return conf, err
}

func (conf *Config) ensureData() error {
	for _, val := range map[string]string{
		"MeilisearchHost":   conf.MeilisearchHost,
		"MeilisearchAPIKey": conf.MeilisearchAPIKey,
		"MeiliIndex":        conf.MeiliIndex,
	} {
		if val == "" {
			return ErrMissingStringConfig
		}
	}

	return nil
}
