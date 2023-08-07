package animations

import (
	"time"
	"trelligo/pkg/neotrellis"
)

type Animation interface {
	Update(now time.Time)
	Draw(dev *neotrellis.Device) error
}

func AnimateFor(dev *neotrellis.Device, a Animation, duration time.Duration) error {
	start := time.Now()
	now := start
	for now.Sub(start) < duration {
		now = time.Now()
		a.Update(now)
		err := a.Draw(dev)
		if err != nil {
			return err
		}
	}
	return nil
}
