package main

import (
	"fmt"
	"io"
	"time"
	"trelligo/dfplayer"
)

type LoggingWriter struct {
	w      dfplayer.RoundTripper
	logger io.Writer
}

func (l *LoggingWriter) Send(tx *dfplayer.Frame, rx *dfplayer.Frame) error {
	fmt.Fprintf(l.logger, "TX: %s\n", tx.String())

	tick := time.Now()
	err := l.w.Send(tx, rx)
	tock := time.Now()

	fmt.Fprintf(l.logger, "RX: %s\n", rx.String())
	ms := tock.Sub(tick).Milliseconds()
	fmt.Fprintf(l.logger, "Latency: %dms\n", ms)
	return err
}
