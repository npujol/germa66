package meiliclient

import (
	"encoding/csv"
	"fmt"
	"germa66/internal/config"
	"germa66/internal/models"
	"germa66/internal/utils"
	"io"
	"os"
	"strings"
	"time"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

type MeiliClient interface {
	HealthCheck() bool
	ImportDictionary(path string) error
	batchProcess(filePath string) ([][]models.Card, error)
}

type Service struct {
	client meilisearch.ServiceManager
	conf   *config.Config
}

// New creates a new MeiliClient using the provided configuration,
// connects to the MeiliSearch instance and creates the index if it doesn't exist.
func New(conf *config.Config) *Service {
	utils.LogInfo(fmt.Sprintf("Creating connection to Meilisearch on %s", conf.MeilisearchHost))
	utils.LogInfo("Creating index %s" + conf.MeiliIndex)

	client := meilisearch.New(conf.MeilisearchHost, meilisearch.WithAPIKey(conf.MeilisearchAPIKey))

	return &Service{
		client: client,
		conf:   conf,
	}
}

// HealthCheck checks the health of the MeiliSearch instance.
func (mc *Service) HealthCheck() bool {
	return mc.client.IsHealthy()
}

func (mc *Service) ImportDictionary(path string) error {
	initTime := time.Now()
	// Batch processing of the CSV file
	batches, err := mc.batchProcess(path)
	if err != nil {
		utils.LogError(fmt.Sprintf("Error processing CSV: %v", err))

		return err
	}
	indexName := utils.GetPathName(path)
	utils.LogInfo(fmt.Sprintf("Index name: %s", indexName))

	index := mc.client.Index(indexName)

	for key, batch := range batches {
		batchInitTime := time.Now()

		utils.LogInfo(fmt.Sprintf("Uploading Batch %d...\n", key+1))

		task, err := index.AddDocuments(batch)
		if err != nil {
			utils.LogError(fmt.Sprintf("Error adding documents to Meilisearch: %v\n", err))

			return err
		}

		utils.LogInfo(fmt.Sprintf(
			"Batch %d successfully uploaded with TaskUID: %d with %d documents took %s\n",
			key+1, task.TaskUID, len(batch), time.Since(batchInitTime)))
	}

	utils.LogInfo(fmt.Sprintf("Total upload time: %s\n", time.Since(initTime)))

	return nil
}

func (mc *Service) batchProcess(filePath string) ([][]models.Card, error) {
	initTime := time.Now()
	outputFile, err := utils.RunPyGlossary(filePath)
	if err != nil {
		utils.LogFatalf("Error running pyglossary: %v", err)
	}

	utils.LogInfo("Conversion complete.")

	fileName := utils.GetPathName(outputFile)
	file, err := os.Open(outputFile)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	reader.FieldsPerRecord = -1

	_, err = reader.Read()

	if err != nil {
		utils.LogWarn(fmt.Sprintf("Warning: skipping header due to error: %v\n", err))

		return nil, err
	}

	var batches [][]models.Card

	var batch []models.Card

	wrongCount := 0

	i := 0
	for record, err := reader.Read(); err != io.EOF; record, err = reader.Read() {
		i++
		if err != nil {
			utils.LogDebug(fmt.Sprintf("Warning: skipping row due to error: %v\n", err))

			continue
		}

		for i := range record {
			record[i] = strings.TrimSpace(record[i])
		}

		product, pErr := models.RowToCard(record, fileName, i)

		if pErr != nil {
			wrongCount++

			continue
		}

		batch = append(batch, product)

		if len(batch) >= mc.conf.BatchSize {
			batches = append(batches, batch)
			batch = []models.Card{}
		}
	}

	utils.LogInfo("Reached EOF, exiting...")

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	utils.LogInfo(fmt.Sprintf("Total time processing batches: %s\n", time.Since(initTime)))

	if wrongCount > 0 {
		utils.LogInfo(fmt.Sprintf("Skipped %d rows due to errors\n", wrongCount))
	}

	utils.LogInfo(fmt.Sprintf("Total time processing CSV with %d batches: %s\n", len(batches), time.Since(initTime)))

	if len(batches) == 0 {
		utils.LogInfo("No batches to upload")

		return nil, nil
	}
	return batches, nil
}
