package config

import (
	"errors"
	"germa66/internal/utils"
)

var ErrMissingStringConfig = errors.New("missing string configuration")

type Config struct {
	MeilisearchHost   string
	MeilisearchAPIKey string
	MeiliIndex        string
	Debug             bool
}

// InitConfig initializes and returns the application configuration
func InitConfig(path string) *Config {
	utils.SetLogger()

	conf, err := setupConfig(path)
	if err != nil {
		utils.LogFatalf("Could not setup the config, due to %s.\n", err)
	}
	return conf
}

func setupConfig(path string) (*Config, error) {
	prov, err := NewProvider(path)
	if err != nil {
		return nil, err
	}

	conf, err := New(prov)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// New creates a new Config using the provided EnvProvider.
func New(prov EnvProvider) (*Config, error) {
	conf := &Config{
		MeilisearchHost:   prov.GetString("MEILISEARCH_HOST"),
		MeilisearchAPIKey: prov.GetString("MEILISEARCH_API_KEY"),
		Debug:             prov.GetBool("DEBUG"),
	}

	err := conf.ensureData()

	return conf, err
}

func (conf *Config) ensureData() error {
	for _, val := range map[string]string{
		"MeilisearchHost":   conf.MeilisearchHost,
		"MeilisearchAPIKey": conf.MeilisearchAPIKey,
	} {
		if val == "" {
			return ErrMissingStringConfig
		}
	}

	return nil
}
