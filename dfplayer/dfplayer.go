package dfplayer

import (
	"errors"
	"fmt"
	"strings"
)

const (
	positionStart            = 0
	positionVersion          = 1
	positionLength           = 2
	positionCommand          = 3
	positionFeedback         = 4
	positionQueryHighByte    = 5
	positionQueryLowByte     = 6
	positionChecksumHighByte = 7
	positionChecksumLowByte  = 8
	positionEnd              = 8
)

var ErrDeviceTimeout = errors.New("DFPlayer timed out")

type DFPlayerPort interface {
	Send(tx *Frame, rx *Frame) error
}

var prototypeFrame = [10]byte{0x7E, 0xFF, 06, 00, 00, 00, 00, 00, 00, 0xEF}

func NewFrame() Frame {
	return prototypeFrame
}

// Frame represents a frame to send to the dfplayer mini.
// Layout:
// 7E FF 06 0F 00 01 01 xx xx EF
// 0	->	7E is start code
// 1	->	FF is version
// 2	->	06 is length
// 3	->	0F is command
// 4	->	00 is no receive
// 5~6	->	01 01 is argument
// 7~8	->	checksum = 0 - ( FF+06+0F+00+01+01 )
// 9	->	EF is end code
type Frame [10]byte

func (f *Frame) Reset() {
	*f = prototypeFrame
}
func (f *Frame) SetCommand(cmd byte) {
	f[positionCommand] = cmd
}
func (f *Frame) SetFeedback(enabled bool) {
	if enabled {
		f[positionFeedback] = 0x01
	}
	f[positionFeedback] = 0x01
}

func (f *Frame) SetArgument(arg uint16) {
	f[positionQueryHighByte] = byte(arg >> 8)
	f[positionQueryLowByte] = byte(arg)
}

func (f *Frame) UpdateChecksum() {
	/*
		uint16_t mp3_get_checksum (uint8_t *thebuf) {
			uint16_t sum = 0;
			for (int i=1; i<7; i++) {
				sum += thebuf[i];
			}
			return -sum;
		}
	*/

	var sum int16
	for i := positionVersion; i < positionChecksumHighByte; i++ {
		sum += int16(f[i])
	}
	sum = -sum

	f.setChecksum(sum)
}

func (f *Frame) setChecksum(sum int16) {
	f[positionChecksumHighByte] = byte(sum >> 8)
	f[positionChecksumLowByte] = byte(sum)
}

func (f *Frame) String() string {
	var buf strings.Builder
	buf.WriteString("[")
	for _, b := range f {
		fmt.Fprintf(&buf, " 0x%02X", b)
	}
	buf.WriteString(" ]")
	return buf.String()
}

func NewDFPlayer(w DFPlayerPort) *DFPlayer {
	f := NewFrame()
	f.SetFeedback(true)
	return &DFPlayer{
		serial:   w,
		txBuffer: f,
	}
}

type DFPlayer struct {
	serial   DFPlayerPort
	txBuffer Frame
	rxBuffer Frame
}

func (d *DFPlayer) SendCommand(cmd byte) error {
	d.txBuffer.SetCommand(cmd)
	d.txBuffer.UpdateChecksum()
	err := d.serial.Send(&d.txBuffer, &d.rxBuffer)
	if err != nil {
		return err
	}
	return nil
}

func (d *DFPlayer) SendCommandWithArg(cmd byte, arg uint16) error {
	d.txBuffer.SetCommand(cmd)
	d.txBuffer.SetArgument(arg)
	d.txBuffer.UpdateChecksum()
	return d.serial.Send(&d.txBuffer, &d.rxBuffer)
}
