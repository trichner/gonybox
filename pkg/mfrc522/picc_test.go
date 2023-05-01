package mfrc522

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestBCC(t *testing.T) {
	//1c000000

	log.Printf("TEST %x\n", 0x1c^0x00^0x00^0x00)
}
func TestDevice_PiccSelect(t *testing.T) {

	cmd := PiccSelectCommand{}
	cmd.setCommand(0x01)
	cmd.setNumberOfValidBits(0x02)
	cmd.setUuidData(&[4]byte{0x3, 0x4, 0x5, 0x6})
	cmd.updateBlockCheckCharacter()

	cmd.updateCrc(func(data []byte, crc []byte) error {
		var buf strings.Builder
		for _, b := range data {
			buf.WriteString(fmt.Sprintf("%0.2X", b))
		}
		crc[0] = 0xBE
		crc[1] = 0xEF
		return nil
	})

	var buf strings.Builder
	for _, b := range cmd.slice() {
		buf.WriteString(fmt.Sprintf("%0.2X ", b))
	}
	log.Println(buf.String())
}
