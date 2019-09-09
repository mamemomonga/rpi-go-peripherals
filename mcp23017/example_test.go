package mcp23017_test

import (
	"github.com/mamemomonga/rpi-go-peripherals/mcp23017"
	"os"
	"log"
	"time"
	"fmt"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"

	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/gpio"
)

func bool2str(b bool) string {
	if b {
		return "HIGH"
	}
	return "LOW"
}


func Example() {

	// 初期化
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// I2Cポートを開く、空欄時は最初にみつかったポート
	// ポート名は i2c-list でわかる
	b, err := i2creg.Open("I2C1")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	// 開始
	// A0,A1,A2 をすべてGNDに接続すると、アドレスは0x20
	mx := mcp23017.New(b, mcp23017.AddrAllLow)

	// ポートAをすべて出力にする
	mx.A.DirectionAllOutput()

	// ポートAをすべてHighにし、反映する
	mx.A.SetAllHigh(true)

	time.Sleep( time.Second )

	// ポートAをすべてLowにし、反映する
	mx.A.SetAllLow(true)

	// ポートBをすべて入力にする
	mx.B.DirectionAllInput()

	// ポートBをすべてプルアップする
	mx.B.PullUpAll()

	// ポートBの割り込みを有効にする
	mx.B.InitInterrupt()

	// ポートBの状態を取得する
	// これでピンの「前の値」が設定される
	mx.B.Fetch()

	// MCP23017のINTB と Raspberry PiのGPIO13 を接続する
	// GPIO13を割り込み入力にする
	intpin := gpioreg.ByName("GPIO13")

	// INTはアクティブLowであるので、入力・フロート(プルアップダウンしない)・エッジ下降
	if err := intpin.In(gpio.Float, gpio.FallingEdge); err != nil {
		log.Fatal(err)
	}

	// 割り込み処理
	go func() {
		for {
			// 割り込みを待つ
			intpin.WaitForEdge(-1)

			// ポートBの値を取得
			mx.B.Fetch()

			// GPB0, GPB1, GPB2, GPB3 の値を表示
			fmt.Fprintf(os.Stderr, " [GPB0 %4s] [GPB1 %4s] [GPB2 %4s] [GPB3 %4s]\n",
				bool2str(mx.B.Get(0)),
				bool2str(mx.B.Get(1)),
				bool2str(mx.B.Get(2)),
				bool2str(mx.B.Get(3)),
			)
		}
	}()

	// LED点滅処理
	for {
		for i:=uint8(0); i<8; i++ {
			// ポートAのピン i の LEDを点灯し、反映する
			mx.A.Set(i,true,true)
			time.Sleep( 100 * time.Millisecond )
			// ポートAのピン i の LEDを消灯し、反映しない
			// (次のポートAの書込の際に反映される)
			mx.A.Set(i,false,false)

		}
	}

	// testingのExample機能は標準出力を
	// 以下の Output: の後と比較するためエラー出力に結果を出している
	// https://golang.org/pkg/testing/#hdr-Examples

	// Output: 
}

