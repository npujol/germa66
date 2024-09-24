package meiliclient

import (
	"fmt"
	"germa66/internal/config"
	"germa66/internal/utils"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

type MeiliClient interface {
	HealthCheck() bool
	ImportDictionary(data []map[string]interface{}) error
}

type Service struct {
	index  meilisearch.IndexManager
	client meilisearch.ServiceManager
	conf   *config.Config
}

// New creates a new MeiliClient using the provided configuration,
// connects to the MeiliSearch instance and creates the index if it doesn't exist.
func New(conf *config.Config, iname string) *Service {
	utils.LogInfo(fmt.Sprintf("Creating connection to Meilisearch on %s", conf.MeilisearchHost))
	client := meilisearch.New(conf.MeilisearchHost, meilisearch.WithAPIKey(conf.MeilisearchAPIKey))
	index := client.Index(iname)

	return &Service{
		client: client,
		index:  index,
		conf:   conf,
	}
}

// HealthCheck checks the health of the MeiliSearch instance.
func (mc *Service) HealthCheck() bool {
	return mc.client.IsHealthy()
}

func (mc *Service) ImportDictionary(data []map[string]interface{}) error {
	// Add the documents from the CSV to the index
	update, err := mc.index.AddDocuments(data)
	if err != nil {
		utils.LogFatalf("Error adding documents to MeiliSearch index: %v", err)
	}

	fmt.Printf("Documents added with update ID: %d\n", update.TaskUID)

	return nil
}
