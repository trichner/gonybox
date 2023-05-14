package main

import "trelligo/pkg/dfplayer"

var _ = Controls(&dfplayer.Player{})

type Controls interface {
	Pause() error
	Play(number uint16) error
	PlayNext() error
	PlayPrevious() error
	SetVolume(volume uint8) error
	Stop() error
	Unpause() error
	VolumeDown() error
	VolumeUp() error
	LoopFolder(folder uint16) error
}
