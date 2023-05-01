package mfrc522

import (
	"strconv"
	"time"
	"trelligo/pkg/debug"
)

type LoggedDriver struct {
	Delegate Driver
	Start    int64
}

func (l *LoggedDriver) WriteRegister(reg Register, tx []byte) error {
	ts := int(time.Now().UnixMilli() - l.Start)
	debug.Log(strconv.Itoa(ts) + " W " + debug.FmtByteToHex(byte(reg)) + " " + debug.FmtSliceToHex(tx))
	err := l.Delegate.WriteRegister(reg, tx)
	if err != nil {
		debug.Log("write err=" + err.Error())
	}
	return err
}

func (l *LoggedDriver) ReadRegister(reg Register, rx []byte) error {
	err := l.Delegate.ReadRegister(reg, rx)
	if err != nil {
		debug.Log("read:  err=" + err.Error())
	} else {
		ts := int(time.Now().UnixMilli() - l.Start)
		debug.Log(strconv.Itoa(ts) + " R " + debug.FmtByteToHex(byte(reg)) + " " + debug.FmtSliceToHex(rx))
	}
	return err
}
