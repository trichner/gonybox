package usbtty

import (
	"go.bug.st/serial"
	"log"
	"time"
	dfplayer2 "trelligo/pkg/dfplayer"
)

var _ = dfplayer2.RoundTripper(&UsbTty{})

type UsbTty struct {
	port      serial.Port
	rxBuffer  []byte
	rxTimeout time.Duration
}

func (u *UsbTty) Send(tx *dfplayer2.Frame, rx *dfplayer2.Frame) error {
	err := u.port.ResetInputBuffer()
	if err != nil {
		return err
	}
	_, err = u.port.Write(tx[:])
	if err != nil {
		return err
	}
	deadline := time.Now().Add(u.rxTimeout)
	return u.readToDeadline(rx, deadline)
}

func (u *UsbTty) readToDeadline(rx *dfplayer2.Frame, deadline time.Time) error {
	count := 0
	u.rxBuffer = u.rxBuffer[:0]
	buf := make([]byte, 10)
	for count < 10 {
		n, err := u.port.Read(buf)
		if err != nil {
			return err
		}
		if n == 0 && time.Now().After(deadline) {
			return dfplayer2.ErrDeviceTimeout
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

func NewUsbTty(device string) *UsbTty {
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

	return &UsbTty{port: port, rxBuffer: make([]byte, 10), rxTimeout: time.Millisecond * 200}
}
