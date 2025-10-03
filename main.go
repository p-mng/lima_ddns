package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const VERSION = "0.1.0"

func main() {
	log.Printf("lima_ddns v%s", VERSION)
	conf, err := LoadConfig()
	if err != nil {
		log.Printf("Error loading configuration: %s", err.Error())
		os.Exit(1)
	}
	log.Print("Loaded configuration")

	log.Print("Starting initial update")
	err = Update(conf)
	if err != nil {
		log.Printf("Error during initial update: %s", err.Error())
	}

	log.Print("Starting update loop")

	ticker := time.NewTicker(time.Second * time.Duration(conf.Interval))
	done := make(chan bool)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			select {
			case <-done:
				log.Print("Shutdown requested, exiting")
				done <- true
				return
			case <-ticker.C:
				log.Print("Updating records")
				err := Update(conf)
				if err != nil {
					log.Printf("Error updating records: %s", err.Error())
				}
			}
		}
	}()

	go func() {
		<-sig
		log.Print("Received interrupt")
		done <- true
	}()

	<-done
	log.Print("Goodbye")

}
