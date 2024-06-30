package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/Snow-00/earthquake-dco/internal/config"
	"github.com/Snow-00/earthquake-dco/internal/controllers"
)

func main() {
	// create context to listen interrupt
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// load config
	config.LoadConfig()

	// auto get gempa
	ticker := time.NewTicker(90 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			<-ticker.C

			new, ok, err := controllers.SendGempa()
			if err != nil {
				if err := controllers.AlertErr(); err != nil {
					log.Printf("Failed send alert: %s", err)
				}
				log.Fatal(err)
			}

			if new {
				log.Println("new coordinate")
			} else {
				log.Println("no new coordinate")
				continue
			}

			if ok {
				log.Println("Message sent")
			} else {
				log.Println("Not around DC")
			}
		}
	}()

	<-ctx.Done()

	log.Println("Service shutdown")
}
