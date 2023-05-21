package main

import (
	"machine"
	"time"
	"trelligo/pkg/debug"
	"trelligo/pkg/dfplayer"
	"trelligo/pkg/dfplayer/uart"
	"trelligo/pkg/hyst"
	"trelligo/pkg/neotrellis"
	"trelligo/pkg/neotrellis/animations"
	"trelligo/pkg/player"
	"trelligo/pkg/prng"
)

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.InitSerial()

	// give some time for Serial to connect
	time.Sleep(2 * time.Second)

	debug.Log("setup ADC pin A0")
	machine.InitADC()
	adc0 := machine.ADC{Pin: machine.A0}
	adc0.Configure(machine.ADCConfig{
		Reference:  0,
		Resolution: 12,
	})

	debug.Log("setup hysteresis")
	h := hyst.New(adc0, 1500)

	debug.Log("setup dfplayer")
	dfp := try(setupDfplayer())

	debug.Log("setup NeoTrellis")
	nt := try(setupNeoTrellis())

	time.Sleep(time.Second * 1)

	debug.Log("blink a bit")
	var err error
	r := try(prng.NewDefault())
	a := animations.NewRandomBlink(r)
	err = animations.AnimateFor(nt, a, time.Second*5)
	if err != nil {
		panic(err)
	}
	a2 := animations.NewInfinityRainbow()
	err = animations.AnimateFor(nt, a2, time.Second*5)
	if err != nil {
		panic(err)
	}

	debug.Log("setup player")
	p := try(player.New(nt, dfp, h))

	for {
		err := p.Process()
		if err != nil {
			panic(err)
		}
	}

}

func setupNeoTrellis() (*neotrellis.Device, error) {

	debug.Log("i2c init")
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{
		SCL: machine.SCL_PIN,
		SDA: machine.SDA_PIN,
	})
	if err != nil {
		return nil, err
	}

	debug.Log("initializing neotrellis")
	nt := try(neotrellis.New(i2c, 0))

	return nt, nil
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

func fatal(s string) {
	debug.Log("FATAL: " + s)
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
