package rpihandlers

import (
	"log"
	"time"

	"github.com/datoga/rpi-go-traffic-light/mqttwrapper"

	"github.com/datoga/rpi-go-traffic-light/fsmtrafficlight"
)

type trafficLightStateMap map[fsmtrafficlight.FSMTrafficLightState]time.Duration

type AutomaticStateHandler struct {
	mqtt            *mqttwrapper.TrafficLightMQTTProxy
	quitAutomaticCh chan int
	Started         bool
	stopping        bool
	states          *fsmtrafficlight.FSMTrafficLight
	stateMap        trafficLightStateMap
}

func NewAutomaticStateHandler() *AutomaticStateHandler {
	mqtt := mqttwrapper.NewTrafficLightMQTTProxy("rpi-automatic-handler", "rpi-automatic-handler", false)

	states := fsmtrafficlight.NewFSMTrafficLight()

	trafficLightStateMap := trafficLightStateMap{
		fsmtrafficlight.Green:  time.Duration(10) * time.Second,
		fsmtrafficlight.Yellow: time.Duration(3) * time.Second,
		fsmtrafficlight.Red:    time.Duration(20) * time.Second,
	}

	quitAutomaticCh := make(chan int)

	return &AutomaticStateHandler{states: states, stateMap: trafficLightStateMap, mqtt: mqtt, quitAutomaticCh: quitAutomaticCh, Started: false}
}

func (auto *AutomaticStateHandler) Start() error {
	if auto.Started {
		log.Println("The handler has been started")
		return nil
	}

	auto.Started = true

	if !auto.mqtt.IsConnected() {
		if err := auto.mqtt.Connect(); err != nil {
			return err
		}

		log.Println("Client connected")
	}

	go func() {
		initialState := auto.states.CurrentState()
		stateFinishedTicker := time.NewTicker(auto.stateMap[initialState])

		end := false

		for end == false {
			select {
			case <-auto.quitAutomaticCh:
				stateFinishedTicker.Stop()
				end = true
			case <-stateFinishedTicker.C:
				stateFinishedTicker.Stop()
				nextState := auto.states.GoNextState()
				stateFinishedTicker = time.NewTicker(auto.stateMap[nextState])
				log.Println("State changed to", nextState)
				err := auto.mqtt.PublishState(nextState.String())

				if err != nil {
					log.Println(err)
				}
			}
		}

		log.Println("Exiting automatic mode")

		auto.Started = false
		auto.stopping = false
	}()

	return nil
}

func (auto *AutomaticStateHandler) AsyncStop() {
	if !auto.Started {
		log.Println("The handler is stopped")

		return
	}

	if auto.stopping {
		log.Println("The handler is stopping")
		return
	}

	auto.stopping = true

	auto.quitAutomaticCh <- 1
}

func (auto *AutomaticStateHandler) Stop() {
	auto.AsyncStop()

	for auto.Started == true {
		time.Sleep(time.Duration(10) * time.Millisecond)
	}
}

func (auto *AutomaticStateHandler) Destroy() {
	auto.mqtt.Disconnect()
}
