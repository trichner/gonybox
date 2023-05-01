package debug

import "machine"

func init() {
}
func Log(s string) {
	//sudo screen /dev/ttyACM0 9600
	machine.Serial.Write([]byte(s + "\r\n"))
}

func FmtByteToBinary(r byte) string {
	formatted := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		if r&0x1 != 0 {
			formatted[i] = '1'
		} else {
			formatted[i] = '0'
		}
		r >>= 1
	}
	return string(formatted)
}

var mapping = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func FmtSliceToHex(s []byte) string {
	formatted := make([]byte, len(s)*2)
	for i := 0; i < len(s); i++ {
		formatted[i*2] = mapping[s[i]>>4]
		formatted[i*2+1] = mapping[s[i]&0xf]
	}
	return string(formatted)
}

func FmtByteToHex(s byte) string {
	formatted := make([]byte, 2)
	formatted[0] = mapping[s>>4]
	formatted[1] = mapping[s&0xf]
	return string(formatted)
}
