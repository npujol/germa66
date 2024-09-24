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

// NewEnvConfigProvider handles configurations using env variables.
// Return an error when it can't read the *.env file
func NewProvider(path string) (*EnvConfigProvider, error) {
	viper.SetConfigType("env")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Error reading config file, %s", err)
	}

	if viper.GetBool(`debug`) {
		utils.LogInfo("Service RUN on DEBUG mode")
		// Log the Debug severity or above.
		log.SetLevel(log.DebugLevel)
	} else {
		// Log the Info severity or above.
		log.SetLevel(log.InfoLevel)
	}

	prov := &EnvConfigProvider{}

	return prov, nil
}

// GetString returns a string value from the environment
func (*EnvConfigProvider) GetString(key string) string {
	return viper.GetString(key)
}

// GetUint returns a uint value from the environment
func (*EnvConfigProvider) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns a bool value from the environment
func (*EnvConfigProvider) GetBool(key string) bool {
	return viper.GetBool(key)
}
