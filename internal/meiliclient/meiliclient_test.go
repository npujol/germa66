package meiliclient_test

import (
	"germa66/internal/meiliclient"
	"germa66/testdata"

	"testing"
)

var (
	MeiliTestClient meiliclient.MeiliClient
)

func TestMeiliClient(t *testing.T) {
	t.Parallel()

	MeiliTestClient = meiliclient.New(testdata.ConfigFixture())

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
}
