package ntdisplay

import (
	"trelligo/pkg/draw"
	"trelligo/pkg/neotrellis"
)

type Display struct {
	dev *neotrellis.Device
	buf []neotrellis.RGB
}

func NewDisplay(dev *neotrellis.Device) *Display {
	buf := make([]neotrellis.RGB, 16)
	return &Display{
		dev: dev,
		buf: buf,
	}
}

func (n *Display) WriteBuffer(b *draw.Buffer4x4) error {

	for i := 0; i < len(b); i++ {
		c := b[i]
		c = draw.GammaCorrect(c)
		n.buf[i] = neotrellis.RGB{R: c.R, G: c.G, B: c.B}
	}

	err := n.dev.WriteColors(n.buf)
	if err != nil {
		return err
	}

	return n.dev.ShowPixels()
}
