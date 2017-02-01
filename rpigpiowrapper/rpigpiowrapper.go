package rpigpiowrapper

import (
	"errors"
	"log"
)

type RPIGPIOWrapper struct {
	gpios map[string]string
}

func NewRPIGPIOWrapper() *RPIGPIOWrapper {
	gpios := map[string]string{
		"green":  "gpioGREEN",
		"yellow": "gpioYELLOW",
		"red":    "gpioRED",
	}

	return &RPIGPIOWrapper{gpios: gpios}
}

func (rpigpio *RPIGPIOWrapper) SetState(color string) error {
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
