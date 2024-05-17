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
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			if err := controllers.SendGempa(); err != nil {
				log.Fatal(err)
			}

			<-ticker.C
		}
	}()

	<-ctx.Done()

	log.Println("Service shutdown")
}
