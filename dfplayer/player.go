package dfplayer

import (
	"errors"
)

const (
	positionStart            = 0
	positionVersion          = 1
	positionLength           = 2
	positionCommand          = 3
	positionFeedback         = 4
	positionQueryHighByte    = 5
	positionQueryLowByte     = 6
	positionChecksumHighByte = 7
	positionChecksumLowByte  = 8
	positionEnd              = 8
)

var ErrDeviceTimeout = errors.New("Player timed out")

type RoundTripper interface {
	Send(tx *Frame, rx *Frame) error
}

func NewPlayer(w RoundTripper) *Player {
	f := NewFrame()
	f.SetFeedback(true)
	return &Player{
		roundTripper: w,
		txBuffer:     f,
	}
}

type Player struct {
	roundTripper RoundTripper
	txBuffer     Frame
	rxBuffer     Frame
}

func (d *Player) SendCommand(cmd byte) error {
	d.txBuffer.SetCommand(cmd)
	d.txBuffer.UpdateChecksum()
	err := d.roundTripper.Send(&d.txBuffer, &d.rxBuffer)
	if err != nil {
		return err
	}
	return nil
}

func (d *Player) SendCommandWithArg(cmd byte, arg uint16) error {
	d.txBuffer.SetCommand(cmd)
	d.txBuffer.SetArgument(arg)
	d.txBuffer.UpdateChecksum()
	return d.roundTripper.Send(&d.txBuffer, &d.rxBuffer)
}
