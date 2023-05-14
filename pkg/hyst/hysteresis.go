// Package hyst implements software hysteresis for a noisy source such as an ADC
package hyst

// generate a lookup-table (LUT) mapping from [0,0xFFFF] to [0,30]
//go:generate go run gen_lut.go -start=2047 -length=31 -name=lut0to30

// Getter for noisy uint16 values
type Getter interface {
	Get() uint16
}

type Hysteresis struct {
	getter    Getter
	lut       []uint16
	pos       int
	threshold uint16
}

func New(g Getter, threshold uint16) *Hysteresis {
	return &Hysteresis{
		getter:    g,
		lut:       lut0to30,
		pos:       -1,
		threshold: threshold,
	}
}

func (d *Hysteresis) Get() (int, bool) {
	next := d.getter.Get()
	if d.pos < 0 {
		d.set(next)
		return d.pos, true
	}

	u := d.update(next)
	return d.pos, u
}

func (d *Hysteresis) update(n uint16) bool {
	lower := d.lut[d.pos]
	//bottomed out?
	if n <= lower && d.pos == 0 {
		return false
	}

	//topped out?
	if n >= lower && d.pos == len(d.lut)-1 {
		return false
	}

	upper := uint16(0xFFFF)
	if len(d.lut) > d.pos+1 {
		upper = d.lut[d.pos+1]
	}
	if csub(lower, d.threshold) <= n && n <= cadd(upper, d.threshold) {
		//same range
		return false
	}

	d.set(n)
	return true
}

func (d *Hysteresis) set(value uint16) {
	prev := 0
	for i, v := range d.lut {
		if value < v {
			d.pos = prev
			return
		}
		prev = i
	}
	d.pos = len(d.lut) - 1
}

func cadd(augend, addend uint16) uint16 {
	if 0xFFFF-augend <= addend {
		return 0xFFFF
	}
	return augend + addend
}
func csub(minuend, subtrahend uint16) uint16 {
	if subtrahend >= minuend {
		return 0
	}
	return minuend - subtrahend
}
