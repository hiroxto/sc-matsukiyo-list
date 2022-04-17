package main

import (
	"encoding/json"
	"os"
)

func exportToJson(stores []Store, filePath string) error {
	jsonData, err := json.Marshal(stores)

	if err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	_, err = f.Write(jsonData)

	if err != nil {
		return err
	}

	return nil
}
