package main

import (
	"machine"
	"time"
	"trelligo/pkg/neotrellis"
)

func main() {
	// an example to use the NeoPixels on the NeoTrellis board
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{})
	if err != nil {
		panic(err)
	}

	dev, err := neotrellis.New(i2c, neotrellis.DefaultNeoTrellisAddress)
	if err != nil {
		panic(err)
	}

	i := uint8(0)
	for {
		r, g, b := colorWheel(i * 13)
		p := i % 16
		x := p / 4
		y := p % 4
		err = dev.SetPixelColor(x, y, r, g, b)
		if err != nil {
			panic(err)
		}
		err = dev.ShowPixels()
		if err != nil {
			panic(err)
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
