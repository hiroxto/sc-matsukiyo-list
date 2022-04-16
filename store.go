package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type RawStore struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Icon              int    `json:"icon"`
	BusinessCompanyId string `json:"business_company_id"`
	PostalCode        string `json:"postal_code"`
	Address           string `json:"address"`
	Url               string `json:"url"`
	BusinessHours     string `json:"businesshours"`
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
	BusinessHours     []MixedSlice `json:"businesshours"`
	Services          []MixedSlice `json:"services"`
	Products          []MixedSlice `json:"products"`
	Payments          []MixedSlice `json:"payments"`
	BusinessCompanyId []MixedSlice `json:"business_company_id"`
}

type InformationAndIcon struct {
	Name string
	Icon string
}

type Store struct {
	Id            int                  `json:"id"`
	Name          string               `json:"name"`
	PostalCode    string               `json:"postal_code"`
	Address       string               `json:"address"`
	Url           string               `json:"url"`
	BusinessHours []InformationAndIcon `json:"businesshours"`
	Services      []InformationAndIcon `json:"services"`
	Products      []InformationAndIcon `json:"products"`
	Payments      []InformationAndIcon `json:"payments"`
}

func getStores() ([]RawStore, error) {
	var rawStores []RawStore

	storesResponse, err := http.Get("https://www.matsukiyo.co.jp/map/s3/json/stores.json")
	defer storesResponse.Body.Close()
	body, err := io.ReadAll(storesResponse.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &rawStores); err != nil {
		return nil, err
	}

	return rawStores, nil
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

func filterOnlyScRawStores(rawStores []RawStore) []RawStore {
	scRawStores := make([]RawStore, 0)

	r := regexp.MustCompile(`\d{8}0\d{2}`)

	for _, rawStore := range rawStores {
		if r.Match([]byte(rawStore.Services)) {
			scRawStores = append(scRawStores, rawStore)
		}
	}

	return scRawStores
}

func convertAttrToInformation(bitsString string, attr []MixedSlice) []InformationAndIcon {
	infos := make([]InformationAndIcon, 0)

	bits := strings.Split(bitsString, "")
	for index, bit := range bits {
		if bit == "1" {
			info := InformationAndIcon{
				Name: attr[index][1].(string),
				Icon: attr[index][2].(string),
			}
			infos = append(infos, info)
		}
	}

	return infos
}

func convertRawStoreToStore(rawStore RawStore, attrs StoreAttributes) Store {
	var store Store

	store.Id = rawStore.Id
	store.Name = rawStore.Name
	store.PostalCode = rawStore.PostalCode
	store.Address = rawStore.Address
	store.Url = rawStore.Url
	store.BusinessHours = convertAttrToInformation(rawStore.BusinessHours, attrs.BusinessHours)
	store.Services = convertAttrToInformation(rawStore.Services, attrs.Services)
	store.Products = convertAttrToInformation(rawStore.Products, attrs.Products)
	store.Payments = convertAttrToInformation(rawStore.Payments, attrs.Payments)

	return store
}

func convertRawStoresToStores(rawStores []RawStore, attrs StoreAttributes) []Store {
	var stores []Store

	for _, rawStore := range rawStores {
		stores = append(stores, convertRawStoreToStore(rawStore, attrs))
	}

	return stores
}
