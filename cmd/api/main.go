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

	defer log.Println("this is test")

	go func() {
		for {
			ok, err := controllers.SendGempa()
			if err != nil {
				log.Fatal(err)
			}

			if ok {
				log.Println("Message sent")
			} else {
				log.Println("Not around DC")
			}

			<-ticker.C
		}
	}()

	<-ctx.Done()

	log.Println("Service shutdown")
}
