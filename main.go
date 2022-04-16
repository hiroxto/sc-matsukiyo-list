package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
)

type Store struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Icon              int    `json:"icon"`
	BusinessCompanyId string `json:"business_company_id"`
	PostalCode        string `json:"postal_code"`
	Address           string `json:"address"`
	Url               string `json:"url"`
	Payments          string `json:"payments"`
	Products          string `json:"products"`
	Services          string `json:"services"`
}

type MixedSlice = []interface{}

type StoreAttributes struct {
	Config struct {
		IconPath string `json:"iconPath"`
	} `json:"config"`
	Icon              []MixedSlice `json:"icon"`
	Businesshours     []MixedSlice `json:"businesshours"`
	Services          []MixedSlice `json:"services"`
	Products          []MixedSlice `json:"products"`
	Payments          []MixedSlice `json:"payments"`
	BusinessCompanyId []MixedSlice `json:"business_company_id"`
}

func getStores() ([]Store, error) {
	var stores []Store

	storesResponse, err := http.Get("https://www.matsukiyo.co.jp/map/s3/json/stores.json")
	defer storesResponse.Body.Close()
	body, err := io.ReadAll(storesResponse.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &stores); err != nil {
		return nil, err
	}

	return stores, nil
}

func getStoreAttributes() (StoreAttributes, error) {
	var storeAttr StoreAttributes

	storeAttributesResponse, err := http.Get("https://www.matsukiyo.co.jp/map/s3/json/storeAttributes.json")
	defer storeAttributesResponse.Body.Close()
	body, err := io.ReadAll(storeAttributesResponse.Body)

	if err != nil {
		return StoreAttributes{}, err
	}

	if err := json.Unmarshal(body, &storeAttr); err != nil {
		return StoreAttributes{}, err
	}

	return storeAttr, nil
}

func filterOnlyScStores(stores []Store) []Store {
	var scStores []Store

	r := regexp.MustCompile(`\d{8}0\d{2}`)

	for _, store := range stores {
		if r.Match([]byte(store.Services)) {
			scStores = append(scStores, store)
		}
	}

	return scStores
}

func main() {
	stores, err := getStores()
	if err != nil {
		log.Fatal("店舗一覧の取得に失敗")
	}
	attrs, err := getStoreAttributes()
	if err != nil {
		log.Fatal("属性の取得に失敗")
	}

	scStores := filterOnlyScStores(stores)
}
