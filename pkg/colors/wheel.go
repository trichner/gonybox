package colors

// ColorAt returns the color within an RGB wheel at the given position.
// This wraps around seamlessly, ideal for rainbows.
func ColorAt(p uint8) (r, g, b uint8) {
	p = 255 - p
	if p < 85 {
		return 255 - p*3, 0, p * 3
	}
	if p < 170 {
		p -= 85
		return 0, p * 3, 255 - p*3
	}
	p -= 170
	return p * 3, 255 - p*3, 0
}
