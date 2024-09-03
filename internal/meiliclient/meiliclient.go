package meiliclient

import (
	"fmt"
	"germa66/internal/config"
	"germa66/internal/utils"
	"os"

	meilisearch "github.com/meilisearch/meilisearch-go"
	log "github.com/sirupsen/logrus"
)

type MeiliClient interface {
	HealthCheck() bool
	ImportDictionary(filepath string) error
}

type Service struct {
	index  meilisearch.IndexManager
	client meilisearch.ServiceManager
	conf   *config.Config
}

// New creates a new MeiliClient using the provided configuration,
// connects to the MeiliSearch instance and creates the index if it doesn't exist.
func New(conf *config.Config, iname string) *Service {
	log.Infof("Creating connection to Meilisearch on %s", conf.MeilisearchHost)
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

func (mc *Service) ImportDictionary(filepath string) error {

	csvData, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading CSV file: %v", err)
	}

	utils.LogInfo(
		fmt.Sprintf(`Adding %d documents to MeiliSearch index...`,
			len(csvData)),
	)

	// Define the CsvDocumentsQuery options if needed (e.g., primary key)
	query := &meilisearch.CsvDocumentsQuery{
		// PrimaryKey: "your_primary_key_column_name",
		CsvDelimiter: ",",
	}

	// Add the documents from the CSV to the index
	update, err := mc.index.AddDocumentsCsv(csvData, query)
	if err != nil {
		log.Fatalf("Error adding documents to MeiliSearch index: %v", err)
	}

	fmt.Printf("Documents added with update ID: %d\n", update.TaskUID)

	return nil
}
