package config

import (
	"errors"
	"fmt"
	"germa66/internal/utils"
)

var (
	ErrMissingStringConfig = errors.New("missing required configuration")
	ErrInvalidBatchSize    = errors.New("invalid batch size")
)

const (
	DefaultBatchSize = 20000
	MinBatchSize     = 1
	MaxBatchSize     = 100000
)

type Config struct {
	MeilisearchHost   string
	MeilisearchAPIKey string
	MeiliIndex        string
	BatchSize         int
	Debug             bool
}

// InitConfig initializes and returns the application configuration
func InitConfig(path string) (*Config, error) {
	utils.SetLogger()

	conf, err := setupConfig(path)
	if err != nil {
		return nil, fmt.Errorf("failed to setup config: %w", err)
	}
	return conf, nil
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
	batchSize := prov.GetInt("BATCH_SIZE")
	if batchSize == 0 {
		batchSize = DefaultBatchSize
	}

	conf := &Config{
		MeilisearchHost:   prov.GetString("MEILISEARCH_HOST"),
		MeilisearchAPIKey: prov.GetString("MEILISEARCH_API_KEY"),
		MeiliIndex:        prov.GetString("MEILI_INDEX"),
		Debug:             prov.GetBool("DEBUG"),
		BatchSize:         batchSize,
	}

	if err := conf.validate(); err != nil {
		return nil, err
	}

	return conf, nil
}

// validate ensures all required configuration is present and valid
func (conf *Config) validate() error {
	requiredFields := map[string]string{
		"MEILISEARCH_HOST":    conf.MeilisearchHost,
		"MEILISEARCH_API_KEY": conf.MeilisearchAPIKey,
	}

	for field, value := range requiredFields {
		if value == "" {
			return fmt.Errorf("%w: %s is required", ErrMissingStringConfig, field)
		}
	}

	if conf.BatchSize < MinBatchSize || conf.BatchSize > MaxBatchSize {
		return fmt.Errorf("%w: batch size must be between %d and %d, got %d",
			ErrInvalidBatchSize, MinBatchSize, MaxBatchSize, conf.BatchSize)
	}

	// Set default index if not provided
	if conf.MeiliIndex == "" {
		conf.MeiliIndex = "cards"
	}

	return nil
}
