package player

import (
	"fmt"
	"strconv"
	"time"
	"trelligo/pkg/debug"
	"trelligo/pkg/dfplayer"
	"trelligo/pkg/errwrap"
	"trelligo/pkg/neotrellis"
	"trelligo/pkg/seesaw/keypad"
)

const minDelay = time.Millisecond * 100

type keyHandlerFunc func(x, y uint8, e keypad.Edge) error

type xy = uint8

func newXy(x, y uint8) xy {
	return y<<2 | x
}

type VolumeGetter interface {
	Get() (int, bool)
}

type Player struct {
	nt  *neotrellis.Device
	dfp *dfplayer.Player

	handlers    []keyHandlerFunc
	needRefresh bool

	vol VolumeGetter

	lastUpdate time.Time
	buf        neotrellis.PixelBuffer
}

// [  0  1  2  3 ]
// [  4  5  6  7 ]
// [  8  9 10 11 ]
// [  <  D  >  x ]

func New(nt *neotrellis.Device, dfp *dfplayer.Player, getter VolumeGetter) (*Player, error) {

	p := &Player{
		nt:          nt,
		dfp:         dfp,
		handlers:    make([]keyHandlerFunc, 16),
		needRefresh: true,
		vol:         getter,
		buf:         neotrellis.NewPixelBuffer(),
	}

	for i := uint8(0); i < 16; i++ {
		err := nt.ConfigureKeypad(i/4, i%4, keypad.EdgeRising, true)
		if err != nil {
			return nil, fmt.Errorf("failed to enable keys: %w", err)
		}
	}

	nt.SetKeyHandleFunc(func(x, y uint8, e keypad.Edge) error {
		f := p.handlers[newXy(x, y)]
		if f == nil {
			return nil
		}
		return f(x, y, e)
	})

	nFolders := 9

	for i := 0; i < nFolders; i++ {
		folder := uint8(i + 1)
		y := uint8(3 - (i / 4))
		x := uint8(i % 4)
		p.buf.SetPixel(x, y, neotrellis.RGB{0, 100, 150})
		p.addHandler(newXy(x, y), func(x, y uint8, e keypad.Edge) error {
			debug.Log("playing folder: " + strconv.Itoa(int(folder)))
			return dfp.PlayFolder(folder, 1)
		})
	}

	// play previous
	p.buf.SetPixel(0, 0, neotrellis.RGB{0, 100, 150})
	p.addHandler(newXy(0, 0), func(x, y uint8, e keypad.Edge) error {
		return dfp.PlayPrevious()
	})

	// play next
	p.buf.SetPixel(1, 0, neotrellis.RGB{0, 150, 100})
	p.addHandler(newXy(1, 0), func(x, y uint8, e keypad.Edge) error {
		return dfp.PlayNext()
	})

	//stop
	p.buf.SetPixel(2, 0, neotrellis.RGB{0xFF, 0, 0})
	p.addHandler(newXy(2, 0), func(x, y uint8, e keypad.Edge) error {
		return dfp.Stop()
	})

	err := p.nt.WriteColors(p.buf)
	if err != nil {
		return nil, err
	}
	err = p.nt.ShowPixels()
	if err != nil {
		return nil, err
	}

	p.lastUpdate = time.Now()
	return p, nil
}

func (p *Player) addHandler(o xy, h keyHandlerFunc) {
	p.handlers[o] = h
}

func (p *Player) Process() error {

	diff := time.Since(p.lastUpdate)
	if diff <= minDelay {
		time.Sleep(minDelay - diff)
	}

	v, updated := p.vol.Get()
	if updated {
		debug.Log(fmt.Sprintf("updating volume: %2d", v))
		err := p.dfp.SetVolume(uint8(v))
		if err != nil {
			return fmt.Errorf("failed to update volume to %d: %w", v, err)
		}
	}

	err := p.nt.ProcessKeyEvents()
	if err != nil {
		err = errwrap.Wrap("player failed to process key events", err)
		debug.Log("warn: " + err.Error())
	}

	err = p.nt.WriteColors(p.buf)
	if err != nil {
		err = errwrap.Wrap("player failed to update pixel values", err)
		debug.Log("warn: " + err.Error())
	}

	err = p.nt.ShowPixels()
	if err != nil {
		return errwrap.Wrap("player failed to process pixel refresh", err)
	}

	return nil
}
