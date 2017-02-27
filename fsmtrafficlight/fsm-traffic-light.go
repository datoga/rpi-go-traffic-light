package fsmtrafficlight

import (
	"container/ring"
)

type FSMTrafficLight struct {
	states *ring.Ring
}

func NewFSMTrafficLight() *FSMTrafficLight {
	trafficLightStates := []FSMTrafficLightState{Green, Yellow, Red}

	r := ring.New(len(trafficLightStates))

	for i := 0; i < r.Len(); i++ {
		r.Value = trafficLightStates[i]
		r = r.Next()
	}

	return &FSMTrafficLight{states: r}
}

func (fsm FSMTrafficLight) CurrentState() FSMTrafficLightState {
	return fsm.states.Value.(FSMTrafficLightState)
}

func (fsm FSMTrafficLight) NextState() FSMTrafficLightState {
	return fsm.states.Next().Value.(FSMTrafficLightState)
}

func (fsm *FSMTrafficLight) GoNextState() FSMTrafficLightState {
	fsm.states = fsm.states.Next()

	return fsm.states.Value.(FSMTrafficLightState)
}

func (fsm FSMTrafficLight) IsNextStateAllowed(state FSMTrafficLightState) bool {
	if fsm.states.Next().Value == state {
		return true
	} else {
		return false
	}
}
