package main

import (
	"germa66/internal/processor"
	"germa66/internal/utils"

	"github.com/spf13/cobra"
)

var (
	inputFile  string
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

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		utils.LogFatalf("Error executing command: %v", err)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Pass the configuration and file paths to the process function
	processor.Process(configPath, inputFile)
}
