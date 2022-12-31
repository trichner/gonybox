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
	dfplayer := dfplayer.NewDFPlayer(lwriter)

	//reset
	//time.Sleep(time.Millisecond * 100)
	err = dfplayer.SendCommand(0x0c)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 2000)

	// volume
	err = dfplayer.SendCommandWithArg(0x06, 10)
	if err != nil {
		log.Fatal(err)
	}
	//IMPORTANT: it seems the cheep needs a bit to actually set the volume
	time.Sleep(time.Millisecond * 300)

	//song 1
	err = dfplayer.SendCommandWithArg(0x03, 2)
	if err != nil {
		log.Fatal(err)
	}

	////play
	//err = dfplayer.SendCommand(0x0d)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//time.Sleep(time.Millisecond * 100)

}
