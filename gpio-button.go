package main

import (
	"log"
	"time"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

// GPIOはプルアップされている
// ボタンをGNDとの間に挟む

func gpioButton() {
    host.Init()

	// https://godoc.org/periph.io/x/periph/host/bcm283x
	// https://godoc.org/periph.io/x/periph/conn/gpio
	// https://godoc.org/periph.io/x/periph/conn/gpio/gpioreg

	pin := gpioreg.ByName("GPIO13")

	// 入力・プルアップ・エッジ下降
	if err := pin.In(gpio.PullUp, gpio.FallingEdge); err != nil {
	    log.Fatal(err)
	}

	log.Println("Ready")
	counter := 0
    for {
		// エッジが変わるまでブロック
		// ボタンが押された、もしくは離された
		pin.WaitForEdge(-1)

		// Highならボタンが離された動作なので無視する
		if pin.Read() {
			continue
		}

		counter++
		log.Printf(" %d 回ボタンが押されました\n",counter)

		// ボタンを押されている状態である
		for !pin.Read() {
			time.Sleep(50 * time.Millisecond)
			log.Println("")
		}

    }
}

