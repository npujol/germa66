package meiliclient_test

import (
	"germa66/internal/meiliclient"
	"germa66/testdata"

	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	MeiliTestClient meiliclient.MeiliClient
)

func TestMeiliClient(t *testing.T) {
	t.Parallel()

	MeiliTestClient = meiliclient.New(testdata.ConfigFixture(), "test")

	testsHealthCheck := []struct {
		name string
	}{
		{
			"it_should_be_healthy",
		},
	}

	for _, test := range testsHealthCheck {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			MeiliTestClient.HealthCheck()
		})
	}

	testsImportDictionary := []struct {
		name string
	}{
		{
			"it_should_be_healthy",
		},
	}

	for _, test := range testsImportDictionary {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			content := []map[string]interface{}{
				{
					"test":   "test",
					"test_2": 1,
				},
			}
			e := MeiliTestClient.ImportDictionary(content)
			assert.NoError(t, e, "Error importing dictionary")
		})
	}
}
