package main

import (
	"encoding/json"
	"io"
	"log"
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

type InformationAndIcon struct {
	Name string
	Icon string
}

type Store struct {
	BusinessHours []InformationAndIcon // 未実装
	Services      []InformationAndIcon
	Products      []InformationAndIcon
	Payments      []InformationAndIcon
	RawStore      RawStore
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
	var scRawStores []RawStore

	r := regexp.MustCompile(`\d{8}0\d{2}`)

	for _, rawStore := range rawStores {
		if r.Match([]byte(rawStore.Services)) {
			scRawStores = append(scRawStores, rawStore)
		}
	}

	return scRawStores
}

func convertAttrToInformation(bitsString string, attr []MixedSlice) []InformationAndIcon {
	var infos []InformationAndIcon

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

	store.Services = convertAttrToInformation(rawStore.Services, attrs.Services)
	store.Products = convertAttrToInformation(rawStore.Products, attrs.Products)
	store.Payments = convertAttrToInformation(rawStore.Payments, attrs.Payments)

	store.RawStore = rawStore

	return store
}

func convertRawStoresToStores(rawStores []RawStore, attrs StoreAttributes) []Store {
	var stores []Store

	for _, rawStore := range rawStores {
		stores = append(stores, convertRawStoreToStore(rawStore, attrs))
	}

	return stores
}

func main() {
	rawStores, err := getStores()
	if err != nil {
		log.Fatal("店舗一覧の取得に失敗")
	}
	attrs, err := getStoreAttributes()
	if err != nil {
		log.Fatal("属性の取得に失敗")
	}

	rawScStores := filterOnlyScRawStores(rawStores)

	// ToDo: 未実装
	_ = convertRawStoresToStores(rawScStores, attrs)
}
