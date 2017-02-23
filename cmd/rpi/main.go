package main

import (
	"flag"
	"log"
	"time"

	"github.com/datoga/rpi-go-traffic-light/rpihandlers"
)

var timeToFinish = flag.Int("time", 0, "Time (minutes) to block until finish (0 to unfinished)")

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()
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

	if *timeToFinish <= 0 {
		select {} // blocks eternally
	} else {
		time.Sleep(time.Duration(*timeToFinish) * time.Minute)

		if auto.Started {
			auto.Stop()
		}

		auto.Destroy()
	}
}
