package main

import (
	"machine"
	"time"
	"trelligo/pkg/debug"
	"trelligo/pkg/dfplayer"
	"trelligo/pkg/dfplayer/uart"
)

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.InitSerial()

	// give some time for Serial to connect
	time.Sleep(2 * time.Second)

	debug.Log("setup ADC pin A0")

	debug.Log("setup dfplayer")
	dfp := try(setupDfplayer())

	debug.Log("setup player")

	err := dfp.SetVolume(20)
	if err != nil {
		fatal(err)
	}

	err = dfp.PlayNext()
	if err != nil {
		fatal(err)
	}

	for {
		//todo
	}

}

func setupDfplayer() (*dfplayer.Player, error) {

	uart1 := machine.UART1
	err := uart1.Configure(machine.UARTConfig{
		BaudRate: 9600,
		TX:       machine.D1,
		RX:       machine.D0,
	})
	if err != nil {
		return nil, err
	}

	rr := uart.NewRoundTripper(uart1)
	player := dfplayer.NewPlayer(rr)

	err = player.Reset()
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second * 3)
	return player, nil
}

func fatal(s error) {
	debug.Log("FATAL: " + s.Error())
	for {
		machine.LED.High()
		time.Sleep(100 * time.Millisecond)
		machine.LED.Low()
		time.Sleep(100 * time.Millisecond)
	}
}

func try[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
