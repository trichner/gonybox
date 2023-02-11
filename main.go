package main

import (
	"machine"
	"time"
	"trelligo/debug"
	"trelligo/dfplayer"
	"trelligo/mcu"
	"trelligo/mfrc522"
)

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.InitSerial()

	debug.Log("setup RC522")
	rc522, err := setupRC522()
	if err != nil {
		fatal()
	}

	_ = rc522

	debug.Log("setup Dfplayer")
	player, err := setupDfplayer()
	if err != nil {
		fatal()
	}

	player.SetVolume(10)

	debug.Log("play song")
	err = player.Play(2)
	if err != nil {
		fatal()
	}

	for {
		machine.LED.High()
		time.Sleep(500 * time.Millisecond)

		machine.LED.Low()
		time.Sleep(500 * time.Millisecond)

		debug.Log("check for card song")
		if rc522.IsNewCardPresent() {

			debug.Log("NEW CARD!")
			machine.LED.High()
			time.Sleep(10 * time.Millisecond)

			machine.LED.Low()
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func setupRC522() (*mfrc522.Device, error) {

	spi := machine.SPI0
	spi.Configure(machine.SPIConfig{
		SCK: machine.SPI0_SCK_PIN,
		SDO: machine.SPI0_SDO_PIN,
		SDI: machine.SPI0_SDI_PIN,
	})

	chipSelect := machine.PA07
	chipSelect.Configure(machine.PinConfig{Mode: machine.PinOutput})
	spiDriver := mfrc522.NewSpiDriver(mfrc522.NewSpi(spi, chipSelect))
	rc522Dev := mfrc522.NewDevice(spiDriver)
	err := rc522Dev.Init()
	return rc522Dev, err
}

func setupDfplayer() (*dfplayer.Player, error) {

	uart := machine.UART1
	uart.Configure(machine.UARTConfig{
		BaudRate: 9600,
		TX:       machine.D1,
		RX:       machine.D0,
	})

	rr := mcu.NewRoundTripper(uart)
	player := dfplayer.NewPlayer(rr)

	err := player.Reset()
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Millisecond * 2000)
	return player, nil
}

func fatal() {
	for {
		machine.LED.High()
		time.Sleep(100 * time.Millisecond)
		machine.LED.Low()
		time.Sleep(100 * time.Millisecond)
	}
}
