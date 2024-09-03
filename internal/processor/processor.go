package processor

import (
	"bytes"
	"fmt"
	"germa66/internal/config"
	"germa66/internal/meiliclient"
	"germa66/internal/utils"
	"os"
	"sync"
)

// Process handles the conversion and import of the dictionary.
func Process(configPath, inputFile string) {
	// Initialize configuration
	conf := config.InitConfig(configPath)

	outputFile, err := utils.RunPyGlossary(inputFile)
	if err != nil {
		utils.LogFatalf("Error running pyglossary: %v", err)
	}

	utils.LogInfo("Conversion complete.")

	fileName := utils.GetPathName(outputFile)

	// Clean and process the CSV file in chunks
	err = processCSVFile(outputFile, func(chunk [][]byte) error {
		return processChunk(conf, chunk, fileName)
	})
	if err != nil {
		utils.LogFatalf("Error processing CSV file: %v", err)
	}

	// Clean up the output file
	if err := os.Remove(outputFile); err != nil {
		utils.LogFatalf("Error removing output file: %v", err)
	}
}


// processChunk processes a chunk of CSV data, for example, by sending it to MeiliSearch or another service.
func processChunk(conf *config.Config, chunk [][]byte, fileName string) error {
	// Initialize the MeiliSearch client
	meili := meiliclient.New(conf, fileName)

	// Create a buffer to accumulate the data for the chunk
	var dataBuffer []map[string]interface{}

	// Process each line in the chunk
	for _, line := range chunk {
		// Here you would parse the line and potentially transform it
		// For simplicity, let's assume you just write it directly to the buffer
		dataBuffer = append(dataBuffer, map[string]interface{}{"content": string(line)})
	}

	// Now, you might send this chunk of data to MeiliSearch
	if err := meili.ImportDictionary(dataBuffer); err != nil {
		return fmt.Errorf("Error importing chunk to MeiliSearch: %v", err)
	}

	utils.LogInfo(fmt.Sprintf("Processed chunk with %d lines", len(chunk)))

	return nil
}

// processCSVFile processes a CSV file and calls the provided processChunk function for each chunk of data.
func processCSVFile(filePath string, processChunk func([][]byte) error) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file error: %v", err)
	}

	var wg sync.WaitGroup
	var rowLength int
	chunkSize := 2000
	var chunk [][]byte

	for _, line := range bytes.Split(content, []byte("\n")) {
		columns := bytes.Split(line, []byte(","))
		if rowLength == 0 {
			rowLength = len(columns)
		}

		if len(columns) != rowLength || len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		chunk = append(chunk, line)

		// When chunk size is reached, process the chunk
		if len(chunk) == chunkSize {
			wg.Add(1)
			go func(chunkToProcess [][]byte) {
				defer wg.Done()
				if err := processChunk(chunkToProcess); err != nil {
					fmt.Printf("Error processing chunk: %v\n", err)
				}
			}(chunk)
			chunk = nil // Reset chunk
		}
	}

	// Process any remaining lines in the last chunk
	if len(chunk) > 0 {
		wg.Add(1)
		go func(chunkToProcess [][]byte) {
			defer wg.Done()
			if err := processChunk(chunkToProcess); err != nil {
				fmt.Printf("Error processing chunk: %v\n", err)
			}
		}(chunk)
	}

	wg.Wait()

	return nil
}
