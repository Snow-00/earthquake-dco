package main

import (
	"log"
	"net/http"
)

func main() {
	// load config

	// start server
	server := &http.Server{Addr: ":8080"}

	log.Println("Listening on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
