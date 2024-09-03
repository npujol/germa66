package testdata

import (
	"log"
	"path/filepath"

	"germa66/internal/config"
)

// ConfigFixture returns a config fixture
func ConfigFixture() *config.Config {
	path, err := filepath.Abs("../../testdata/test.env")
	if err != nil {
		log.Fatal(err)
	}
	prov, pErr := config.NewProvider(path)
	if pErr != nil {
		log.Fatal(pErr)
	}
	conf, cErr := config.New(prov)
	if cErr != nil {
		log.Fatal(cErr)
	}
	return conf
}

func DictionaryFixturePath() string {
	path, err := filepath.Abs("../../testdata/import/deutsch_spanisch.BGL")
	if err != nil {
		log.Fatal(err)
	}
	return path
}
