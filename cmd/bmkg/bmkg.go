package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Infogempa struct{}
}

func main() {
	res, err := http.Get("https://data.bmkg.go.id/DataMKG/TEWS/autogempa.json?000")
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var respBody Response

	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(respBody.Infogempa)
}