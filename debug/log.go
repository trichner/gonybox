package debug

import "machine"

func init() {
}
func Log(s string) {
	//sudo screen /dev/ttyACM0 9600
	machine.Serial.Write([]byte(s + "\n\r"))
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
