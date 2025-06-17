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
}

type Service struct {
	index  meilisearch.IndexManager
	client meilisearch.ServiceManager
	conf   *config.Config
}

// New creates a new MeiliClient using the provided configuration.
// It establishes a connection to MeiliSearch and sets up the index.
func New(conf *config.Config) *Service {
	utils.LogInfo(fmt.Sprintf("Creating connection to Meilisearch on %s", conf.MeilisearchHost))

	client := meilisearch.New(conf.MeilisearchHost, meilisearch.WithAPIKey(conf.MeilisearchAPIKey))
	index := client.Index(conf.MeiliIndex)

	// Set primary key for the index
	if _, err := index.UpdateIndex("id"); err != nil {
		utils.LogWarn(fmt.Sprintf("Warning: failed to set primary key for index: %v", err))
	}

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

func (mc *Service) ImportDictionary(path string) error {
	startTime := time.Now()
	utils.LogInfo(fmt.Sprintf("Starting dictionary import from: %s", path))

	// Convert BGL to CSV and process into batches
	batches, err := mc.batchProcess(path)
	if err != nil {
		return fmt.Errorf("failed to process dictionary file: %w", err)
	}

	utils.LogInfo(fmt.Sprintf("Processed %d batches for upload", len(batches)))

	// Upload batches to MeiliSearch
	for i, batch := range batches {
		if err := mc.uploadBatch(i+1, batch); err != nil {
			return fmt.Errorf("failed to upload batch %d: %w", i+1, err)
		}
	}

	utils.LogInfo(fmt.Sprintf("Dictionary import completed in %s", time.Since(startTime)))
	return nil
}

// uploadBatch uploads a single batch of documents to MeiliSearch
func (mc *Service) uploadBatch(batchNum int, batch []models.Card) error {
	batchStart := time.Now()
	utils.LogInfo(fmt.Sprintf("Uploading batch %d (%d documents)...", batchNum, len(batch)))

	task, err := mc.index.AddDocuments(batch)
	if err != nil {
		return fmt.Errorf("failed to add documents to MeiliSearch: %w", err)
	}

	utils.LogInfo(fmt.Sprintf(
		"Batch %d uploaded successfully (TaskUID: %d) in %s",
		batchNum, task.TaskUID, time.Since(batchStart)))

	return nil
}

// batchProcess converts BGL file to CSV and processes it into batches
func (mc *Service) batchProcess(filePath string) ([][]models.Card, error) {
	startTime := time.Now()

	// Convert BGL to CSV using pyglossary
	csvFile, err := utils.RunPyGlossary(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to convert BGL to CSV: %w", err)
	}
	utils.LogInfo("BGL to CSV conversion completed")

	// Process CSV file into batches
	batches, err := mc.processCSVFile(csvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to process CSV file: %w", err)
	}

	utils.LogInfo(fmt.Sprintf("Batch processing completed in %s", time.Since(startTime)))
	return batches, nil
}

// processCSVFile reads CSV file and creates batches of cards
func (mc *Service) processCSVFile(csvPath string) ([][]models.Card, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	fileName := utils.GetPathName(csvPath)
	var batches [][]models.Card
	var currentBatch []models.Card
	skippedRows := 0
	totalRows := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.LogDebug(fmt.Sprintf("Skipping row due to read error: %v", err))
			skippedRows++
			continue
		}

		totalRows++

		// Trim whitespace from all fields
		for i := range record {
			record[i] = strings.TrimSpace(record[i])
		}

		// Convert row to Card
		card, err := models.RowToCard(record, fileName)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("Skipping invalid row: %v", err))
			skippedRows++
			continue
		}

		currentBatch = append(currentBatch, card)

		// Create new batch when current batch reaches configured size
		if len(currentBatch) >= mc.conf.BatchSize {
			batches = append(batches, currentBatch)
			currentBatch = []models.Card{}
		}
	}

	// Add remaining cards as final batch
	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	utils.LogInfo(fmt.Sprintf("Processed %d total rows, created %d batches, skipped %d invalid rows",
		totalRows, len(batches), skippedRows))

	return batches, nil
}
