package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	//	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
)

// BME280 温湿度・気圧センサ
// AE-BME280を使用した http://akizukidenshi.com/catalog/g/gK-09421/
// J3 のみブリッジすると I2C が有効になる(CSB=HIGH)
// プルアップは使用しない
// VDD -> 3.3V
// GND -> GND
// CSB -> NC(未接続)
// SDI -> SDA
// SDO -> アドレス設定: LOWで0x76(デフォルト) / HIGHで0x77 ここではデフォルトなのでGND
// SCK -> SDL

func i2cBME280() {
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

	d, err := bmxx80.NewI2C(b, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}
	e := physic.Env{}
	if err := d.Sense(&e); err != nil {
		log.Fatal(err)
	}

	t := time.NewTicker(30 * time.Second)

	for {
		log.Printf("%8s %9s %10s\n", e.Temperature, e.Humidity, e.Pressure)

		// 数値に再変換
		temperature, _ := strconv.ParseFloat(strings.TrimRight(e.Temperature.String(), "°C"), 64)
		humidity, _ := strconv.ParseFloat(strings.TrimRight(e.Humidity.String(), "%rH"), 64)
		pressure, _ := strconv.ParseFloat(strings.TrimRight(e.Pressure.String(), "kPa"), 64)

		// 不快指数の計算 https://keisan.casio.jp/exec/system/1202883065
		di := 0.81*temperature + 0.01*humidity*(0.99*temperature-14.3) + 46.3
		dit := ""
		switch {
		case di <= 55:
			dit = "寒い"
		case di > 55 && di <= 60:
			dit = "肌寒い"
		case di > 60 && di <= 65:
			dit = "何も感じない"
		case di > 65 && di <= 70:
			dit = "快適"
		case di > 70 && di <= 75:
			dit = "暑くない"
		case di > 75 && di <= 80:
			dit = "やや暑い"
		case di > 80 && di <= 85:
			dit = "暑くて汗が出る"
		case di > 85:
			dit = "暑くてたまらない"
		}

		fmt.Printf("  気温: 摂氏 %2.2f 度\n", temperature)
		fmt.Printf("  湿度: %3.2f パーセント\n", humidity)
		fmt.Printf("  気圧: %4.4f ヘクトパスカル (%2.4f 気圧) \n", pressure*10, pressure/101.325)
		fmt.Printf("  不快指数: %2.2f (%s) \n", di, dit)

		<-t.C
	}
}
