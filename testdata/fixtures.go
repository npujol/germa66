package testdata

import (
	"path/filepath"
	"testing"

	"germa66/internal/config"
)

// ConfigFixture returns a config fixture for testing
func ConfigFixture(t *testing.T) *config.Config {
	t.Helper()
	
	path, err := filepath.Abs("../../testdata/test.env")
	if err != nil {
		t.Fatalf("Failed to get absolute path for test config: %v", err)
	}
	
	conf, err := config.InitConfig(path)
	if err != nil {
		t.Fatalf("Failed to initialize test config: %v", err)
	}
	
	return conf
}

// DictionaryFixturePath returns the path to the test dictionary file
func DictionaryFixturePath(t *testing.T) string {
	t.Helper()
	
	path, err := filepath.Abs("../../testdata/deutsch_spanisch.BGL")
	if err != nil {
		t.Fatalf("Failed to get absolute path for test dictionary: %v", err)
	}
	
	return path
}
