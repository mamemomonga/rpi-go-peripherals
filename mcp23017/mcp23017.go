// mcp23017 ポートマルチプレクサ
//
package mcp23017

import (
	"log"
	"periph.io/x/periph/conn/i2c"
)

const (
	AddrAllLow byte = 0x20 // A0,A1,A2 をすべてGNDに接続すると、アドレスは0x20
)

type MCP23017 struct {
	dev *i2c.Dev
	A   *PortT
	B   *PortT
}

type PortT struct {
	p     *MCP23017
	addr  addrT
	State byte
}

type addrT struct {
	iodir   byte
	gppu    byte
	olat    byte
	gpio    byte
	gpinten byte
	defval  byte
	intcon  byte
}


// S2BV Slice to Bit vector
func S2BV(in []uint8) uint8 {
	var v uint8
	var i uint8
	for i = 0; i < 8; i++ {
		v = v | (in[i] << i)
	}
	return v
}

// BV2S Bit vector to Slice
func BV2S(in uint8) []uint8 {
	v := []uint8{0,0,0,0,0,0,0,0}
	var i uint8
	for i = 0; i < 8; i++ {
		if (in & (1 << i)) > 0 {
			v[i] = 1
		} else {
			v[i] = 0
		}
	}
	return v
}

func New(bus i2c.BusCloser, addr byte) *MCP23017 {
	t := new(MCP23017)
	t.dev = &i2c.Dev{Addr: uint16(addr), Bus: bus}
	t.A = newPort("A", t)
	t.B = newPort("B", t)
	return t
}

func (t *MCP23017) Write(addr uint8, data uint8) {
	if _,err := t.dev.Write([]byte{addr,data}); err != nil {
		log.Fatal(err)
	}
}

func (t *MCP23017) Read(addr uint8) byte {
	v := make([]byte, 1,1)
	if err := t.dev.Tx([]byte{addr},v); err != nil {
		log.Fatal(err)
	}
	return v[0]
}

func newPort(port string, parent *MCP23017) *PortT {
	t := new(PortT)
	t.p = parent
	t.State = 0x00

	switch port {
	// ICON.BANK=0 専用
	case "A":
		t.addr = addrT{
			iodir:   0x00, // IODIRA
			gppu:    0x0c, // GPPUA
			olat:    0x14, // OLATA
			gpio:    0x12, // GPIOA
			gpinten: 0x04, // GPINTENA
			defval:  0x06, // DEFVALA
			intcon:  0x08, // INTCONA
		}
	case "B":
		t.addr = addrT{
			iodir:   0x01, // IODIRB
			gppu:    0x0d, // GPPUB
			olat:    0x15, // OLATB
			gpio:    0x13, // GPIOB
			gpinten: 0x05, // GPINTENB
			defval:  0x07, // DEFVALB
			intcon:  0x09, // INTCONB
		}
	}
	return t
}

func (t *PortT) Direction(v ...uint8) {
	t.p.Write(t.addr.iodir, S2BV(v))
}
func (t *PortT) DirectionAllOutput() {
	t.Direction(0,0,0,0,0,0,0,0)
}
func (t *PortT) DirectionAllInput() {
	t.Direction(1,1,1,1,1,1,1,1)
}
func (t *PortT) PullUp(v ...uint8) {
	t.p.Write(t.addr.gppu, S2BV(v))
}
func (t *PortT) PullUpAll() {
	t.PullUp(1,1,1,1,1,1,1,1)
}
func (t *PortT) Apply() {
	t.p.Write(t.addr.olat, t.State)
}
func (t *PortT) Fetch() {
	t.State = t.p.Read(t.addr.gpio)
}
func (t *PortT) SetAllLow(a bool) {
	t.State = 0x00
	if a {
		t.Apply()
	}
}
func (t *PortT) SetAllHigh(a bool) {
	t.State = 0xFF
	if a {
		t.Apply()
	}
}
func (t *PortT) Set(p byte, v bool, a bool) {
	if v {
		t.State = t.State |  ( 1 << p )
	} else {
		t.State = t.State &^ ( 1 << p )
	}
	if a {
		t.Apply()
	}
}

func (t *PortT) Get(p byte) bool {
	if t.State & ( 1 << p ) > 0 {
		return true
	} else {
		return false
	}
}

func (t *PortT) InitInterrupt() {
	// 以下、すべてのピンに対して設定する
	// INTnピンはデフォルトはアクティブLow
	// IOCON.INTPOLで変更できる

	// 割り込みを有効にする
	t.p.Write(t.addr.gpinten, 0xFF)

	// 前の状態と比較する
	t.p.Write(t.addr.intcon, 0x00)

}

