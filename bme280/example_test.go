package bme280_test



import (
	"os"
	"fmt"
	"github.com/mamemomonga/rpi-go-peripherals/bme280"
	"log"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

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
	bm, err := bme280.New(b, bme280.AddrLow)
	if err != nil {
		log.Fatal(err)
	}

	// 検出
	err = bm.Sense()
	if err != nil {
		log.Fatal(err)
	}

	// 結果
	thi := bm.THI()
	fmt.Fprintf(os.Stderr,"  気温: 摂氏 %2.2f 度\n", bm.Temperature())
	fmt.Fprintf(os.Stderr,"  湿度: %3.2f パーセント\n", bm.Humidity())
	fmt.Fprintf(os.Stderr,"  気圧: %4.4f ヘクトパスカル (%2.4f 気圧)\n", bm.Pressure(), bm.Atm())
	fmt.Fprintf(os.Stderr,"  不快指数: %2.2f %s(%s)\n", thi.Value, thi.FeelJa, thi.FeelEn)

	// Output: 
}
