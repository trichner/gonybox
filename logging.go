package main

import (
	"fmt"
	"io"
	"time"
	dfplayer2 "trelligo/pkg/dfplayer"
)

type LoggingWriter struct {
	w      dfplayer2.RoundTripper
	logger io.Writer
}

func (l *LoggingWriter) Send(tx *dfplayer2.Frame, rx *dfplayer2.Frame) error {
	fmt.Fprintf(l.logger, "TX: %s\n", tx.String())

	tick := time.Now()
	err := l.w.Send(tx, rx)
	tock := time.Now()

	fmt.Fprintf(l.logger, "RX: %s\n", rx.String())
	ms := tock.Sub(tick).Milliseconds()
	fmt.Fprintf(l.logger, "Latency: %dms\n", ms)
	if err != nil {
		fmt.Fprintf(l.logger, "Error: %v\n", err)
	}
	return err
}
