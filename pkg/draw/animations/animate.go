package animations

import (
	"time"
	"trelligo/pkg/draw"
)

func AnimateFor(display draw.Display, a draw.Animation, duration time.Duration) error {
	start := time.Now()
	now := start
	for now.Sub(start) < duration {
		now = time.Now()
		a.Update(now)
		err := a.Draw(display)
		if err != nil {
			return err
		}
	}
	return nil
}
