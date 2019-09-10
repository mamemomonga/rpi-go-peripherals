// bme280 BME280 温湿度・気圧センサ
//
// VDD -> 3.3V
// GND -> GND
// CSB -> NC(未接続)
// SDI -> SDA
// SDO -> アドレス設定: LOWで0x76(デフォルト) / HIGHで0x77 ここではデフォルトなのでGND
// SCK -> SDL
//
package bme280

import (
	"strconv"
	"strings"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
)

const (
	AddrLow  byte = 0x76 // SDOをLowに接続した場合のアドレス
	AddrHigh byte = 0x77 // SDOをHighに接続した場合のアドレス
)

type BME280 struct {
	dev         *bmxx80.Dev // BME280
	temperature float64     // 温度(摂氏)
	humidity    float64     // 湿度(%rH)
	pressure    float64     // 気圧(hPa)
}

type THIT struct {
	Value  float64 // 不快指数
	Number int     // 番号
	FeelJa string  // 日本語
	FeelEn string  // 英語
}

// New コンストラクタ
func New(bus i2c.BusCloser, addr byte) (t *BME280, err error) {
	t = new(BME280)
	t.dev, err = bmxx80.NewI2C(bus, uint16(addr), &bmxx80.DefaultOpts)
	if err != nil {
		return t, err
	}
	return t, nil
}

// Sense 取得
func (t *BME280) Sense() error {

	e := physic.Env{}
	if err := t.dev.Sense(&e); err != nil {
		return err
	}

	var err error = nil
	t.temperature, err = strconv.ParseFloat(strings.TrimRight(e.Temperature.String(), "°C"), 64)
	if err != nil {
		return err
	}

	t.humidity, err = strconv.ParseFloat(strings.TrimRight(e.Humidity.String(), "%rH"), 64)
	if err != nil {
		return err
	}

	t.pressure, err = strconv.ParseFloat(strings.TrimRight(e.Pressure.String(), "kPa"), 64)
	t.pressure = t.pressure * 10
	if err != nil {
		return err
	}
	return nil
}

// Temperature 温度(摂氏)
func (t *BME280) Temperature() float64 {
	return t.temperature
}

// Humidity 湿度(%rH)
func (t *BME280) Humidity() float64 {
	return t.humidity
}

// Pressure 気圧(hPa)
func (t *BME280) Pressure() float64 {
	return t.pressure
}

// Atm 大気圧(atm)
func (t *BME280) Atm() float64 {
	return t.pressure / 1013.25
}

// THI 不快指数
func (t *BME280) THI() THIT {
	thi := THIT{}

	// 不快指数の計算 https://keisan.casio.jp/exec/system/1202883065
	thi.Value = 0.81*t.temperature + 0.01*t.humidity*(0.99*t.temperature-14.3) + 46.3
	switch {
	case thi.Value <= 55:
		thi.Number = 1
		thi.FeelJa = "寒い"
		thi.FeelEn = "cold"

	case thi.Value > 55 && thi.Value <= 60:
		thi.Number = 2
		thi.FeelJa = "肌寒い"
		thi.FeelEn = "chilly"

	case thi.Value > 60 && thi.Value <= 65:
		thi.Number = 3
		thi.FeelJa = "何も感じない"
		thi.FeelEn = "no feel anything"

	case thi.Value > 65 && thi.Value <= 70:
		thi.Number = 4
		thi.FeelJa = "快適"
		thi.FeelEn = "comfortable"

	case thi.Value > 70 && thi.Value <= 75:
		thi.Number = 5
		thi.FeelJa = "暑くない"
		thi.FeelEn = "not hot"

	case thi.Value > 75 && thi.Value <= 80:
		thi.Number = 6
		thi.FeelJa = "やや暑い"
		thi.FeelEn = "slightly hot"

	case thi.Value > 80 && thi.Value <= 85:
		thi.Number = 7
		thi.FeelJa = "暑くて汗が出る"
		thi.FeelEn = "Hot and sweaty"

	case thi.Value > 85:
		thi.Number = 8
		thi.FeelJa = "暑くてたまらない"
		thi.FeelEn = "hot and irresistible"
	}

	return thi
}
