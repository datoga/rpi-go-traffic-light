package actuators

import (
	"log"

	"github.com/datoga/rpi-go-traffic-light/mqttwrapper"
)

type RemoteActuator struct {
	mqtt *mqttwrapper.TrafficLightMQTTProxy
}

func NewRemoteActuator() *RemoteActuator {
	mqtt := mqttwrapper.NewTrafficLightMQTTProxy("remote-actuator", "remote-actuator", true)

	return &RemoteActuator{mqtt: mqtt}
}

func (remoteActuator *RemoteActuator) Start() error {
	if !remoteActuator.mqtt.IsConnected() {
		if err := remoteActuator.mqtt.Connect(); err != nil {
			return err
		}

		log.Println("Client connected")
	}

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
