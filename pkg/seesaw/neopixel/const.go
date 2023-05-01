package neopixel

// Constants for NeoPixels on the Seesaw
// https://github.com/adafruit/Adafruit_Seesaw/blob/master/seesaw_neopixel.h#L49

// PixelType defines the order of primary colors in the NeoPixel data stream, which can vary
// among device types, manufacturers and even different revisions of
// the same item.  The PixelType encodes the per-pixel byte offsets of the red, green
// and blue primaries (plus white, if present) in the data stream.
//
// Below an easier-to-use named version for
// each permutation.  e.g. NEO_GRB indicates a NeoPixel-compatible
// device expecting three bytes per pixel, with the first byte
// containing the green value, second containing red and third
// containing blue.  The in-memory representation of a chain of
// NeoPixels is the same as the data-stream order; no re-ordering of
// bytes is required when issuing data to the chain.
//
// Bits 5,4 of this value are the offset (0-3) from the first byte of
// a pixel to the location of the red color byte.  Bits 3,2 are the
// green offset and 1,0 are the blue offset.  If it is an RGBW-type
// device (supporting a white primary in addition to R,G,B), bits 7,6
// are the offset to the white byte...otherwise, bits 7,6 are set to
// the same value as 5,4 (red) to indicate an RGB (not RGBW) device.
//
// i.e. binary representation:
//
//	0bWWRRGGBB for RGBW devices
//	0bRRRRGGBB for RGB
type PixelType uint16

func (p PixelType) BlueOffset() int {
	return int(byte(p) & 0b11)
}

func (p PixelType) GreenOffset() int {
	return int((byte(p) >> 2) & 0b11)
}

func (p PixelType) RedOffset() int {
	return int((byte(p) >> 4) & 0b11)
}

func (p PixelType) WhiteOffset() int {
	return int((byte(p) >> 6) & 0b11)
}

func (p PixelType) IsRGBW() bool {
	return p.RedOffset() != p.WhiteOffset()
}

func (p PixelType) EncodedLen() int {
	if p.IsRGBW() {
		return 4
	}
	return 3
}

func (p PixelType) Is800KHz() bool {
	return (uint16(p) & 0xFF00) == NEO_KHZ800
}

func (p PixelType) PutRGBW(buf []byte, r, g, b, w uint8) int {
	buf[p.RedOffset()] = r
	buf[p.GreenOffset()] = g
	buf[p.BlueOffset()] = b
	if !p.IsRGBW() {
		// if we don't have white, skip it
		return 3
	}
	buf[p.WhiteOffset()] = w
	return 4
}

// RGB NeoPixel permutations; white and red offsets are always same
// Offset:                       W          R          G          B
const (
	NEO_RGB PixelType = (0 << 6) | (0 << 4) | (1 << 2) | (2)
	NEO_RBG PixelType = (0 << 6) | (0 << 4) | (2 << 2) | (1)
	NEO_GRB PixelType = (1 << 6) | (1 << 4) | (0 << 2) | (2)
	NEO_GBR PixelType = (2 << 6) | (2 << 4) | (0 << 2) | (1)
	NEO_BRG PixelType = (1 << 6) | (1 << 4) | (2 << 2) | (0)
	NEO_BGR PixelType = (2 << 6) | (2 << 4) | (1 << 2) | (0)
)

// RGBW NeoPixel permutations; all 4 offsets are distinct
// Offset:                        W          R          G          B
const (
	NEO_WRGB PixelType = (0 << 6) | (1 << 4) | (2 << 2) | (3)
	NEO_WRBG PixelType = (0 << 6) | (1 << 4) | (3 << 2) | (2)
	NEO_WGRB PixelType = (0 << 6) | (2 << 4) | (1 << 2) | (3)
	NEO_WGBR PixelType = (0 << 6) | (3 << 4) | (1 << 2) | (2)
	NEO_WBRG PixelType = (0 << 6) | (2 << 4) | (3 << 2) | (1)
	NEO_WBGR PixelType = (0 << 6) | (3 << 4) | (2 << 2) | (1)
	//
	NEO_RWGB PixelType = (1 << 6) | (0 << 4) | (2 << 2) | (3)
	NEO_RWBG PixelType = (1 << 6) | (0 << 4) | (3 << 2) | (2)
	NEO_RGWB PixelType = (2 << 6) | (0 << 4) | (1 << 2) | (3)
	NEO_RGBW PixelType = (3 << 6) | (0 << 4) | (1 << 2) | (2)
	NEO_RBWG PixelType = (2 << 6) | (0 << 4) | (3 << 2) | (1)
	NEO_RBGW PixelType = (3 << 6) | (0 << 4) | (2 << 2) | (1)
	//
	NEO_GWRB PixelType = (1 << 6) | (2 << 4) | (0 << 2) | (3)
	NEO_GWBR PixelType = (1 << 6) | (3 << 4) | (0 << 2) | (2)
	NEO_GRWB PixelType = (2 << 6) | (1 << 4) | (0 << 2) | (3)
	NEO_GRBW PixelType = (3 << 6) | (1 << 4) | (0 << 2) | (2)
	NEO_GBWR PixelType = (2 << 6) | (3 << 4) | (0 << 2) | (1)
	NEO_GBRW PixelType = (3 << 6) | (2 << 4) | (0 << 2) | (1)
	//
	NEO_BWRG PixelType = (1 << 6) | (2 << 4) | (3 << 2) | (0)
	NEO_BWGR PixelType = (1 << 6) | (3 << 4) | (2 << 2) | (0)
	NEO_BRWG PixelType = (2 << 6) | (1 << 4) | (3 << 2) | (0)
	NEO_BRGW PixelType = (3 << 6) | (1 << 4) | (2 << 2) | (0)
	NEO_BGWR PixelType = (2 << 6) | (3 << 4) | (1 << 2) | (0)
	NEO_BGRW PixelType = (3 << 6) | (2 << 4) | (1 << 2) | (0)
)

// If 400 KHz support is enabled, the third parameter to the constructor
// requires a 16-bit value (in order to select 400 vs 800 KHz speed).
// If only 800 KHz is enabled (as is default on ATtiny), an 8-bit value
// is sufficient to encode pixel color order, saving some space.
const (
	NEO_KHZ800 = 0x0000 // 800 KHz datastream
	NEO_KHZ400 = 0x0100 // 400 KHz datastream
)
