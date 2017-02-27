# rpi-go-traffic-light
Raspberry Pi - powered traffic light control implented in Golang using MQTT

Known issues:
- The automatic state handler is unconnected from activate state, when it starts the cycle, it always overrides the current state.
- The actuator AAR file must be compiled separately and copied to the Android Studio project.
- No support for Linux different to Raspberry Pi for GPIO handler.
- FSM integration in low-level mqttwrapper
