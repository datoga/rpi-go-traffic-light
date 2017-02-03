package rpigpiowrapper

import (
	"errors"
	"log"
)

type RPIGPIOMock struct {
	gpios map[string]string
}

func NewRPIGPIOMock() *RPIGPIOMock {

	return &RPIGPIOMock{}
}

func (rpigpio *RPIGPIOMock) Prepare() error {
	gpios := map[string]string{
		"green":  "gpioGREEN",
		"yellow": "gpioYELLOW",
		"red":    "gpioRED",
	}

	rpigpio.gpios = gpios

	return nil
}

func (rpigpio *RPIGPIOMock) SetState(color string) error {
	if _, ok := rpigpio.gpios[color]; !ok {
		return errors.New("Color " + color + " not found")
	}

	for colorK, gpio := range rpigpio.gpios {
		if color == colorK {
			log.Println("GPIO", gpio, "ACTIVATED")
		} else {
			log.Println("GPIO", gpio, "DISABLED")
		}
	}

	return nil
}

func (rpigpio *RPIGPIOMock) Destroy() {
	log.Println("Destroying")
}
