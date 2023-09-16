package minigames

import (
	"time"
	"trelligo/pkg/draw"
	"trelligo/pkg/rbuf"
)

type EventType uint8

const (
	Unknown EventType = iota
	KeyDown
	KeyUp
)

type Event interface {
	X() uint8
	Y() uint8
	Type() EventType
}

type Game interface {
	Update(now time.Time, events rbuf.RingBuffer[Event])
	Draw(d draw.Display) error
}
type WhackAMole struct {
}
