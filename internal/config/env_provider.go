package config

import (
	"fmt"
	"germa66/internal/utils"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type EnvProvider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}

type EnvConfigProvider struct {
}

// NewProvider creates a new environment configuration provider.
// It reads configuration from the specified file path and sets up logging.
func NewProvider(path string) (*EnvConfigProvider, error) {
	viper.SetConfigType("env")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	// Configure logging based on debug setting
	if viper.GetBool("DEBUG") {
		utils.LogInfo("Service running in DEBUG mode")
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	return &EnvConfigProvider{}, nil
}

// GetString returns a string value from the environment
func (*EnvConfigProvider) GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns an integer value from the environment
func (*EnvConfigProvider) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns a boolean value from the environment
func (*EnvConfigProvider) GetBool(key string) bool {
	return viper.GetBool(key)
}
