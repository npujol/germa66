package main

import (
	"fmt"
	"germa66/internal/config"
	"germa66/internal/meiliclient"
	"germa66/internal/utils"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
	configPath string
)

func main() {
	// Initialize Cobra root command
	rootCmd := &cobra.Command{
		Use:   "app",
		Short: "A CLI application for converting and importing dictionaries",
		Run:   run,
	}

	// Define flags for the root command
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "./.env", "Path to configuration file")
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "./import/deutsch_spanisch.BGL", "Input file path")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "./import/deutsch_spanisch.csv", "Output file path")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		utils.LogFatalf("Error executing command: %v", err)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Initialize configuration
	conf := config.InitConfig(configPath)

	// Run the pyglossary command
	if err := utils.RunPyGlossary(inputFile, outputFile); err != nil {
		utils.LogFatalf("Error running pyglossary: %v", err)
	}

	utils.LogInfo("Conversion complete.")

	fileName := strings.Split(strings.Split(outputFile, "/")[len(strings.Split(outputFile, "/"))-1], ".")[0]

	// Import the dictionary concurrently
	var wg sync.WaitGroup
	errCh := make(chan error, 1) // Buffer size 1 to handle a single error

	wg.Add(1)
	go func() {
		defer wg.Done()
		meili := meiliclient.New(conf, fileName)
		if err := meili.ImportDictionary(outputFile); err != nil {
			errCh <- fmt.Errorf("Error importing dictionary: %v", err)
		}
	}()

	// Wait for the import to complete or fail
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Handle any errors from the import process
	if err, ok := <-errCh; ok {
		utils.LogFatalf("%v", err)
	}

	// Clean up the output file
	if err := os.Remove(outputFile); err != nil {
		utils.LogFatalf("Error removing output file: %v", err)
	}
}
