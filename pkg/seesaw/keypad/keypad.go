package keypad

import (
	"time"
	"trelligo/pkg/seesaw"
	"unsafe"
)

type Edge uint8

const (
	EdgeHigh Edge = iota
	EdgeLow
	EdgeFalling
	EdgeRising
)

// KeyEvent represents a pressed or released key
type KeyEvent uint8

func NewKeyEvent(b byte) KeyEvent {
	return KeyEvent(b)
}

func (k KeyEvent) Edge() Edge {
	return Edge(k & 0b11)
}

func (k KeyEvent) Key() uint8 {
	return uint8(k >> 2)
}

type SeesawKeypad struct {
	seesaw *seesaw.Device
}

func New(dev *seesaw.Device) *SeesawKeypad {
	return &SeesawKeypad{seesaw: dev}
}

// KeyEventCount returns the number of pending KeyEvent s in the FIFO queue
func (s *SeesawKeypad) KeyEventCount() (uint8, error) {
	//https://github.com/adafruit/Adafruit_Seesaw/blob/master/Adafruit_seesaw.cpp#L721
	buf := make([]byte, 1)
	err := s.seesaw.Read(seesaw.SEESAW_KEYPAD_BASE, seesaw.SEESAW_KEYPAD_COUNT, buf, 500*time.Microsecond)
	return buf[0], err
}

// SetKeypadInterrupt enables or disables interrupts for key events
func (s *SeesawKeypad) SetKeypadInterrupt(enable bool) error {
	if enable {
		return s.seesaw.WriteRegister(seesaw.SEESAW_KEYPAD_BASE, seesaw.SEESAW_KEYPAD_INTENSET, 0x1)
	}
	return s.seesaw.WriteRegister(seesaw.SEESAW_KEYPAD_BASE, seesaw.SEESAW_KEYPAD_INTENCLR, 0x1)
}

// Read reads pending KeyEvent s from the FIFO
func (s *SeesawKeypad) Read(buf []KeyEvent) error {
	//https://github.com/adafruit/Adafruit_Seesaw/blob/master/Adafruit_seesaw.cpp#LL732C21-L732C21
	bytesBuf := *(*[]byte)(unsafe.Pointer(&buf))
	return s.seesaw.Read(seesaw.SEESAW_KEYPAD_BASE, seesaw.SEESAW_KEYPAD_FIFO, bytesBuf, time.Millisecond)
}

// ConfigureKeypad enables or disables a key and edge on the keypad module
func (s *SeesawKeypad) ConfigureKeypad(key uint8, edge Edge, enable bool) error {

	//set STATE
	state := byte(0)
	if enable {
		state |= 0x01
	}

	//set ACTIVE
	state |= (1 << edge) << 1
	return s.seesaw.Write(seesaw.SEESAW_KEYPAD_BASE, seesaw.SEESAW_KEYPAD_EVENT, []byte{key, state})
}
