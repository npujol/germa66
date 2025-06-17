package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// RunPyGlossary runs the pyglossary command to convert a dictionary file to CSV format
func RunPyGlossary(inputFile string) (string, error) {
	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return "", fmt.Errorf("input file does not exist: %s", inputFile)
	}

	outputFile := ChangePathExt(inputFile, ".csv")
	LogInfof("Converting %s to %s using pyglossary", inputFile, outputFile)

	cmd := exec.Command("pyglossary", inputFile, outputFile, "--write-format=Csv")

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pyglossary conversion failed: %w, output: %s", err, output.String())
	}

	LogDebugf("pyglossary conversion completed successfully: %s", output.String())
	return outputFile, nil
}
