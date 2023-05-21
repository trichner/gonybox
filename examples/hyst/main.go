package main

import (
	"machine"
	"strconv"
	"time"
	"trelligo/pkg/debug"
	"trelligo/pkg/hyst"
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

	for {
		v, updated := h.Get()
		if updated {
			debug.Log("v=" + strconv.Itoa(v))
		}
	}

}
