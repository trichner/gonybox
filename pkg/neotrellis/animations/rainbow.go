package animations

import (
	"time"
	"trelligo/pkg/neotrellis"
)

type infinityRainbow struct {
	buf        []neotrellis.RGB
	iteration  uint8
	lastUpdate time.Time
}

func NewInfinityRainbow() Animation {
	buf := make([]neotrellis.RGB, 16)

	return &infinityRainbow{
		buf: buf,
	}
}

func (i *infinityRainbow) Draw(dev *neotrellis.Device) error {
	return drawBuffer(dev, i.buf)
}
func (i *infinityRainbow) Update(now time.Time) {

	if now.Sub(i.lastUpdate) < 50*time.Millisecond {
		return
	}
	i.lastUpdate = now

	p := i.iteration % 16

	color := colorWheel(i.iteration * 13)
	i.buf[p] = color

	i.iteration++
}
