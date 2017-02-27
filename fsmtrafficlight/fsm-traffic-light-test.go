package fsmtrafficlight

import (
	"testing"
)

func TestInitialState(t *testing.T) {
	fsm := NewFSMTrafficLight()

	if fsm.CurrentState() != Green {
		t.Errorf("Error, expected state", Green, "got", fsm.CurrentState())
	}
}

func TestGetNextState(t *testing.T) {
	fsm := NewFSMTrafficLight()

	nextState := fsm.NextState()

	if nextState != Yellow {
		t.Errorf("Error, expected state", Yellow, "got", nextState)
		return
	}

	if fsm.CurrentState() != Green {
		t.Errorf("Current state modified, expected", Green, "got", fsm.CurrentState())
	}
}

func TestGoNextState(t *testing.T) {
	fsm := NewFSMTrafficLight()

	nextState := fsm.GoNextState()

	if nextState != Yellow {
		t.Errorf("Error, expected state", Yellow, "got", nextState)
		return
	}

	if fsm.CurrentState() != Yellow {
		t.Errorf("Current state unexpected, expecting", Yellow, "got", fsm.CurrentState())
	}
}

func TestIsAllowed(t *testing.T) {
	fsm := NewFSMTrafficLight()

	if !fsm.IsNextStateAllowed(Yellow) {
		t.Errorf("Error, expected next state allowed is yellow")
		return
	}

	if fsm.IsNextStateAllowed(Green) {
		t.Errorf("Error, expected next state not allowed green")
		return
	}

	if fsm.IsNextStateAllowed(Red) {
		t.Errorf("Error, expected next state not allowed red")
		return
	}

	_ = fsm.GoNextState()

	if !fsm.IsNextStateAllowed(Red) {
		t.Errorf("Error, expected next state allowed is red")
		return
	}

	if fsm.IsNextStateAllowed(Green) {
		t.Errorf("Error, expected next state not allowed green")
		return
	}

	if fsm.IsNextStateAllowed(Yellow) {
		t.Errorf("Error, expected next state not allowed yellow")
		return
	}
}
