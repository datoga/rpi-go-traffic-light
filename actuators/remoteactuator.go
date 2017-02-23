package actuators

import (
	"errors"
	"log"

	"github.com/datoga/rpi-go-traffic-light/mqttwrapper"
)

type RemoteActuator struct {
	mqtt          *mqttwrapper.TrafficLightMQTTProxy
	changeColorCb mqttwrapper.ChangeColorMQTTCallback
	remoteState   string
}

func NewRemoteActuator() *RemoteActuator {
	mqtt := mqttwrapper.NewTrafficLightMQTTProxy("remote-actuator", "remote-actuator", true)

	return &RemoteActuator{mqtt: mqtt}
}

func (remoteActuator *RemoteActuator) WithChangeColorCallback(cb mqttwrapper.ChangeColorMQTTCallback) *RemoteActuator {
	remoteActuator.setChangeColorCallback(cb)

	return remoteActuator
}

func (remoteActuator *RemoteActuator) Start() error {
	if !remoteActuator.mqtt.IsConnected() {
		if err := remoteActuator.mqtt.Connect(); err != nil {
			return err
		}

		log.Println("Client connected")
	}

	realCb := func(state string) {
		if remoteActuator.changeColorCb != nil {
			remoteActuator.changeColorCb(state)
		}

		remoteActuator.remoteState = state
	}

	if err := remoteActuator.mqtt.ListenStatusChanges(realCb); err != nil {
		return err
	}

	log.Println("Listening changes")

	return nil
}

func (remoteActuator *RemoteActuator) Destroy() {
	remoteActuator.mqtt.Disconnect()
}

func (remoteActuator *RemoteActuator) SetManual() error {
	return remoteActuator.mqtt.PublishMode("manual")
}

func (remoteActuator *RemoteActuator) SetAutomatic() error {
	return remoteActuator.mqtt.PublishMode("auto")

}

func (remoteActuator *RemoteActuator) SetGreen() error {
	return remoteActuator.mqtt.PublishState("green")

}

func (remoteActuator *RemoteActuator) SetYellow() error {
	return remoteActuator.mqtt.PublishState("yellow")
}

func (remoteActuator *RemoteActuator) SetRed() error {
	return remoteActuator.mqtt.PublishState("red")
}

func (remoteActuator *RemoteActuator) GetRemoteState() (string, error) {
	if remoteActuator.remoteState == "" {
		return "", errors.New("Not initialized state")
	}

	if !remoteActuator.mqtt.IsConnected() {
		return "", errors.New("MQTT is disconnected")
	}

	return remoteActuator.remoteState, nil
}

func (remoteActuator *RemoteActuator) setChangeColorCallback(cb mqttwrapper.ChangeColorMQTTCallback) {
	remoteActuator.changeColorCb = cb
}
