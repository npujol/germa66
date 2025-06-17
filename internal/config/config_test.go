package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"germa66/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_New(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			envVars: map[string]string{
				"MEILISEARCH_HOST":    "http://localhost:7700",
				"MEILISEARCH_API_KEY": "test-key",
				"BATCH_SIZE":          "1000",
				"DEBUG":               "true",
			},
			expectError: false,
		},
		{
			name: "missing host",
			envVars: map[string]string{
				"MEILISEARCH_API_KEY": "test-key",
			},
			expectError: true,
			errorMsg:    "MEILISEARCH_HOST is required",
		},
		{
			name: "missing API key",
			envVars: map[string]string{
				"MEILISEARCH_HOST": "http://localhost:7700",
			},
			expectError: true,
			errorMsg:    "MEILISEARCH_API_KEY is required",
		},
		{
			name: "invalid batch size",
			envVars: map[string]string{
				"MEILISEARCH_HOST":    "http://localhost:7700",
				"MEILISEARCH_API_KEY": "test-key",
				"BATCH_SIZE":          "200000", // Too large
			},
			expectError: true,
			errorMsg:    "batch size must be between",
		},
		{
			name: "default values",
			envVars: map[string]string{
				"MEILISEARCH_HOST":    "http://localhost:7700",
				"MEILISEARCH_API_KEY": "test-key",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tempFile := createTempConfigFile(t, tt.envVars)
			defer os.Remove(tempFile)

			// Test the config initialization
			conf, err := config.InitConfig(tempFile)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, conf)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, conf)
				assert.Equal(t, tt.envVars["MEILISEARCH_HOST"], conf.MeilisearchHost)
				assert.Equal(t, tt.envVars["MEILISEARCH_API_KEY"], conf.MeilisearchAPIKey)
				
				// Check default values
				if tt.envVars["BATCH_SIZE"] == "" {
					assert.Equal(t, config.DefaultBatchSize, conf.BatchSize)
				}
				if conf.MeiliIndex == "" {
					assert.Equal(t, "cards", conf.MeiliIndex)
				}
			}
		})
	}
}

func createTempConfigFile(t *testing.T, envVars map[string]string) string {
	t.Helper()
	
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.env")
	
	var content string
	for key, value := range envVars {
		content += key + "=" + value + "\n"
	}
	
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)
	
	return tmpFile
}
