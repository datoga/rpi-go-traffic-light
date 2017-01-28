package main

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type ChangeColorCallback func(string)

type TrafficLightProxy struct {
	mqtt          MQTT.Client
	changeColorCb ChangeColorCallback
}

func NewTrafficLightProxy() *TrafficLightProxy {

	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://m21.cloudmqtt.com:17122")
	opts.SetCleanSession(true)
	opts.SetClientID("wbtvsoxz")
	opts.SetUsername("rpi")
	opts.SetPassword("rpi")

	client := MQTT.NewClient(opts)

	trafficLightProxy := TrafficLightProxy{mqtt: client}

	return &trafficLightProxy
}

func (trafficLightProxy *TrafficLightProxy) Connect() error {
	if token := trafficLightProxy.mqtt.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Client connected")

	return nil
}

func (trafficLightProxy *TrafficLightProxy) ListenStatusChanges(cb ChangeColorCallback) error {
	if token := trafficLightProxy.mqtt.Subscribe("traffic-light/status", 0, trafficLightProxy.statusChangedHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Client subscribed")

	trafficLightProxy.setChangeColorCallback(cb)

	return nil
}

func (trafficLightProxy *TrafficLightProxy) UnlistenStatusChanges() error {
	if token := trafficLightProxy.mqtt.Unsubscribe("traffic-light/status"); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Client unsubscribed")

	trafficLightProxy.setChangeColorCallback(nil)

	return nil
}

func (trafficLightProxy *TrafficLightProxy) Disconnect() {
	trafficLightProxy.mqtt.Disconnect(0)
	log.Println("Client disconnected")
}

func (trafficLightProxy *TrafficLightProxy) setChangeColorCallback(cb ChangeColorCallback) {
	trafficLightProxy.changeColorCb = cb
}

func (trafficLightProxy *TrafficLightProxy) statusChangedHandler(client MQTT.Client, msg MQTT.Message) {
	newColor := string(msg.Payload())
	log.Println("New color:", newColor)

	if trafficLightProxy.changeColorCb != nil {
		trafficLightProxy.changeColorCb(newColor)
	}
}
