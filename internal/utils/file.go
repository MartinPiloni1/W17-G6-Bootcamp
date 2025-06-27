package utils

import (
	"fmt"
	"os"
)

// based on a file, transform it to a json
func ReadFile(FilePath string) ([]byte, error) {
	data, err := os.ReadFile(FilePath)
	if err != nil {
		return nil, fmt.Errorf("Error at reading file -> %s", err)
	}
	return data, nil
}

// write the json data into the file
func WriteFile(FilePath string, data []byte) error {
	err := os.WriteFile(FilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("Error at writing file -> %s", err)
	}
	return nil
}
