package usbtty

import (
	"go.bug.st/serial"
	"log"
	"time"
	"trelligo/dfplayer"
)

var _ = dfplayer.RoundTripper(&NucleoRoundTripper{})

type NucleoRoundTripper struct {
	port     serial.Port
	rxBuffer []byte
}

func (u *NucleoRoundTripper) Send(tx *dfplayer.Frame, rx *dfplayer.Frame) error {
	err := u.port.ResetInputBuffer()
	if err != nil {
		return err
	}
	_, err = u.port.Write(tx[:])
	if err != nil {
		return err
	}
	deadline := time.Now().Add(time.Millisecond * 200)
	return u.readToDeadline(rx, deadline)
}

func (u *NucleoRoundTripper) readToDeadline(rx *dfplayer.Frame, deadline time.Time) error {
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

func NewUsbTty(device string) *NucleoRoundTripper {
	//ex: /dev/ttyUSB0
	// 9600 8N1
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(device, mode)
	if err != nil {
		log.Fatal(err)
	}
	port.SetReadTimeout(time.Millisecond * 500)

	//buf := make([]byte, 10)
	//n, err := port.Read(buf)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("RX: %v", buf[:n])

	return &NucleoRoundTripper{port: port, rxBuffer: make([]byte, 10)}
}
