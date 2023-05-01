package dfplayer

import (
	"fmt"
	"strings"
)

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
	} else {
		f[positionFeedback] = 0x00
	}
}

func (f *Frame) SetArgument(arg uint16) {
	f[positionQueryHighByte] = byte(arg >> 8)
	f[positionQueryLowByte] = byte(arg)
}

func (f *Frame) UpdateChecksum() {
	/*
		// Reference implementation:
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
