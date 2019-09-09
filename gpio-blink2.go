package main

import (
    "time"
	"log"
    "periph.io/x/periph/conn/gpio"
    "periph.io/x/periph/conn/gpio/gpioreg"
    "periph.io/x/periph/host"
)

// 発光ダイオードをGNDとGPIO26の間に1KΩの抵抗を介して接続

func gpioBlink2() {

	// 初期化
	if _,err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// GPIO26を使用する
	pin := gpioreg.ByName("GPIO26")

	// 最初はHighにする
	pin.Out(gpio.High)

	// 1秒寝る
	time.Sleep( time.Second )

	// 点滅
    t := time.NewTicker(100 * time.Millisecond)
    for l := gpio.Low; ; l = !l {
        pin.Out(l)
        <-t.C
    }
}

