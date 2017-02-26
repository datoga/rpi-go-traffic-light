package rpigpiowrapper

import (
	"errors"
	"log"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

const (
	PIN_GREEN = 21
	PIN_YELLOW = 20
	PIN_RED = 16
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

	var host embd.Host
	
	var rev int

	host, rev, err = embd.DetectHost()

	if err != nil {
		return err
	} 

	if host != embd.HostRPi {
		return errors.New("Error, host detected" + string(host) + ", it only works on RPi")
	}

	log.Println("Host detected:", host, "with revision number", rev)

	if err := embd.InitGPIO(); err != nil {
		return err
	}
	
	log.Println("GPIO INIT")

	if ledGreen, err = embd.NewDigitalPin(PIN_GREEN); err != nil {
		return err
	}

	if ledGreen.SetDirection(embd.Out); err != nil {
		return err
	}

	log.Println("GPIO Green enabled")

	if ledYellow, err = embd.NewDigitalPin(PIN_YELLOW); err != nil {
		return err
	}

	if ledYellow.SetDirection(embd.Out); err != nil {
		return err
	}

	log.Println("GPIO Yellow enabled")

	if ledRed, err = embd.NewDigitalPin(PIN_RED); err != nil {
		return err
	}

	if ledRed.SetDirection(embd.Out); err != nil {
		return err
	}

	log.Println("GPIO Red enabled")

	gpios := map[string]embd.DigitalPin{
		"green":  ledGreen,
		"yellow": ledYellow,
		"red":    ledRed,
	}

	rpigpio.gpios = gpios

	return nil
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
