package utils

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

func Read[T any](filePath string) (map[int]T, error) {
	dataAsBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tmpData map[string]T
	err = json.Unmarshal(dataAsBytes, &tmpData)
	if err != nil {
		return nil, err
	}

	data := make(map[int]T)
	for k, v := range tmpData {
		key, err := strconv.Atoi(k)
		if err != nil {
			return nil, errors.New("Error while loading Json as data")
		}
		data[key] = v
	}

	return data, nil
}

// Escribe un map[int]T al filePath codificado como JSON (keys string)
func Write[T any](filePath string, data map[int]T) error {
	// Convertir de map[int]T → map[string]T para JSON
	tmpData := make(map[string]T, len(data))
	for k, v := range data {
		tmpData[strconv.Itoa(k)] = v
	}

	// Marshal a JSON
	dataAsBytes, err := json.MarshalIndent(tmpData, "", "  ")
	if err != nil {
		return err
	}

	// Escribir al archivo (permiso 0644: rw-r--r--)
	return os.WriteFile(filePath, dataAsBytes, 0644)
}
