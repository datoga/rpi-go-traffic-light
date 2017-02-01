package main

import (
	"log"
)

type ModeListenerHandler struct {
	mqtt         *TrafficLightMQTTProxy
	changeModeCb ChangeModeMQTTCallback
}

func NewModeListenerHandler() *ModeListenerHandler {
	mqtt := NewTrafficLightMQTTProxy("rpi-mode-listener", "rpi-mode-listener")

	return &ModeListenerHandler{mqtt: mqtt}
}

func (modeListener *ModeListenerHandler) WithChangeModeCallback(cb ChangeModeMQTTCallback) *ModeListenerHandler {
	modeListener.SetChangeModeCallback(cb)

	return modeListener
}

func (modeListener *ModeListenerHandler) Start() error {
	if !modeListener.mqtt.client.IsConnected() {
		if err := modeListener.mqtt.Connect(); err != nil {
			return err
		}

		log.Println("Client connected")
	}

	if err := modeListener.mqtt.ListenModeChanges(modeListener.changeModeCallback); err != nil {
		return err
	}

	log.Println("Listening mode changes")

	return nil
}

func (modeListener *ModeListenerHandler) Stop() error {
	if err := modeListener.mqtt.UnlistenModeChanges(); err != nil {
		return err
	}

	log.Println("Unlistening mode changes")

	return nil
}

func (modeListener *ModeListenerHandler) Destroy() {
	modeListener.mqtt.Disconnect()
}

func (modeListener *ModeListenerHandler) SetChangeModeCallback(cb ChangeModeMQTTCallback) {
	modeListener.changeModeCb = cb
}

func (modeListener *ModeListenerHandler) changeModeCallback(mode string) {
	log.Println("Mode changed to", mode)

	if modeListener.changeModeCb != nil {
		modeListener.changeModeCb(mode)
	}
}
