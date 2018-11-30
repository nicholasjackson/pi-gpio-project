package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nicholasjackson/periph-gpio-simulator/host/rpi"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
)

func main() {
	logger := log.New(os.Stdout, "", log.Lmicroseconds)
	logger.Println("Hello World")

	// Load all drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	go flipPinState(rpi.P1_15)

	c := make(chan os.Signal, 1)
	signal.Notify(c)

	// Block until a signal is received.
	<-c
}

func flipPinState(pin gpio.PinIO) {
	state := gpio.High

	for {
		pin.Out(state)

		time.Sleep(500 * time.Millisecond)

		// flip the state
		if state == gpio.High {
			state = gpio.Low
		} else {
			state = gpio.High
		}
	}
}
