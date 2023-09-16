package animations

import (
	"time"
	"trelligo/pkg/draw"
)

type infinityRainbow struct {
	buf        draw.Buffer4x4
	iteration  uint8
	lastUpdate time.Time
}

func NewInfinityRainbow() draw.Animation {
	return &infinityRainbow{}
}

func (i *infinityRainbow) Draw(display draw.Display) error {
	return display.WriteBuffer(&i.buf)
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
