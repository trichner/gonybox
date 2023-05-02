package neopixel

import (
	"fmt"
	"strconv"
	"trelligo/pkg/seesaw"
)

type RGBWColor [4]uint8

func NewRGBW(r, g, b, w uint8) RGBWColor {
	return [4]uint8{r, g, b, w}
}

type SeesawNeopixel struct {
	seesaw    *seesaw.Device
	numLEDs   int
	pixels    []RGBWColor
	pin       uint8
	pixelType PixelType
}

func New(dev *seesaw.Device, pin uint8, numLEDs int, pixelType PixelType) (*SeesawNeopixel, error) {

	pixel := &SeesawNeopixel{
		seesaw: dev,
	}

	err := pixel.updatePin(pin)
	if err != nil {
		return nil, fmt.Errorf("failed to update seesaw NeoPixel pin: %w", err)
	}

	err = pixel.setNumberOfLEDs(numLEDs)
	if err != nil {
		return nil, fmt.Errorf("failed to update LED count: %w", err)
	}

	err = pixel.updatePixelType(pixelType)
	if err != nil {
		return nil, fmt.Errorf("failed to update pixel type: %w", err)
	}

	return pixel, nil
}

func (s *SeesawNeopixel) setNumberOfLEDs(n int) error {

	s.numLEDs = n
	s.reinitLeds()

	lenBytes := calculateBufferLength(n, s.pixelType)
	buf := []byte{byte(lenBytes >> 8), byte(lenBytes & 0xFF)}
	return s.seesaw.Write(seesaw.SEESAW_NEOPIXEL_BASE, seesaw.SEESAW_NEOPIXEL_BUF_LENGTH, buf)
}

func calculateBufferLength(ledCount int, pixelType PixelType) int {
	return ledCount * pixelType.EncodedLen()
}

func (s *SeesawNeopixel) reinitLeds() {
	s.pixels = make([]RGBWColor, s.numLEDs)
}

func (s *SeesawNeopixel) updatePixelType(t PixelType) error {
	old := s.pixelType
	s.pixelType = t

	if old.IsRGBW() != t.IsRGBW() {
		//byte-size changed, re-init buffer
		s.reinitLeds()
	}

	speed := byte(0)
	if t.Is800KHz() {
		speed = 1
	}

	return s.seesaw.WriteRegister(seesaw.SEESAW_NEOPIXEL_BASE, seesaw.SEESAW_NEOPIXEL_SPEED, speed)
}

func (s *SeesawNeopixel) updatePin(pin uint8) error {
	s.pin = pin
	return s.seesaw.WriteRegister(seesaw.SEESAW_NEOPIXEL_BASE, seesaw.SEESAW_NEOPIXEL_PIN, pin)
}

func (s *SeesawNeopixel) SetPixelColor(offset uint16, r, g, b, w uint8) error {

	encodedLen := s.pixelType.EncodedLen()

	buf := make([]byte, 2+encodedLen)
	l := s.pixelType.PutRGBW(buf[2:], r, g, b, w)
	if l != encodedLen {
		panic("unexpected encoded length: " + strconv.Itoa(l) + " != " + strconv.Itoa(encodedLen))
	}
	byteOffset := offset * uint16(encodedLen)
	buf[0] = uint8(byteOffset >> 8)
	buf[1] = uint8(byteOffset)
	return s.seesaw.Write(seesaw.SEESAW_NEOPIXEL_BASE, seesaw.SEESAW_NEOPIXEL_BUF, buf)
}

func (s *SeesawNeopixel) ShowPixels() error {
	return s.seesaw.Write(seesaw.SEESAW_NEOPIXEL_BASE, seesaw.SEESAW_NEOPIXEL_SHOW, nil)
}
