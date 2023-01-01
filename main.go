package main

import (
	"log"
	"os"
	"time"
	"trelligo/dfplayer"
	"trelligo/usbtty"
)

func main() {
	var err error
	writer := usbtty.NewUsbTty("/dev/ttyUSB0")

	lwriter := &LoggingWriter{w: writer, logger: os.Stderr}
	player := dfplayer.NewPlayer(lwriter)

	err = player.Reset()
	if err != nil {
		log.Fatal(err)
	}
	// chip needs some time to actually reset, even though ACK comes faster :/
	time.Sleep(time.Millisecond * 1000)

	// volume
	err = player.SetVolume(15)
	if err != nil {
		log.Fatal(err)
	}
	//IMPORTANT: it seems the chip needs a bit to actually set the volume
	//time.Sleep(time.Millisecond * 300)

	//play file 2
	err = player.Play(2)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 3)
	err = player.Play(1)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 3)
	err = player.Stop()
	if err != nil {
		log.Fatal(err)
	}

	////play
	//err = player.sendCommand(0x0d)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//time.Sleep(time.Millisecond * 100)

}
