package main

import (
	"fmt"
	"io"
	"trelligo/dfplayer"
)

type LoggingWriter struct {
	w      dfplayer.DFPlayerPort
	logger io.Writer
}

func (l *LoggingWriter) Send(tx *dfplayer.Frame, rx *dfplayer.Frame) error {
	fmt.Fprintf(l.logger, "TX: %s\n", tx.String())
	err := l.w.Send(tx, rx)
	fmt.Fprintf(l.logger, "RX: %s\n", rx.String())
	return err
}
