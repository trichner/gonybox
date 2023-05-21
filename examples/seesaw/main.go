package main

import (
	"machine"
	"time"
	"trelligo/pkg/seesaw"
)

// seesawWriteDelay the seesaw is quite timing sensitive and times out if not given enough time,
// this is an empirically determined delay that seems to have good results
const seesawWriteDelay = time.Millisecond * 100
const defaultNeoTrellisAddress = 0x2E
const neoPixelPin = 3

// usedPixelCount the number of pixels we're gonna play around
const usedPixelCount = 4

var ledBuffer = make([]byte, 3*usedPixelCount)

func main() {
	machine.InitSerial()

	// give some time to attach to Serial
	time.Sleep(3 * time.Second)

	// an example to use the NeoPixels on the NeoTrellis board
	log("init i2c")
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{})
	if err != nil {
		fatal(err)
	}

	time.Sleep(1000 * time.Millisecond)

	log("new seesaw")
	ss := seesaw.New(defaultNeoTrellisAddress, i2c)

	log("reset seesaw")
	err = ss.SoftReset()
	if err != nil {
		fatal(err)
	}

	log("init neopixel")

	err = ss.WriteRegister(seesaw.ModuleNeoPixelBase, seesaw.FunctionNeopixelPin, neoPixelPin)
	if err != nil {
		fatal(err)
	}

	time.Sleep(seesawWriteDelay)

	bl := len(ledBuffer)
	tx := []byte{0, byte(bl)}
	err = ss.Write(seesaw.ModuleNeoPixelBase, seesaw.FunctionNeopixelBufLength, tx)
	if err != nil {
		fatal(err)
	}

	time.Sleep(seesawWriteDelay)

	err = ss.WriteRegister(seesaw.ModuleNeoPixelBase, seesaw.FunctionNeopixelSpeed, 0x01)
	if err != nil {
		fatal(err)
	}

	time.Sleep(seesawWriteDelay)

	log("start rainbow")
	i := uint8(0)
	for {
		r, g, b := colorWheel(i * 13)
		p := (i % usedPixelCount) * 3
		ledBuffer[p] = g
		ledBuffer[p+1] = r
		ledBuffer[p+2] = b

		log("writing rainbow")
		err = ss.Write(seesaw.ModuleNeoPixelBase, seesaw.FunctionNeopixelBuf, ledBuffer)
		if err != nil {
			fatal(err)
		}

		time.Sleep(100 * time.Millisecond)

		log("showing rainbow")
		err = ss.Write(seesaw.ModuleNeoPixelBase, seesaw.FunctionNeopixelShow, nil)
		if err != nil {
			fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
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
