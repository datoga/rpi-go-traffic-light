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
	rpi := rpihandlers.NewRPIStateHandler()

	if err := rpi.Start(); err != nil {
		log.Fatalln(err)
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Println("Sig received:", sig)
			rpi.Stop()
			rpi.Destroy()

			os.Exit(0)
		}
	}()

	select {}
}
