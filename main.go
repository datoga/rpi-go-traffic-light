package main

import (
	"container/ring"
	"log"
	"time"
)

var trafficLightProxy *TrafficLightProxy

var quitAutomaticCh = make(chan int)

func main() {
	trafficLightProxy = NewTrafficLightProxy()

	if err := trafficLightProxy.Connect(); err != nil {
		log.Fatalln("Error connecting", err)
	}

	defer trafficLightProxy.Disconnect()

	automaticMode()

	time.Sleep(1 * time.Minute)

	quitAutomaticCh <- 1

	time.Sleep(time.Duration(1) * time.Minute)
}

func changeColorCallback(state string) {
	log.Println("Callback state", state)
}

func manualMode() {
	if err := trafficLightProxy.ListenStatusChanges(changeColorCallback); err != nil {
		log.Println("Error listening status changes", err)
	}
}

func automaticMode() {
	if err := trafficLightProxy.UnlistenStatusChanges(); err != nil {
		log.Println("Error unlistening status changes", err)
	}

	greenState := TrafficLightState{"green", 10}
	yellowState := TrafficLightState{"yellow", 1}
	redState := TrafficLightState{"red", 20}

	automaticLoop([]TrafficLightState{greenState, yellowState, redState})
}

type TrafficLightState struct {
	state string
	time  int
}

func automaticLoop(states []TrafficLightState) {
	r := ring.New(len(states))

	for i := 0; i < r.Len(); i++ {
		r.Value = states[i]
		r = r.Next()
	}

	go func() {
		currentState := r

		initialState := currentState.Value.(TrafficLightState)
		stateFinishedTicker := time.NewTicker(time.Duration(initialState.time) * time.Second)

		end := false

		for end == false {
			select {
			case <-quitAutomaticCh:
				stateFinishedTicker.Stop()
				end = true
			case <-stateFinishedTicker.C:
				stateFinishedTicker.Stop()
				currentState = currentState.Next()
				trafficLightState := currentState.Value.(TrafficLightState)
				stateFinishedTicker = time.NewTicker(time.Duration(trafficLightState.time) * time.Second)
				log.Println("State changed to", trafficLightState.state)
			}
		}

		log.Println("Exiting automatic mode")
	}()
}
