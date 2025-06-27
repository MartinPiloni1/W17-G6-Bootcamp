package utils

import (
	"encoding/json"
	"os"
)

func Read[T Identifiable](filePath string) (map[int]T, error) {
	dataAsBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var slicedData []T
	err = json.Unmarshal(dataAsBytes, &slicedData)
	if err != nil {
		return nil, err
	}

	data := make(map[int]T)
	for _, e := range slicedData {
		data[e.GetID()] = e
	}

	return data, nil
}

func Write[T any](filePath string, newData T) error {
	dataAsBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var slicedData []T
	err = json.Unmarshal(dataAsBytes, &slicedData)
	if err != nil {
		return err
	}

	slicedData = append(slicedData, newData)
	newDataAsBytes, err := json.MarshalIndent(slicedData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, newDataAsBytes, 0644)
	return err
}

func GetNextID[T Identifiable](filePath string) (int, error) {
	dataAsBytes, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	var slicedData []T
	err = json.Unmarshal(dataAsBytes, &slicedData)
	if err != nil {
		return 0, err
	}

	maxId := 0
	for _, v := range slicedData {
		maxId = max(maxId, v.GetID())
	}
	return maxId + 1, nil
}
