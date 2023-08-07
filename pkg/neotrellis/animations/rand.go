package animations

import (
	"time"
	"trelligo/pkg/neotrellis"
	"trelligo/pkg/shims/rand"
)

type randomBlink struct {
	buf []neotrellis.RGB
	rnd *rand.Rand

	lastUpdate time.Time
}

func NewRandomBlink(r *rand.Rand) Animation {
	buf := make([]neotrellis.RGB, 16)
	for i := range buf {
		buf[i] = colorWheel(uint8(r.Uint32()))
	}

	return &randomBlink{
		buf: buf,
		rnd: r,
	}
}

func (r *randomBlink) Draw(dev *neotrellis.Device) error {
	return drawBuffer(dev, r.buf)
}
func (r *randomBlink) Update(now time.Time) {

	if now.Sub(r.lastUpdate) < 100*time.Millisecond {
		return
	}
	r.lastUpdate = now

	v := r.rnd.Uint32()
	color := colorWheel(uint8(v))
	p := (v >> 8) % 16
	r.buf[p] = color

}
