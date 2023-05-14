package animations

import (
	"time"
	"trelligo/pkg/neotrellis"
	"trelligo/pkg/shims/rand"
)

func RandomBlink(r rand.Rand, dev *neotrellis.Device) {
	buf := make([]neotrellis.RGB, 16)

	for i := range buf {
		buf[i] = colorWheel(uint8(r.Uint32()))
	}

	err := dev.WriteColors(buf)
	if err != nil {
		panic(err)
	}

	err = dev.ShowPixels()
	if err != nil {
		panic(err)
	}

	for {

		v := r.Uint32()
		color := colorWheel(uint8(v))
		x := uint8((v >> 8) & 0x11)
		y := uint8((v >> 12) & 0x11)
		err = dev.SetPixelColor(x, y, color)
		if err != nil {
			panic(err)
		}
		err = dev.ShowPixels()
		if err != nil {
			panic(err)
		}
		time.Sleep(50 * time.Millisecond)
	}

}
func InfiniteRainbowRun(dev *neotrellis.Device) {
	buf := make([]neotrellis.RGB, 16)

	err := dev.WriteColors(buf)
	if err != nil {
		panic(err)
	}
	err = dev.ShowPixels()
	if err != nil {
		panic(err)
	}

	i := uint8(0)
	for {
		p := i % 16

		x := p / 4
		y := p % 4

		color := colorWheel(i * 13)
		err = dev.SetPixelColor(x, y, color)
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

func colorWheel(p uint8) neotrellis.RGB {
	p = 255 - p
	if p < 85 {
		return neotrellis.RGB{255 - p*3, 0, p * 3}
	}
	if p < 170 {
		p -= 85
		return neotrellis.RGB{0, p * 3, 255 - p*3}
	}
	p -= 170
	return neotrellis.RGB{p * 3, 255 - p*3, 0}
}
