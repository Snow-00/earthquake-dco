package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/Snow-00/earthquake-dco/internal/controllers"
)

func main() {
	// create context to listen interrupt
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// load config

	// start server
	// server := &http.Server{Addr: ":8080"}

	// log.Println("Listening on :8080")
	// if err := server.ListenAndServe(); err != nil {
	// 	log.Fatal(err)
	// }

	// auto get gempa
	ticker := time.NewTicker(5 * time.Second)

out:
	for {
		select {
		case <-ticker.C:
			controllers.GetGempa()
		case <-ctx.Done():
			ticker.Stop()
			break out
		}
	}

	log.Println("Service shutdown")
}
