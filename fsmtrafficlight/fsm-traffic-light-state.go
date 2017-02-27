package fsmtrafficlight

type FSMTrafficLightState string

const (
	Green  FSMTrafficLightState = "green"
	Yellow FSMTrafficLightState = "yellow"
	Red    FSMTrafficLightState = "red"
)

func (state FSMTrafficLightState) String() string {
	return string(state)
}
