package main

import "log"

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

	stores := convertRawStoresToStores(rawScStores, attrs)

	err = exportToJson(stores, "dist/sc-matsukiyo-list.json")

	if err != nil {
		log.Fatal("JSONエクスポートに失敗", err)
	}
}
