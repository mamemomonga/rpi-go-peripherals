# periph-examples

ここは、このリポジトリのライブラリにしない、[periph](https://periph.io) ライブラリを単独で利用するサンプルです。

* [gpioBlink1 LED点滅 (bcm283x 直接アクセス)](gpioBlink1/main.go)
* [gpioBlink2 LED点滅](gpioBlink2/main.go)
* [gpioButton gpioButton ボタン反応 エッジ待機](gpioButton/main.go)
* [i2cBME280 BME280 温湿度・気圧センサ](i2cBME280/main.go)
* [i2cMCP23017 i2cの基本サンプル (MCP23017ポートマルチプレクサ)](i2cMCP23017/main.go)

# 実行方法

それぞれのディレクトリで

	go run .

を実行してください。CTRL-C で終了です。

