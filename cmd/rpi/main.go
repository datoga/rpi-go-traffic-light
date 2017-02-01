package main

import (
	"log"
	"time"

	"github.com/datoga/rpi-go-traffic-light/rpihandlers"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {

	rpi := rpihandlers.NewRPIStateHandler()

	if err := rpi.Start(); err != nil {
		log.Fatalln(err)
	}

	defer rpi.Stop()
	defer rpi.Destroy()

	auto := rpihandlers.NewAutomaticStateHandler()

	modeListener := rpihandlers.NewModeListenerHandler().WithChangeModeCallback(func(mode string) {
		log.Println("Change mode to", mode)

		switch mode {
		case "auto":
			if err := auto.Start(); err != nil {
				log.Fatalln(err)
			}

		case "manual":
			auto.AsyncStop()

		default:
			log.Println("Mode", mode, "not allowed")
		}
	})

	if err := modeListener.Start(); err != nil {
		log.Fatalln(err)
	}

	defer modeListener.Stop()
	defer modeListener.Destroy()

	time.Sleep(5 * time.Minute)

	if auto.Started {
		auto.Stop()
	}

	auto.Destroy()
}
