package main

import (
	"log"
	"time"

	"github.com/datoga/rpi-go-traffic-light/actuators"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	remoteActuator := actuators.NewRemoteActuator()

	go func() {
		for {
			state, err := remoteActuator.GetRemoteState()

			if err != nil {
				log.Println(err)
			} else {
				log.Println("Current remote state:", state)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	if err := remoteActuator.Start(); err != nil {
		log.Fatalln(err)
	}

	remoteActuator.SetManual()
	remoteActuator.SetGreen()
	time.Sleep(time.Duration(7) * time.Second)
	remoteActuator.SetYellow()
	time.Sleep(time.Duration(3) * time.Second)
	remoteActuator.SetRed()
	time.Sleep(time.Duration(5) * time.Second)
	remoteActuator.SetAutomatic()
	remoteActuator.Destroy()
}
