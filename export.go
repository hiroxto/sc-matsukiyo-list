package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// fileExists はファイルの存在を確認する
func fileExists(name string) bool {
	_, err := os.Stat(name)

	return !os.IsNotExist(err)
}

// exportToJson は店舗情報をJSONにエクスポートする
func exportToJson(stores []Store, filePath string) error {
	jsonData, err := json.Marshal(stores)
	if err != nil {
		return err
	}

	dirname := filepath.Dir(filePath)
	if !fileExists(dirname) {
		if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
			return err
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	if _, err = f.Write(jsonData); err != nil {
		return err
	}

	return nil
}
