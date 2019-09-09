# BME280

気温・湿度・気圧センサ BME280 を利用して、日本向けの単位と、関連情報を取得します

* [BME280](https://www.bosch-sensortec.com/bst/products/all_products/bme280)
* [データシート](https://ae-bst.resource.bosch.com/media/_tech/media/datasheets/BST-BME280-DS002.pdf)

# 配線

* [AE-BME280](http://akizukidenshi.com/catalog/g/gK-09421/)を使用 
* J3 のみブリッジすると I2C が有効になる(CSB=HIGH)
* プルアップは使用しない
* VDD -> 3.3V
* GND -> GND
* CSB -> NC(未接続)
* SDI -> SDA
* SDO -> アドレス設定: LOWで0x76(デフォルト) / HIGHで0x77 ここではデフォルトなのでGND
* SCK -> SDL

# サンプルコード

[example_test.go](example_test.go)

[testingのExample](https://golang.org/pkg/testing/#hdr-Examples)は標準出力を以下のOutput: の後と較するためエラー出力に結果を出している

実行

	go test -v -count=1 .

