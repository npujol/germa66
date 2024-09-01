package main

import (
	"bytes"
	"fmt"
	"germa66/internal/config"
	"germa66/internal/meiliclient"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func InitConfig() *config.Config {

	// Setup logger
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Setup the app's config
	conf, confErr := setupConfig()
	if confErr != nil {
		log.Fatalf("Could not setup the config, due to %s.\n", confErr)
	}
	return conf
}

func setupConfig() (*config.Config, error) {
	path := "./.env"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	prov, err := config.NewProvider(path)
	if err != nil {
		return nil, err
	}

	conf, err2 := config.New(prov)
	if err2 != nil {
		return nil, err2
	}

	return conf, nil
}


func main() {

	conf := InitConfig()
	// Input and output file paths
	inputFile := "./data/deutsch_spanisch.BGL"
	outputFile := "./data/deutsch_spanisch.csv"

	// Run the pyglossary command
	err := runPyGlossary(inputFile, outputFile)
	if err != nil {
		log.Fatalf("Error running pyglossary: %v", err)
	}

	log.Println("Conversion complete.")

	meili := meiliclient.New(conf)
	importErr := meili.ImportDictionary(outputFile)

	if importErr != nil {
		log.Fatalf("Error importing dictionary: %v", importErr)
	}
}

// runPyGlossary runs the pyglossary command and converts the file to CSV
func runPyGlossary(inputFile, outputFile string) error {
	// Prepare the pyglossary command
	cmd := exec.Command("pyglossary", inputFile, outputFile, "--write-format=Csv")

	// Get the output from the command
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("pyglossary error: %v, output: %s", err, out.String())
	}

	// Read the output to confirm successful execution
	log.Println("pyglossary output:", out.String())

	return nil
}