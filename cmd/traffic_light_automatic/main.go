package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/datoga/rpi-go-traffic-light/rpihandlers"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	auto := rpihandlers.NewAutomaticStateHandler()

	modeListener := rpihandlers.NewModeListenerHandler().WithChangeModeCallback(func(mode string) {
		log.Println("Activate mode", mode)

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

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Println("Sig received:", sig)

			auto.Stop()
			auto.Destroy()

			modeListener.Stop()
			modeListener.Destroy()

			os.Exit(0)
		}
	}()

	select {}
}
