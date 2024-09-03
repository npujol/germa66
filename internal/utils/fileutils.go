package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// RunPyGlossary runs the pyglossary command and converts the file to CSV
func RunPyGlossary(inputFile string) (string, error) {

 	outputFile := ChangePathExt(inputFile, ".csv")
	LogInfo("pyglossary input:", inputFile)
	cmd := exec.Command("pyglossary", inputFile, outputFile,  "--write-format=Csv")

	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	if err != nil {
		LogError("pyglossary error:", out.String())
		return outputFile, fmt.Errorf("pyglossary error: %v, output: %s", err, out.String())
	}

	LogDebug("pyglossary output: ", out.String())
	return  outputFile, nil
}
