package animations

import (
	"trelligo/pkg/draw"
	"trelligo/pkg/neotrellis"
)

func colorWheel(p uint8) draw.RGB {
	p = 255 - p
	if p < 85 {
		return draw.RGB{R: 255 - p*3, B: p * 3}
	}
	if p < 170 {
		p -= 85
		return draw.RGB{G: p * 3, B: 255 - p*3}
	}
	p -= 170
	return draw.RGB{R: p * 3, G: 255 - p*3}
}

func drawBuffer(dev *neotrellis.Device, buf []neotrellis.RGB) error {
	err := dev.WriteColors(buf)
	if err != nil {
		return err
	}
	return dev.ShowPixels()
}
