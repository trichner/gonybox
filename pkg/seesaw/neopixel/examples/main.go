package main

import (
	"machine"
	"time"
	"trelligo/pkg/seesaw"
	"trelligo/pkg/seesaw/neopixel"
)

const DefaultNeoTrellisAddress = 0x2E
const neoPixelPin = 3

func main() {
	machine.InitSerial()
	time.Sleep(3 * time.Second)

	// an example to use the NeoPixels on the NeoTrellis board
	log("init i2c")
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{})
	if err != nil {
		fatal(err)
	}

	log("init seesaw")
	ss := seesaw.New(DefaultNeoTrellisAddress, i2c)
	err = ss.Begin()
	if err != nil {
		fatal(err)
	}

	//v, err := ss.ReadVersion()
	//if err != nil {
	//	fatal(err)
	//}
	//log("version: " + strconv.FormatInt(int64(v), 10))

	time.Sleep(100 * time.Millisecond)

	log("init neopixel")
	pix, err := neopixel.New(ss, neoPixelPin, 16, neopixel.PixelTypeGRB)
	if err != nil {
		fatal(err)
	}

	log("start rainbow")
	buf := make([]neopixel.RGBW, 16)
	i := uint8(0)
	for {
		r, g, b := colorWheel(i * 13)
		c := neopixel.RGBW{r, g, b, 0}
		p := uint16(i % 16)
		buf[p] = c
		err = pix.WriteColors(buf)
		if err != nil {
			fatal(err)
		}

		err = pix.ShowPixels()
		if err != nil {
			fatal(err)
		}
		time.Sleep(50 * time.Millisecond)
		i++
	}

}

func colorWheel(p uint8) (r, g, b uint8) {
	p = 255 - p
	if p < 85 {
		return 255 - p*3, 0, p * 3
	}
	if p < 170 {
		p -= 85
		return 0, p * 3, 255 - p*3
	}
	p -= 170
	return p * 3, 255 - p*3, 0
}

func fatal(err error) {
	log("fatal: " + err.Error())
	panic(err)
}
func log(s string) {
	machine.Serial.Write([]byte(s + "\r\n"))
}
