package rpihandlers

import (
	"container/ring"
	"log"
	"time"

	"github.com/datoga/rpi-go-traffic-light/mqttwrapper"
)

type AutomaticStateHandler struct {
	currentState    *ring.Ring
	mqtt            *mqttwrapper.TrafficLightMQTTProxy
	quitAutomaticCh chan int
	Started         bool
	stopping        bool
}

type TrafficLightState struct {
	state string
	time  int
}

func NewAutomaticStateHandler() *AutomaticStateHandler {
	mqtt := mqttwrapper.NewTrafficLightMQTTProxy("rpi-automatic-handler", "rpi-automatic-handler", false)

	greenState := TrafficLightState{"green", 10}
	yellowState := TrafficLightState{"yellow", 3}
	redState := TrafficLightState{"red", 20}

	states := []TrafficLightState{greenState, yellowState, redState}

	r := ring.New(len(states))

	for i := 0; i < r.Len(); i++ {
		r.Value = states[i]
		r = r.Next()
	}

	quitAutomaticCh := make(chan int)

	return &AutomaticStateHandler{currentState: r, mqtt: mqtt, quitAutomaticCh: quitAutomaticCh, Started: false}
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
		initialState := auto.currentState.Value.(TrafficLightState)
		stateFinishedTicker := time.NewTicker(time.Duration(initialState.time) * time.Second)

		end := false

		for end == false {
			select {
			case <-auto.quitAutomaticCh:
				stateFinishedTicker.Stop()
				end = true
			case <-stateFinishedTicker.C:
				stateFinishedTicker.Stop()
				auto.currentState = auto.currentState.Next()
				trafficLightState := auto.currentState.Value.(TrafficLightState)
				stateFinishedTicker = time.NewTicker(time.Duration(trafficLightState.time) * time.Second)
				log.Println("State changed to", trafficLightState.state)
				err := auto.mqtt.PublishState(trafficLightState.state)

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
