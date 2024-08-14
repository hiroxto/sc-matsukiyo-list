package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// RawStore APIから取得した生の店舗情報
type RawStore struct {
	Id                int         `json:"id"`
	Name              string      `json:"name"`
	Icon              int         `json:"icon"`
	BusinessCompanyId string      `json:"business_company_id"`
	PostalCode        string      `json:"postal_code"`
	Address           string      `json:"address"`
	Latitude          float64     `json:"latitude"`
	Longitude         float64     `json:"longitude"`
	ClosedDay         string      `json:"closed_day"`
	Comment           interface{} `json:"comment"`
	Url               string      `json:"url"`
	BusinessHours     string      `json:"businesshours"`
	Payments          string      `json:"payments"`
	Products          string      `json:"products"`
	Services          string      `json:"services"`
}

// MixedSlice StoreAttributesで利用されるslice。stringやintが混じっているため，これで表現する。
type MixedSlice = []interface{}

// StoreAttributes StoreAttributesのレスポンス。
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

// InformationAndIcon 店舗情報の名前とアイコンを表す
type InformationAndIcon struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// Store 扱いやすいように加工した店舗情報。
type Store struct {
	Id            int                  `json:"id"`
	Name          string               `json:"name"`
	PostalCode    string               `json:"postal_code"`
	Address       string               `json:"address"`
	Latitude      float64              `json:"latitude"`
	Longitude     float64              `json:"longitude"`
	Url           string               `json:"url"`
	ClosedDay     string               `json:"closed_day"`
	Comment       interface{}          `json:"comment"`
	BusinessHours []InformationAndIcon `json:"businesshours"`
	Services      []InformationAndIcon `json:"services"`
	Products      []InformationAndIcon `json:"products"`
	Payments      []InformationAndIcon `json:"payments"`
}

// getStores 店舗一覧を取得する。
func getStores() ([]RawStore, error) {
	var rawStores []RawStore

	req, err := http.NewRequest("GET", "https://www.matsukiyococokara-online.com/map/s3/json/stores.json", nil)
	if err != nil {
		return nil, err
	}
	// NOTE: ヘッダーをセットしないと取得できない
	req.Header.Set("DNT", "1")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	client := new(http.Client)
	storesResponse, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(storesResponse.Body)
	defer storesResponse.Body.Close()

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &rawStores); err != nil {
		return nil, err
	}

	return rawStores, nil
}

// getStoreAttributes 店舗の属性情報を取得する。
func getStoreAttributes() (StoreAttributes, error) {
	var storeAttr StoreAttributes

	req, err := http.NewRequest("GET", "https://www.matsukiyococokara-online.com/map/s3/json/storeAttributes.json", nil)
	if err != nil {
		return StoreAttributes{}, err
	}
	// NOTE: ヘッダーをセットしないと取得できない
	req.Header.Set("DNT", "1")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	client := new(http.Client)
	storeAttributesResponse, err := client.Do(req)
	if err != nil {
		return StoreAttributes{}, err
	}

	body, err := io.ReadAll(storeAttributesResponse.Body)
	defer storeAttributesResponse.Body.Close()

	if err != nil {
		return StoreAttributes{}, err
	}

	if err := json.Unmarshal(body, &storeAttr); err != nil {
		return StoreAttributes{}, err
	}

	return storeAttr, nil
}

// filterOnlyScRawStores 店舗のスライスからSC店舗のみを抽出する。dカード特約店(クレジットカード/iD)の値を見て判断する。
func filterOnlyScRawStores(rawStores []RawStore) []RawStore {
	scRawStores := make([]RawStore, 0)

	// NOTE: ハードコーディングしているからサービスの増減があると動かなくなる
	r := regexp.MustCompile(`\d{7}0\d{2}`)

	for _, rawStore := range rawStores {
		if r.Match([]byte(rawStore.Services)) {
			scRawStores = append(scRawStores, rawStore)
		}
	}

	return scRawStores
}

// convertAttrToInformation 店舗の属性情報を InformationAndIcon に変換する。
func convertAttrToInformation(bitsString string, attr []MixedSlice) ([]InformationAndIcon, error) {
	infos := make([]InformationAndIcon, 0)

	if len(bitsString) != len(attr) {
		return nil, errors.New("ビット数と属性数が一致していません。")
	}

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

	return infos, nil
}

// convertRawStoreToStore 生の店舗情報を扱いやすい Store に変換する。
func convertRawStoreToStore(rawStore RawStore, attrs StoreAttributes) (Store, error) {
	var store Store

	store.Id = rawStore.Id
	store.Name = rawStore.Name
	store.PostalCode = rawStore.PostalCode
	store.Address = rawStore.Address
	store.Latitude = rawStore.Latitude
	store.Longitude = rawStore.Longitude
	store.Url = rawStore.Url
	store.ClosedDay = rawStore.ClosedDay
	store.Comment = rawStore.Comment

	businessHours, err := convertAttrToInformation(rawStore.BusinessHours, attrs.BusinessHours)
	if err != nil {
		return store, err
	}
	store.BusinessHours = businessHours

	services, err := convertAttrToInformation(rawStore.Services, attrs.Services)
	if err != nil {
		return store, err
	}
	store.Services = services

	products, err := convertAttrToInformation(rawStore.Products, attrs.Products)
	if err != nil {
		return store, err
	}
	store.Products = products

	payments, err := convertAttrToInformation(rawStore.Payments, attrs.Payments)
	if err != nil {
		return store, err
	}
	store.Payments = payments

	return store, nil
}

// convertRawStoresToStores 生の店舗情報のスライスを扱いやすい店舗情報のスライスに変換する。
func convertRawStoresToStores(rawStores []RawStore, attrs StoreAttributes) ([]Store, error) {
	var stores []Store

	for _, rawStore := range rawStores {
		store, err := convertRawStoreToStore(rawStore, attrs)
		if err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	return stores, nil
}
