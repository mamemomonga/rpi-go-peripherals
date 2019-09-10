package mcp23017_test

import (
	"log"
	"time"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/mamemomonga/rpi-go-peripherals/mcp23017"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"

	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/gpio"
)

// boolを文字列にする
func bool2str(b bool) string {
	if b {
		return "HIGH"
	}
	return "LOW"
}

// 割り込み処理
func btnInterruptLoop(intpin gpio.PinIO, mx *mcp23017.MCP23017) {
	fmt.Fprintln(os.Stderr,"START btnInterruptLoop")

	// mx.B.Get(0) という形でも呼び出せるが
	// mcp23017.NewPinという形でピン番号も含んだ形でインスタンス化し
	// Get() だけで取得できるようにすることができる

	// ポートB ピン0 GPB0
	button0 := mcp23017.NewPin(mx.B, 0)

	// ポートB ピン1 GPB1
	button1 := mcp23017.NewPin(mx.B, 1)

	// ポートB ピン2 GPB2
	button2 := mcp23017.NewPin(mx.B, 2)

	// ポートB ピン3 GPB3
	button3 := mcp23017.NewPin(mx.B, 3)

	for {
		// 割り込みを待つ
		intpin.WaitForEdge(-1)

		// ポートBの値を取得
		mx.B.Fetch()

		// GPB0, GPB1, GPB2, GPB3 の値を表示
		fmt.Fprintf(os.Stderr, " [GPB0 %4s] [GPB1 %4s] [GPB2 %4s] [GPB3 %4s]\n",
			bool2str(button0.Get()),
			bool2str(button1.Get()),
			bool2str(button2.Get()),
			bool2str(button3.Get()),
		)
	}
}

// LED処理
func ledLoop(mx *mcp23017.MCP23017) {
	fmt.Fprintln(os.Stderr,"START ledLoop")
	for {
		for i:=uint8(0); i<8; i++ {
			// ポートAのピン i の LEDを点灯し、反映する
			mx.A.Set(i,true).Apply()
			time.Sleep( 100 * time.Millisecond )
			// ポートAのピン i の LEDを消灯し、反映しない
			// (次のポートAの反映の際まで遅延)
			mx.A.Set(i,false)
		}
	}
}

// メイン
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
	mx.A.SetAllHigh().Apply()

	time.Sleep( time.Second )

	// ポートAをすべてLowにし、反映する
	mx.A.SetAllLow().Apply()

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

	fmt.Fprintln(os.Stderr," *** START ***")


	// 割り込み処理
	go btnInterruptLoop(intpin,mx)

	// LED点滅処理
	go ledLoop(mx)

	// シグナル割り込みの設定
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT)

	// シグナルがくるまでここでブロック
	<-quit

	fmt.Fprintln(os.Stderr," *** STOP ***")

	// ポートAをすべてLowにし、反映する
	mx.A.SetAllLow().Apply()

	fmt.Println("Done")
	// testingのExample機能は標準出力を
	// 以下の Output: の後と比較するためエラー出力に結果を出している
	// https://golang.org/pkg/testing/#hdr-Examples

	// Output: Done
}

