package main

import (
    "time"
	"log"
    "periph.io/x/periph/conn/gpio"
    "periph.io/x/periph/host"
    "periph.io/x/periph/host/bcm283x"
)

// 発光ダイオードをGNDとGPIO26の間に1KΩの抵抗を介して接続

func gpioBlink1() {

	// 初期化
	if _,err := host.Init(); err != nil {
		log.Fatal(err)
	}

    t := time.NewTicker(100 * time.Millisecond)

    for l := gpio.Low; ; l = !l {
		// 低レベルのbcm283xを直接呼び出す
		// gpioregをつかったほうがいろいろと便利
        bcm283x.GPIO26.Out(l)
        <-t.C
    }
}

