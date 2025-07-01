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
			return nil, errors.New("error while loading Json as data")
		}
		data[key] = v
	}

	return data, nil
}
