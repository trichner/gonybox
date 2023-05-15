package animations

import (
	"time"
	"trelligo/pkg/neotrellis"
	"trelligo/pkg/shims/rand"
)

type Animation interface {
	Update(now time.Time)
	Draw(dev *neotrellis.Device) error
}

func AnimateFor(dev *neotrellis.Device, a Animation, duration time.Duration) error {
	start := time.Now()
	now := start
	for now.Sub(start) < duration {
		now = time.Now()
		a.Update(now)
		err := a.Draw(dev)
		if err != nil {
			return err
		}
	}
	return nil
}

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

func drawBuffer(dev *neotrellis.Device, buf []neotrellis.RGB) error {
	err := dev.WriteColors(buf)
	if err != nil {
		return err
	}
	return dev.ShowPixels()
}
