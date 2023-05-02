package main

import (
	"machine"
	"time"
	"trelligo/pkg/debug"
	"trelligo/pkg/dfplayer"
	"trelligo/pkg/dfplayer/uart"
	"trelligo/pkg/neotrellis"
	"trelligo/pkg/seesaw/keypad"
	"trelligo/pkg/shims/rand"
)

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.InitSerial()

	time.Sleep(time.Second * 3)

	debug.Log("setup Dfplayer")
	player, err := setupDfplayer()
	if err != nil {
		fatal(err.Error())
	}

	player.SetVolume(10)

	debug.Log("play song")
	err = player.Play(2)
	if err != nil {
		fatal(err.Error())
	}

	// seesaw
	debug.Log("setup seesaw")

	debug.Log("i2c init")
	i2c := machine.I2C0
	err = i2c.Configure(machine.I2CConfig{
		SCL: machine.SCL_PIN,
		SDA: machine.SDA_PIN,
	})
	if err != nil {
		fatal(err.Error())
	}

	debug.Log("initializing neotrellis")
	nt, err := neotrellis.New(i2c)
	if err != nil {
		fatal(err.Error())
	}

	debug.Log("enabling keys")
	for i := uint8(0); i < 16; i++ {
		err := nt.ConfigureKeypad(i, keypad.EdgeRising, true)
		if err != nil {
			fatal(err.Error())
		}
	}

	hi, err := machine.GetRNG()
	if err != nil {
		fatal(err.Error())
	}
	lo, err := machine.GetRNG()
	if err != nil {
		fatal(err.Error())
	}
	rsrc := rand.NewSource(int64(hi)<<32 | int64(lo))
	prng := rand.New(rsrc)
	for {
		debug.Log("turning on neopixels")
		for i := 0; i < 16; i++ {
			c := prng.Uint32()
			err := nt.SetPixelColor(uint16(i), byte(c), byte(c>>8), byte(c>>16), 0)
			if err != nil {
				warn("setpixel " + err.Error())
			}
		}
		time.Sleep(100 * time.Millisecond)
		err = nt.ShowPixels()
		if err != nil {
			warn("showpixels " + err.Error())
		}

		debug.Log("reading keypresses")
		err := nt.ProcessKeyEvents()
		if err != nil {
			warn("readkeys " + err.Error())
		}
	}

	for {
		machine.LED.Toggle()
		time.Sleep(500 * time.Millisecond)
	}
}

func setupDfplayer() (*dfplayer.Player, error) {

	uart1 := machine.UART1
	uart1.Configure(machine.UARTConfig{
		BaudRate: 9600,
		TX:       machine.D1,
		RX:       machine.D0,
	})

	rr := uart.NewRoundTripper(uart1)
	player := dfplayer.NewPlayer(rr)

	err := player.Reset()
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Millisecond * 2000)
	return player, nil
}

func fatal(s string) {
	debug.Log("FATAL: " + s)
	for {
		machine.LED.High()
		time.Sleep(100 * time.Millisecond)
		machine.LED.Low()
		time.Sleep(100 * time.Millisecond)
	}
}

func warn(s string) {
	debug.Log("WARN: " + s)
}
