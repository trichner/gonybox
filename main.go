package main

import (
	"machine"
	"strconv"
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

	time.Sleep(3 * time.Second)

	// seesaw
	debug.Log("setup seesaw")

	debug.Log("i2c init")
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{
		SCL: machine.SCL_PIN,
		SDA: machine.SDA_PIN,
	})
	if err != nil {
		fatal(err.Error())
	}

	debug.Log("initializing neotrellis")
	nt, err := neotrellis.New(i2c, 0)
	if err != nil {
		fatal(err.Error())
	}

	debug.Log("enabling keys")
	for i := uint8(0); i < 16; i++ {
		err := nt.ConfigureKeypad(i/4, i%4, keypad.EdgeRising, true)
		if err != nil {
			fatal(err.Error())
		}
	}

	debug.Log("setup PRNG")
	prng, err := setupPRNG()
	if err != nil {
		fatal(err.Error())
	}

	nt.SetKeyHandleFunc(func(x, y uint8, e keypad.Edge) error {
		debug.Log("keypress: " + strconv.Itoa(int(x)) + "/" + strconv.Itoa(int(y)) + " (" + strconv.Itoa(int(e)) + ")")
		c := prng.Uint32()
		return nt.SetPixelColor(x, y, byte(0), byte(c>>8), byte(c>>16))
	})

	for i := uint8(0); i < 16; i++ {
		err := nt.SetPixelColor(i/4, i%4, 0, 0, 0)
		if err != nil {
			warn("setpixel " + err.Error())
		}
	}

	for {
		err = nt.ShowPixels()
		if err != nil {
			warn("showpixels " + err.Error())
		}
		time.Sleep(100 * time.Millisecond)

		err := nt.ProcessKeyEvents()
		if err != nil {
			warn("readkeys " + err.Error())
		}
	}

}

func setupPRNG() (*rand.Rand, error) {

	hi, err := machine.GetRNG()
	if err != nil {
		return nil, err
	}
	lo, err := machine.GetRNG()
	if err != nil {
		return nil, err
	}
	rsrc := rand.NewSource(int64(hi)<<32 | int64(lo))
	return rand.New(rsrc), nil
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
