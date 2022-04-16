package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func exportToJson(stores []Store, filePath string) error {
	jsonData, err := json.Marshal(stores)

	if err != nil {
		return err
	}

	fmt.Printf("%s\n", jsonData)

	f, err := os.Create(filePath)
	_, err = f.Write(jsonData)

	if err != nil {
		return err
	}

	return nil
}
