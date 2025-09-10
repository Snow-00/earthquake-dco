package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/Snow-00/earthquake-dco/internal/config"
	"github.com/Snow-00/earthquake-dco/internal/controllers"
)

func main() {
	// load config
	config.LoadConfig()

	http.HandleFunc("/trigger_check", controllers.TriggerCheck)

	log.Println("Service started")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Service shutdown")
			return
		}
		log.Fatal("Something went wrong:", err)
	}

	log.Println("Service shutdown")
}
