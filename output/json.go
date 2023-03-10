package output

import (
	"encoding/json"
	"fmt"
	"os"
	"yoda/y_types"
)

func CreateJSON(p *y_types.PackageInfo, outputName string) error {

	// Marshal the PackageInfo struct into a pretty-printed JSON string.
	jsonData, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling PackageInfo: %v", err)
	}
	// Check if the file already exists.
	if _, err := os.Stat(outputName); !os.IsNotExist(err) {
		// Delete the file if it exists.
		if err := os.Remove(outputName); err != nil {
			return fmt.Errorf("error deleting file: %v", err)
		}
	}

	// Create a new file with the same name.
	file, err := os.Create(outputName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write the compacted JSON data to the file.
	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("error writing JSON data to file: %v", err)
	}

	return nil
}
