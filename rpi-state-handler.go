package main

import (
	"log"
)

type RPIStateHandler struct {
	mqtt *TrafficLightMQTTProxy
}

func NewRPIStateHandler() *RPIStateHandler {
	mqtt := NewTrafficLightMQTTProxy("rpi-listener", "rpi-listener")

	return &RPIStateHandler{mqtt: mqtt}
}

func (rpi *RPIStateHandler) Start() error {
	if !rpi.mqtt.client.IsConnected() {
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

	return nil
}

func (rpi *RPIStateHandler) Destroy() {
	rpi.mqtt.Disconnect()
}

func (rpi *RPIStateHandler) changeStateCallback(state string) {
	log.Println("State changed to", state)
	log.Println("Activating GPIOS")
}
