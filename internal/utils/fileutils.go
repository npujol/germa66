package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// RunPyGlossary runs the pyglossary command and converts the file to CSV
func RunPyGlossary(inputFile, outputFile string) error {
	cmd := exec.Command("pyglossary", inputFile, outputFile, "--write-format=Csv")

	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("pyglossary error: %v, output: %s", err, out.String())
	}

	clean, err := cleanCSVFile(outputFile)

	if err != nil {
		return err
	}
	os.Remove(outputFile)

	err = os.WriteFile(outputFile, clean.Bytes(), 0644)
	if err != nil {
		return err
	}

	LogInfo("pyglossary output:", out.String())
	return nil
}

func cleanCSVFile(filePath string) (bytes.Buffer, error) {
	var clean bytes.Buffer

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return clean, fmt.Errorf("file %s does not exist", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return clean, err
	}

	var row_length int

	for _, line := range bytes.Split(content, []byte("\n")) {
		length := len(bytes.Split(line, []byte(",")))
		if row_length == 0 {
			row_length = length
		}

		if length != row_length {
			continue

		}

		if len(bytes.TrimSpace(line)) == 0 {
			return clean, fmt.Errorf("empty line found in file %s", filePath)
		}

		clean.Write(line)

	}

	return clean, nil
}
