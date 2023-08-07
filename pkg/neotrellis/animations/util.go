package animations

import "trelligo/pkg/neotrellis"

func colorWheel(p uint8) neotrellis.RGB {
	p = 255 - p
	if p < 85 {
		return neotrellis.RGB{R: 255 - p*3, B: p * 3}
	}
	if p < 170 {
		p -= 85
		return neotrellis.RGB{G: p * 3, B: 255 - p*3}
	}
	p -= 170
	return neotrellis.RGB{R: p * 3, G: 255 - p*3}
}

func drawBuffer(dev *neotrellis.Device, buf []neotrellis.RGB) error {
	err := dev.WriteColors(buf)
	if err != nil {
		return err
	}
	return dev.ShowPixels()
}
