package utils

import (
	"encoding/json"
	"os"
)

func MapToSlice[T any](data map[int]T) []T {
	var slicedData []T
	for _, v := range data {
		slicedData = append(slicedData, v)
	}
	return slicedData
}

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

func Write[T any](filePath string, data map[int]T) error {
	slicedData := MapToSlice(data)

	dataAsBytes, err := json.MarshalIndent(slicedData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, dataAsBytes, 0644)
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
