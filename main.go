package main

import (
	"log"
	"os"
	"os/signal"

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

	rpi.P1_15.Out(gpio.High)

	c := make(chan os.Signal, 1)
	signal.Notify(c)

	// Block until a signal is received.
	<-c
}
