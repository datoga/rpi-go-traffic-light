package rpihandlers

import (
	"log"

	"github.com/datoga/rpi-go-traffic-light/mqttwrapper"
	"github.com/datoga/rpi-go-traffic-light/rpigpiowrapper"
)

type RPIStateHandler struct {
	mqtt    *mqttwrapper.TrafficLightMQTTProxy
	rpigpio rpigpiowrapper.RPIGPIO
}

func NewRPIStateHandler() *RPIStateHandler {
	mqtt := mqttwrapper.NewTrafficLightMQTTProxy("rpi-listener", "rpi-listener", false)

	simulated := false

	var rpigpio rpigpiowrapper.RPIGPIO

	if simulated {
		rpigpio = rpigpiowrapper.NewRPIGPIOMock()
	} else {
		rpigpio = rpigpiowrapper.NewRPIGPIOPhy()
	}

	return &RPIStateHandler{mqtt: mqtt, rpigpio: rpigpio}
}

func (rpi *RPIStateHandler) Start() error {
	if err := rpi.rpigpio.Prepare(); err != nil {
		return err
	}

	if !rpi.mqtt.IsConnected() {
		if err := rpi.mqtt.Connect(); err != nil {
			return err
		}

		log.Println("Client connected")
	}

	if err := rpi.mqtt.ListenStatusChanges(rpi.changeStateCallback); err != nil {
		return err
	}

	log.Println("Listening changes")

	return nil
}

func (rpi *RPIStateHandler) Stop() error {
	if err := rpi.mqtt.UnlistenStatusChanges(); err != nil {
		return err
	}

	log.Println("Stopped")

	rpi.rpigpio.Destroy()

	return nil
}

func (rpi *RPIStateHandler) Destroy() {
	rpi.mqtt.Disconnect()
}

func (rpi *RPIStateHandler) changeStateCallback(state string) {
	log.Println("State changed to", state)
	log.Println("Activating GPIOS")

	err := rpi.rpigpio.SetState(state)

	if err != nil {
		log.Println("Error activating GPIOs:", err)
	}
}
