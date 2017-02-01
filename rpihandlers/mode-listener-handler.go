package rpihandlers

import (
	"log"

	"github.com/datoga/rpi-go-traffic-light/mqttwrapper"
)

type ModeListenerHandler struct {
	mqtt         *mqttwrapper.TrafficLightMQTTProxy
	changeModeCb mqttwrapper.ChangeModeMQTTCallback
}

func NewModeListenerHandler() *ModeListenerHandler {
	mqtt := mqttwrapper.NewTrafficLightMQTTProxy("rpi-mode-listener", "rpi-mode-listener", false)

	return &ModeListenerHandler{mqtt: mqtt}
}

func (modeListener *ModeListenerHandler) WithChangeModeCallback(cb mqttwrapper.ChangeModeMQTTCallback) *ModeListenerHandler {
	modeListener.SetChangeModeCallback(cb)

	return modeListener
}

func (modeListener *ModeListenerHandler) Start() error {
	if !modeListener.mqtt.IsConnected() {
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

func (modeListener *ModeListenerHandler) SetChangeModeCallback(cb mqttwrapper.ChangeModeMQTTCallback) {
	modeListener.changeModeCb = cb
}

func (modeListener *ModeListenerHandler) changeModeCallback(mode string) {
	log.Println("Mode changed to", mode)

	if modeListener.changeModeCb != nil {
		modeListener.changeModeCb(mode)
	}
}
