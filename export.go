package main

import (
	"encoding/json"
	"os"
)

// exportToJson は店舗情報をJSONにエクスポートする
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
