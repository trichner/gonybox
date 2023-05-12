package neotrellis

import (
	"fmt"
	"trelligo/pkg/seesaw"
	"trelligo/pkg/seesaw/keypad"
	"trelligo/pkg/seesaw/neopixel"
)

// DefaultNeoTrellisAddress is the i2c address of a NeoTrellis board without any changes to its address
// The address can be changed with solder bridges, allowing daisy-chaining multiple NeoTrellis boards.
const DefaultNeoTrellisAddress = 0x2E
const neoPixelPin = 3
const neoPixelType = neopixel.PixelTypeGRB
const yCount = 4
const xCount = 4
const keyCount = yCount * xCount

type RGB struct {
	R, G, B uint8
}

type Device struct {
	dev        *seesaw.Device
	pix        *neopixel.Device
	kpd        *keypad.SeesawKeypad
	events     []keypad.KeyEvent
	keyHandler func(x, y uint8, edge keypad.Edge) error
}

func New(dev I2C, addr uint16) (*Device, error) {

	if addr == 0 {
		addr = DefaultNeoTrellisAddress
	}

	ss := seesaw.New(addr, dev)
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

// SetPixelColor sets the color of a pixel at position x/y
//
// Note: ShowPixels MUST be called to actually show the updated color.
func (d *Device) SetPixelColor(x, y uint8, color RGB) error {
	// `w` is always 0, the NeoTrellis only has RGB NeoPixels
	return d.pix.WriteColorAtOffset(newXy(x, y).neoPixel(), neopixel.RGBW{
		R: color.R,
		G: color.G,
		B: color.B,
	})
}

// WriteColors writes the color for multiple pixels at once. At most all 16 LEDs can
// be updated.
// Note: ShowPixels MUST be called to actually show the updated colors.
func (d *Device) WriteColors(colors []RGB) error {
	if len(colors) > keyCount {
		return fmt.Errorf("too many colors: %d > %d", len(colors), keyCount)
	}

	buf := make([]neopixel.RGBW, len(colors))
	for i, c := range colors {
		buf[i] = neopixel.RGBW{
			R: c.R,
			G: c.G,
			B: c.B,
		}
	}
	return d.pix.WriteColors(buf)
}

// ShowPixels instructs the NeoPixel buffer to update and display the set colors
func (d *Device) ShowPixels() error {
	return d.pix.ShowPixels()
}

// ConfigureKeypad enables or disables a key and edge on the keypad module. Events can be handled by setting a handler
// with SetKeyHandleFunc.
func (d *Device) ConfigureKeypad(x, y uint8, edge keypad.Edge, enable bool) error {
	return d.kpd.ConfigureKeypad(newXy(x, y).seesawKey(), edge, enable)
}

// SetKeyHandleFunc sets a callback for key events
//
// Note: In order for the handler to be called, the keypads MUST be configured via ConfigureKeypad and the events
// MUST be processed via ProcessKeyEvents.
func (d *Device) SetKeyHandleFunc(handler func(x, y uint8, e keypad.Edge) error) {
	d.keyHandler = handler
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

	if d.keyHandler == nil {
		return nil
	}

	for _, e := range buf {
		p := newXyFromSeesawKey(e.Key())
		err := d.keyHandler(p.X(), p.Y(), e.Edge())
		if err != nil {
			return err
		}
	}
	return nil
}

func minu8(a, b uint8) uint8 {
	if a > b {
		return b
	}
	return a
}
