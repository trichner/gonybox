package main

import (
	"machine"
	"time"
	"trelligo/dfplayer"
	"trelligo/mcu"
)

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	uart := machine.UART1
	uart.Configure(machine.UARTConfig{
		BaudRate: 9600,
		TX:       machine.D10,
		RX:       machine.D11,
	})

	rr := mcu.NewRoundTripper(uart)
	player := dfplayer.NewPlayer(rr)

	err := player.Reset()
	if err != nil {
		fatal()
	}
	time.Sleep(time.Millisecond * 2000)

	player.SetVolume(10)

	err = player.Play(2)
	if err != nil {
		fatal()
	}

	for {
		machine.LED.High()
		time.Sleep(500 * time.Millisecond)

		machine.LED.Low()
		time.Sleep(500 * time.Millisecond)
	}
}

func fatal() {
	for {
		machine.LED.High()
		time.Sleep(100 * time.Millisecond)
		machine.LED.Low()
		time.Sleep(100 * time.Millisecond)
	}
}
