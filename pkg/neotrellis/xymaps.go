package neotrellis

// Position maps between X/Y coordinates, seesaw key and neoTrellis pixels.
// There are three representations of keys on the NeoTrellis:
//
// 1. NeoTrellis pixel offsets, sequential number through 0 to 15
//
// 2. Seesaw keys, the key IDs as wired up to the seesaw chip
//
// 3. x/y coordinates, starting at the bottom left
//
// NeoTrellis pixel offset arrangement:
//
//	03 07 11 15
//	02 06 10 14
//	01 05 09 13
//	00 04 08 12
//
// Seesaw key arrangement:
//
//	03 11 19 27
//	02 10 18 26
//	01 09 17 25
//	00 08 16 24
//
// x/y coordinate arrangement:
//
//	03 13 23 33
//	02 12 22 32
//	01 11 21 31
//	00 10 20 30
type Position uint8

func PositionFromXY(x, y uint8) Position {
	return Position(x<<4 | y&0x0F)
}
func newXyFromSeesawKey(k uint8) Position {
	return PositionFromXY(seesawKeyToXy(k))
}
func (v Position) X() uint8 {
	return uint8(v >> 4)
}

func (v Position) Y() uint8 {
	return uint8(v & 0x0F)
}

func (v Position) KeyID() uint8 {
	return xyToSeesawKey(v.X(), v.Y())
}

func (v Position) PixelOffset() uint16 {
	// returning uint16 because that's what NeoPixel offsets want
	return uint16(v.X()*4 + v.Y())
}

func xyToSeesawKey(x, y uint8) uint8 {
	return ((x*4+y)/4)*8 + ((x*4 + y) % 4)
}

func seesawKeyToXy(k uint8) (uint8, uint8) {
	offset := ((k)/8)*4 + ((k) % 8)
	return offset / 4, offset % 4
}
