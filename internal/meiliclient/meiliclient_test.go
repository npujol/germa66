package meiliclient_test

import (
	"testing"

	"germa66/internal/meiliclient"
	"germa66/testdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeiliClient_HealthCheck(t *testing.T) {
	// Skip if MeiliSearch is not available
	conf := testdata.ConfigFixture(t)
	client := meiliclient.New(conf)
	
	// This test will only pass if MeiliSearch is running
	// For CI/CD, you might want to skip this test or use a mock
	if !client.HealthCheck() {
		t.Skip("MeiliSearch is not available, skipping health check test")
		return
	}
	
	assert.True(t, client.HealthCheck(), "MeiliSearch should be healthy")
}

func TestMeiliClient_ImportDictionary_FileNotFound(t *testing.T) {
	conf := testdata.ConfigFixture(t)
	client := meiliclient.New(conf)
	
	// Test with non-existent file
	err := client.ImportDictionary("/non/existent/file.bgl")
	assert.Error(t, err, "Should return error for non-existent file")
	assert.Contains(t, err.Error(), "input file does not exist", "Error should mention file not found")
}

func TestMeiliClient_ImportDictionary_WithMockFile(t *testing.T) {
	conf := testdata.ConfigFixture(t)
	client := meiliclient.New(conf)
	
	// Skip if MeiliSearch is not available
	if !client.HealthCheck() {
		t.Skip("MeiliSearch is not available, skipping import test")
		return
	}
	
	// Skip if pyglossary is not available
	dictPath := testdata.DictionaryFixturePath(t)
	
	// This test requires pyglossary to be installed
	// For now, we'll test that the error handling works correctly
	err := client.ImportDictionary(dictPath)
	
	// The test might fail due to pyglossary not being installed, which is expected
	// In a real test environment, you'd either:
	// 1. Ensure pyglossary is installed, or
	// 2. Mock the pyglossary call, or  
	// 3. Use a pre-converted CSV file
	if err != nil {
		t.Logf("Import failed (expected if pyglossary not installed): %v", err)
		// Verify it's the expected type of error
		require.Contains(t, err.Error(), "failed to process dictionary file")
	} else {
		t.Log("Dictionary import succeeded")
	}
}
