// +build linux

package rpigpiowrapper

import (
	"errors"
	"log"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

type RPIGPIOPhy struct {
	gpios map[string]embd.DigitalPin
}

func NewRPIGPIOPhy() *RPIGPIOPhy {
	return &RPIGPIOPhy{}
}

func (rpigpio *RPIGPIOPhy) Prepare() error {
	var err error
	var ledGreen, ledYellow, ledRed embd.DigitalPin

	if err := embd.InitGPIO(); err != nil {
		return err
	}

	if ledGreen, err := embd.NewDigitalPin(PIN_GREEN); err != nil {
		return err
	}

	if ledGreen.SetDirection(embd.Out); err != nil {
		return err
	}

	if ledYellow, err := embd.NewDigitalPin(PIN_YELLOW); err != nil {
		return err
	}

	if ledYellow.SetDirection(embd.Out); err != nil {
		return err
	}

	if ledRed, err := embd.NewDigitalPin(PIN_RED); err != nil {
		return err
	}

	if ledRed.SetDirection(embd.Out); err != nil {
		return err
	}

	gpios := map[string]embd.DigitalPin{
		"green":  ledGreen,
		"yellow": ledYellow,
		"red":    ledRed,
	}

	rpigpio.gpios = gpios
}

func (rpigpio *RPIGPIOPhy) SetState(color string) error {
	if _, ok := rpigpio.gpios[color]; !ok {
		return errors.New("Color " + color + " not found")
	}

	for colorK, gpio := range rpigpio.gpios {
		if color == colorK {
			log.Println("GPIO", colorK, "ACTIVATED")
			gpio.Write(embd.High)
		} else {
			log.Println("GPIO", colorK, "DISABLED")
			gpio.Write(embd.Low)
		}
	}

	return nil
}

func (rpigpio *RPIGPIOPhy) Destroy() {
	embd.CloseGPIO()
}
