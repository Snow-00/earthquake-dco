package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Gempa struct {
	Tanggal     string `json:"Tanggal"`
	Jam         string `json:"Jam"`
	Coordinates string `json:"Coordinates"`
	Magnitude   string `json:"Magnitude"`
	Kedalaman   string `json:"Kedalaman"`
	Wilayah     string `json:"Wilayah"`
	Potensi     string `json:"Potensi"`
	Dirasakan   string `json:"Dirasakan"`
	Shakemap    string `json:"Shakemap"`
}

type Response struct {
	Infogempa struct {
		Gempa Gempa `json:"gempa"`
	} `json:"Infogempa"`
}

// this is link for get image
// const image = "https://data.bmkg.go.id/DataMKG/TEWS/${Shakemap}?000"

func main() {
	resp, err := http.Get("https://data.bmkg.go.id/DataMKG/TEWS/autogempa.json?000")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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
