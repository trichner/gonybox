package neotrellis

type PixelBuffer []RGB

func NewPixelBuffer() PixelBuffer {
	return make(PixelBuffer, xCount*yCount)
}
func (p PixelBuffer) SetPixel(x, y uint8, c RGB) {
	offset := PositionFromXY(x, y).PixelOffset()
	p[offset] = c
}

func (p PixelBuffer) Pixel(x, y uint8) RGB {
	offset := PositionFromXY(x, y).PixelOffset()
	return p[offset]
}
