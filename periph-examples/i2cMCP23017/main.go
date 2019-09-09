// ic2MCP23017 i2cの基本サンプル (MCP23017ポートマルチプレクサ)
// 
// 配線とi2c* コマンドでの動作確認はこちら
// https://gist.github.com/mamemomonga/bb915ea66904605598a9331cdbb4ac18
//
package main

import (
	"log"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
	"time"
)

func main() {
	// 初期化
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// I2Cポートを開く、省略時は最初にみつかったポート
	// ポート名は i2c-list でわかる
	b, err := i2creg.Open("I2C1")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	// アドレスの指定
	// MCP23017はADDRピンをすべてGNDに落とすと
	// 0x20となる
	d := &i2c.Dev{Addr: 0x20, Bus: b}

	// データアドレスはICON.BANKで変わる
	// ここではデフォルトの ICON.BANK=0 とする

	// 書込 IODIRAをすべて出力
	if _, err := d.Write([]byte{0x00, 0x00}); err != nil {
		log.Fatal(err)
	}

	// 書込 OLATAをすべてHIGHにする
	if _, err := d.Write([]byte{0x14, 0xFF}); err != nil {
		log.Fatal(err)
	}

	// 1秒寝る
	time.Sleep(time.Second)

	// 書込 OLATAをすべてLOWにする
	if _, err := d.Write([]byte{0x14, 0x00}); err != nil {
		log.Fatal(err)
	}

	// 書込 IODIRBをすべて入力
	if _, err := d.Write([]byte{0x01, 0xFF}); err != nil {
		log.Fatal(err)
	}

	// 読書 GPIOBの値を読む
	read := make([]byte, 1)
	if err := d.Tx([]byte{0x13}, read); err != nil {
		log.Fatal(err)
	}
	log.Printf("GPIOB: 0x%02x", read[0])

}
