// gpioBlink1 LED点滅
//
// 配線:発光ダイオードをGNDとGPIO26の間に1KΩの抵抗を介して接続する
//
package main

import (
	"log"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	"time"
)

// 発光ダイオードをGNDとGPIO26の間に1KΩの抵抗を介して接続

func main() {

	// 初期化
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// GPIO26を使用する
	pin := gpioreg.ByName("GPIO26")

	// 最初はHighにする
	pin.Out(gpio.High)

	// 1秒寝る
	time.Sleep(time.Second)

	// 点滅
	t := time.NewTicker(100 * time.Millisecond)

	log.Println("running")
	for l := gpio.Low; ; l = !l {
		pin.Out(l)
		<-t.C
	}
}
