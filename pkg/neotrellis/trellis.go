package neotrellis

import (
	"strconv"
	"trelligo/pkg/debug"
	"trelligo/pkg/seesaw"
	"trelligo/pkg/seesaw/keypad"
	"trelligo/pkg/seesaw/neopixel"
)

const DefaultNeotrellisAddress = 0x2E
const neoPixelPin = 3
const neoPixelType = neopixel.NEO_GRB
const rowCount = 4
const columnCount = 4
const keyCount = rowCount * columnCount

func seesawKeyToNeoTrellis(x uint8) uint8 {
	return ((x)/8)*4 + ((x) % 8)
}
func neoTrellisKeyToSeesawKey(x uint8) uint8 {
	return ((x)/4)*8 + ((x) % 4)
}

type Device struct {
	dev    *seesaw.Device
	pix    *neopixel.SeesawNeopixel
	kpd    *keypad.SeesawKeypad
	events []keypad.KeyEvent
}

func New(dev I2C) (*Device, error) {

	ss := seesaw.New(DefaultNeotrellisAddress, dev)
	if err := ss.Begin(); err != nil {
		return nil, err
	}

	pix, err := neopixel.New(ss, neoPixelPin, keyCount, neoPixelType)
	if err != nil {
		return nil, err
	}

	kbd := keypad.New(ss)
	err = kbd.SetKeypadInterrupt(true)
	if err != nil {
		return nil, err
	}

	return &Device{
		dev:    ss,
		pix:    pix,
		kpd:    kbd,
		events: make([]keypad.KeyEvent, 16),
	}, nil
}

// TODO: translate to x/y
func (d *Device) SetPixelColor(offset uint16, r, g, b, w uint8) error {
	return d.pix.SetPixelColor(offset, r, g, b, w)
}

func (d *Device) ShowPixels() error {
	return d.pix.ShowPixels()
}

// ConfigureKeypad enables or disables a key and edge on the keypad module
func (d *Device) ConfigureKeypad(key uint8, edge keypad.Edge, enable bool) error {
	return d.kpd.ConfigureKeypad(neoTrellisKeyToSeesawKey(key), edge, enable)
}

// ProcessKeyEvents reads pending keypad.KeyEvent s from the FIFO and processes them
func (d *Device) ProcessKeyEvents() error {

	n, err := d.kpd.KeyEventCount()
	if err != nil {
		return err
	}

	buf := d.events[:minu8(uint8(cap(d.events)), n)]

	err = d.kpd.Read(buf)
	if err != nil {
		return err
	}
	for _, e := range buf {
		key := seesawKeyToNeoTrellis(e.Key())
		debug.Log("keypress: " + strconv.Itoa(int(key)) + " (" + strconv.Itoa(int(e.Edge())) + ")")
	}
	return nil
}

func minu8(a, b uint8) uint8 {
	if a > b {
		return b
	}
	return a
}
