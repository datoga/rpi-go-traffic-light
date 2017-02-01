package main

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const MQTT_TOPIC_STATUS = "traffic-light/status"
const MQTT_TOPIC_MODE = "traffic-light/mode"

type ChangeColorMQTTCallback func(string)
type ChangeModeMQTTCallback func(string)

type TrafficLightMQTTProxy struct {
	User          string
	client        MQTT.Client
	changeColorCb ChangeColorMQTTCallback
	changeModeCb  ChangeModeMQTTCallback
}

func NewTrafficLightMQTTProxy(user string, password string) *TrafficLightMQTTProxy {

	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://m21.cloudmqtt.com:17122")
	opts.SetClientID(user)
	opts.SetUsername(user)
	opts.SetPassword(password)

	opts.SetAutoReconnect(true)

	client := MQTT.NewClient(opts)

	mqtt := TrafficLightMQTTProxy{client: client, User: user}

	return &mqtt
}

func (mqtt *TrafficLightMQTTProxy) Connect() error {
	if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Client", mqtt.User, "connected")

	return nil
}

func (mqtt *TrafficLightMQTTProxy) ListenStatusChanges(cb ChangeColorMQTTCallback) error {
	if token := mqtt.client.Subscribe(MQTT_TOPIC_STATUS, 0, mqtt.statusChangedHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Client", mqtt.User, "subscribed to status")

	mqtt.setChangeColorCallback(cb)

	return nil
}

func (mqtt *TrafficLightMQTTProxy) ListenModeChanges(cb ChangeModeMQTTCallback) error {
	if token := mqtt.client.Subscribe(MQTT_TOPIC_MODE, 0, mqtt.modeChangedHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Client", mqtt.User, "subscribed to mode")

	mqtt.setChangeModeCallback(cb)

	return nil
}

func (mqtt *TrafficLightMQTTProxy) UnlistenStatusChanges() error {
	if token := mqtt.client.Unsubscribe(MQTT_TOPIC_STATUS); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Client", mqtt.User, "unsubscribed to status")

	mqtt.setChangeColorCallback(nil)

	return nil
}

func (mqtt *TrafficLightMQTTProxy) UnlistenModeChanges() error {
	if token := mqtt.client.Unsubscribe(MQTT_TOPIC_MODE); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Client", mqtt.User, "unsubscribed to mode")

	mqtt.setChangeModeCallback(nil)

	return nil
}

func (mqtt *TrafficLightMQTTProxy) PublishState(state string) error {
	if token := mqtt.client.Publish(MQTT_TOPIC_STATUS, 0, true, state); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("State published:", state)

	return nil
}

func (mqtt *TrafficLightMQTTProxy) Disconnect() {
	mqtt.client.Disconnect(0)

	log.Println("Client", mqtt.User, "disconnected")
}

func (mqtt *TrafficLightMQTTProxy) setChangeColorCallback(cb ChangeColorMQTTCallback) {
	mqtt.changeColorCb = cb
}

func (mqtt *TrafficLightMQTTProxy) statusChangedHandler(client MQTT.Client, msg MQTT.Message) {
	newColor := string(msg.Payload())
	log.Println("New color:", newColor)

	if mqtt.changeColorCb != nil {
		mqtt.changeColorCb(newColor)
	}
}

func (mqtt *TrafficLightMQTTProxy) setChangeModeCallback(cb ChangeModeMQTTCallback) {
	mqtt.changeModeCb = cb
}

func (mqtt *TrafficLightMQTTProxy) modeChangedHandler(client MQTT.Client, msg MQTT.Message) {
	newMode := string(msg.Payload())
	log.Println("New mode:", newMode)

	if mqtt.changeModeCb != nil {
		mqtt.changeModeCb(newMode)
	}
}
