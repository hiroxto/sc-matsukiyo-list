# sc-matsukiyo-list

SCのマツモトキヨシを抽出するプログラム

## 収集方法

[マツキヨの公式サイトの店舗検索ページ](https://www.matsukiyo.co.jp/map/search)で使用しているJSONデータを使う。
- stores.json
    - 店舗の一覧が返ってくる
    - 決済方法や取り扱い商品，サービスなどは2進数で表されている
- storeAttributes.json
    - 店舗一覧のJSONで2進数で表されている物の詳細が入っている
    - 配列の0番目はid，1番目は決済方法など名前，2番目はアイコンのファイル名

店舗一覧の取り扱いサービスを取得し，サービスの「dカード特約店（クレジットカード/iD）」が無効になっている店舗をSCと見なす。
殆どの場合はこれで判別出来るはず。

これらを自動的に判別し，JSONファイルに書き出す。

## Licence

MIT Licence
