package dfplayer

import (
	"errors"
)

const (
	CommandNext                   = 0x01
	CommandPrevious               = 0x02
	CommandPlayFile               = 0x03
	CommandVolumeUp               = 0x04
	CommandVolumeDown             = 0x05
	CommandSetVolume              = 0x06
	CommandSetEQ                  = 0x07
	CommandLoopPlayFile           = 0x08
	CommandSetOutputDevice        = 0x09
	CommandSleep                  = 0x0A
	CommandReset                  = 0x0C
	CommandStart                  = 0x0D
	CommandPause                  = 0x0E
	CommandPlayFolder             = 0x0F
	CommandConfigureOutputSetting = 0x10
	CommandSetLoopAll             = 0x11
	CommandPlayMP3Folder          = 0x12
	CommandAdvertiseFile          = 0x13
	CommandPlayLargeFolder        = 0x14
	CommandStopAdvertise          = 0x15
	CommandStop                   = 0x16
	CommandLoopFolder             = 0x17
	CommandRandomAll              = 0x18
	CommandSetLoop                = 0x19
	CommandSetDAC                 = 0x1A
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

func (d *Player) PlayNext() error {
	return d.sendCommand(CommandNext)
}

func (d *Player) PlayPrevious() error {
	return d.sendCommand(CommandPrevious)
}

func (d *Player) Play(file uint16) error {
	return d.sendCommandWithArg(CommandPlayFile, file)
}

func (d *Player) VolumeUp() error {
	return d.sendCommand(CommandVolumeUp)
}

func (d *Player) VolumeDown() error {
	return d.sendCommand(CommandVolumeDown)
}

// SetVolume sets the volume of the player, volume must be in the range [0,30]
func (d *Player) SetVolume(volume uint8) error {
	return d.sendCommandWithArg(CommandSetVolume, uint16(volume))
}

func (d *Player) SetEQ(eq uint8) error {
	return d.sendCommandWithArg(CommandSetEQ, uint16(eq))
}

func (d *Player) LoopFile(file uint16) error {
	return d.sendCommandWithArg(CommandLoopPlayFile, file)
}

// SetOutputDevice - undocumented
func (d *Player) SetOutputDevice(device uint16) error {
	return d.sendCommandWithArg(CommandSetOutputDevice, device)
}

func (d *Player) Sleep() error {
	return d.sendCommand(CommandSleep)
}

// Reset resets the chip
// NOTE: It is advisable to wait for 1-3 seconds after sending a reset command before sending new commands
func (d *Player) Reset() error {
	return d.sendCommand(CommandReset)
}

func (d *Player) Unpause() error {
	return d.sendCommand(CommandStart)
}

func (d *Player) Pause() error {
	return d.sendCommand(CommandPause)
}

// PlayFolder plays specific mp3 in SD:/15/004.mp3; Folder Name(1~99); File Name(1~255)
func (d *Player) PlayFolder(folder uint8, file uint8) error {
	arg := (uint16(folder) << 8) | uint16(file)
	return d.sendCommandWithArg(CommandPlayFolder, arg)
}

func (d *Player) ConfigureOutput(enable bool, gain uint8) error {
	var isEnabled uint16
	if enable {
		isEnabled = 1
	}
	arg := (isEnabled << 8) | uint16(gain)
	return d.sendCommandWithArg(CommandConfigureOutputSetting, arg)
}

func (d *Player) SetLoopAll(enable bool) error {
	var isEnabled uint16
	if enable {
		isEnabled = 1
	}
	return d.sendCommandWithArg(CommandSetLoopAll, isEnabled)
}

// PlayMP3Folder plays specific mp3 in SD:/MP3/0004.mp3; File Name(0~65535)
func (d *Player) PlayMP3Folder(file uint16) error {
	return d.sendCommandWithArg(CommandPlayMP3Folder, file)
}

// Advertise specific mp3 in SD:/ADVERT/0003.mp3; File Name(0~65535)
func (d *Player) Advertise(file uint16) error {
	return d.sendCommandWithArg(CommandAdvertiseFile, file)
}
func (d *Player) StopAdvertise() error {
	return d.sendCommand(CommandStopAdvertise)
}

// PlayLargeFolder plays specific mp3 in SD:/02/004.mp3; Folder Name(1~10); File Name(1~1000)
func (d *Player) PlayLargeFolder(folder uint8, file uint16) error {
	arg := (uint16(folder) << 12) | file
	return d.sendCommandWithArg(CommandPlayLargeFolder, arg)
}

// Stop is not documented
func (d *Player) Stop() error {
	return d.sendCommand(CommandStop)
}

// LoopFolder loops all mp3 files in folder SD:/05 where 5 is the folder provided
func (d *Player) LoopFolder(folder uint16) error {
	return d.sendCommandWithArg(CommandLoopFolder, folder)
}

// RandomAll plays all mp3 files in random order
func (d *Player) RandomAll() error {
	return d.sendCommand(CommandRandomAll)
}

func (d *Player) SetLoop(enable bool) error {
	var isEnabled uint16
	if enable {
		isEnabled = 1
	}
	return d.sendCommandWithArg(CommandSetLoop, isEnabled)
}

func (d *Player) SetDAC(enable bool) error {
	var isEnabled uint16
	if enable {
		isEnabled = 1
	}
	return d.sendCommandWithArg(CommandSetDAC, isEnabled)
}

func (d *Player) sendCommand(cmd byte) error {
	d.txBuffer.SetCommand(cmd)
	d.txBuffer.UpdateChecksum()
	err := d.roundTripper.Send(&d.txBuffer, &d.rxBuffer)
	if err != nil {
		return err
	}
	return nil
}

func (d *Player) sendCommandWithArg(cmd byte, arg uint16) error {
	d.txBuffer.SetCommand(cmd)
	d.txBuffer.SetArgument(arg)
	d.txBuffer.UpdateChecksum()
	return d.roundTripper.Send(&d.txBuffer, &d.rxBuffer)
}
