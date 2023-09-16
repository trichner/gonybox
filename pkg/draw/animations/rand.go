package animations

import (
	"time"
	"trelligo/pkg/draw"
	"trelligo/pkg/shims/rand"
)

type randomBlink struct {
	buf draw.Buffer4x4
	rnd *rand.Rand

	lastUpdate time.Time
}

func NewRandomBlink(r *rand.Rand) draw.Animation {
	b := &randomBlink{rnd: r}

	for i := range b.buf {
		b.buf[i] = colorWheel(uint8(r.Uint32()))
	}

	return b
}

func (r *randomBlink) Draw(display draw.Display) error {
	return display.WriteBuffer(&r.buf)
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
