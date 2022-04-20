package main

import "log"

func main() {
	rawStores, err := getStores()
	if err != nil {
		log.Fatal("店舗一覧の取得に失敗", err)
	}
	attrs, err := getStoreAttributes()
	if err != nil {
		log.Fatal("属性の取得に失敗", err)
	}

	rawScStores := filterOnlyScRawStores(rawStores)

	stores, err := convertRawStoresToStores(rawScStores, attrs)
	if err != nil {
		log.Fatal("店舗情報の変換に失敗", err)
	}

	err = exportToJson(stores, "dist/sc-matsukiyo-list.json")

	if err != nil {
		log.Fatal("JSONエクスポートに失敗", err)
	}
}
