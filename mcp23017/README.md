# MCP23017

MCP23017はシリアルインターフェイス内臓 16ビットI/Oエキスパンダです

* [データシート 日本語](http://ww1.microchip.com/downloads/jp/DeviceDoc/20001952C_JP.pdf)
* [データシート 英語](http://ww1.microchip.com/downloads/en/devicedoc/20001952c.pdf)

# 配線

* i2cを有効にする
* SDA,SCLを接続する
* A0~A2 はプルダウン
* /RESETはプルアップ
* i2cのプルアップは不要(Raspberry Pi側でプルアップされている)

# サンプルコード

[example\_test.go](example_test.go)

	go test -v -count=1 .

