package main

import (
	"fmt"
	"os"

	"germa66/internal/config"
	"germa66/internal/meiliclient"
	"germa66/internal/utils"

	"github.com/spf13/cobra"
)

var (
	inputFile  string
	configPath string
	version    = "dev" // Will be set during build
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "germa66",
		Short:   "A CLI application for converting and importing dictionaries to MeiliSearch",
		Long:    "germa66 converts BGL dictionary files to CSV format and imports them into MeiliSearch for fast searching.",
		Version: version,
		RunE:    run,
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "./.env", "Path to configuration file")
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "./import/deutsch_spanisch.BGL", "Input file path")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Validate input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Initialize configuration
	conf, err := config.InitConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}

	// Create MeiliSearch client
	client := meiliclient.New(conf)

	// Perform health check
	if !client.HealthCheck() {
		return fmt.Errorf("MeiliSearch health check failed")
	}

	// Import dictionary
	if err := client.ImportDictionary(inputFile); err != nil {
		return fmt.Errorf("failed to import dictionary: %w", err)
	}

	utils.LogInfo("Dictionary import completed successfully")
	return nil
}
