package draw

import (
	"time"
)

type RGB struct {
	R, G, B uint8
}
type Buffer4x4 [16]RGB

func (p *Buffer4x4) Set(x, y uint8, c RGB) {
	offset := 4*x + y
	p[offset] = c
}

type Display interface {
	WriteBuffer(b *Buffer4x4) error
}

type Animation interface {
	Update(now time.Time)
	Draw(d Display) error
}
