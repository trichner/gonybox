package mcu

import (
	"log"
	"machine"
	"time"
	"trelligo/dfplayer"
)

var _ = dfplayer.RoundTripper(&RoundTripper{})

type RoundTripper struct {
	port     *machine.UART
	rxBuffer []byte
}

func NewRoundTripper(uart *machine.UART) *RoundTripper {
	return &RoundTripper{
		port:     uart,
		rxBuffer: make([]byte, 10),
	}
}

func (u *RoundTripper) Send(tx *dfplayer.Frame, rx *dfplayer.Frame) error {
	// clear the input buffer as we don't care about things already queued
	u.port.Buffer.Clear()

	_, err := u.port.Write(tx[:])
	if err != nil {
		return err
	}
	deadline := time.Now().Add(time.Millisecond * 200)
	return u.readToDeadline(rx, deadline)
}

func (u *RoundTripper) readToDeadline(rx *dfplayer.Frame, deadline time.Time) error {
	count := 0
	u.rxBuffer = u.rxBuffer[:0]
	buf := make([]byte, 10)
	for count < 10 {
		n, err := u.port.Read(buf)
		if err != nil {
			return err
		}
		if n == 0 && time.Now().After(deadline) {
			return dfplayer.ErrDeviceTimeout
		}
		count += n
		u.rxBuffer = append(u.rxBuffer, buf[:n]...)
	}
	if count != 10 {
		log.Fatalf("expected 10 bytes, got: %d", count)
	}
	arr := (*[10]byte)(u.rxBuffer)
	*rx = *arr
	return nil
}
