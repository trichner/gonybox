package main

import (
	"machine"
	"strconv"
	"time"
	"trelligo/pkg/debug"
	"trelligo/pkg/dfplayer"
	"trelligo/pkg/dfplayer/uart"
	"trelligo/pkg/neotrellis"
	"trelligo/pkg/seesaw"
	"trelligo/pkg/seesaw/keypad"
	"trelligo/pkg/seesaw/neopixel"
	"trelligo/pkg/shims/rand"
)

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.InitSerial()

	time.Sleep(time.Second * 10)

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

	debug.Log("seesaw new")
	seesawDev := seesaw.New(neotrellis.DefaultNeotrellisAddress, i2c)

	debug.Log("seesaw begin")
	err = seesawDev.Begin()
	if err != nil {
		fatal(err.Error())
	}

	debug.Log("initializing neopixels")
	// https://github.com/adafruit/Adafruit_Seesaw/blob/8a2dc5e0645239cb34e23a4b62c456436b098ab3/Adafruit_NeoTrellis.cpp#L10
	const NeoTrellisSeesawPin = 3
	const nPixels = 16
	pix, err := neopixel.New(seesawDev, NeoTrellisSeesawPin, nPixels, neopixel.NEO_GRB)
	if err != nil {
		fatal(err.Error())
	}

	debug.Log("initializing keypad")
	kpd := keypad.New(seesawDev)

	debug.Log("enabling keys")
	for i := 0; i < nPixels; i++ {
		err := kpd.ConfigureKeypad(uint8(i), keypad.EdgeRising, true)
		if err != nil {
			fatal(err.Error())
		}
	}

	debug.Log("enable keypad interrupt")
	err = kpd.SetKeypadInterrupt(true)
	if err != nil {
		fatal(err.Error())
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
		for i := 0; i < nPixels; i++ {
			c := prng.Uint32()
			err := pix.SetPixelColor(uint16(i), byte(c), byte(c>>8), byte(c>>16), 0)
			if err != nil {
				fatal(err.Error())
			}
		}
		time.Sleep(100 * time.Millisecond)
		debug.Log("showing")
		err = pix.ShowPixels()
		if err != nil {
			fatal(err.Error())
		}
		debug.Log("done!")

		debug.Log("reading keypresses")
		n, err := kpd.KeyEventCount()
		if err != nil {
			fatal(err.Error())
		}
		debug.Log("events: " + strconv.Itoa(int(n)))
		events := make([]keypad.KeyEvent, n)
		err = kpd.Read(events)
		if err != nil {
			fatal(err.Error())
		}
		for _, e := range events {
			debug.Log("keypress: " + strconv.Itoa(int(e.Key())) + " (" + strconv.Itoa(int(e.Edge())) + ")")
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
