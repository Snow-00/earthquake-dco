package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Snow-00/earthquake-dco/internal/config"
	"github.com/Snow-00/earthquake-dco/internal/models"
)

// func SendMessage() (*models.RespMessage, error) {
// 	// prepare request

//		// send message with photo
//	}

func SendGempa() error {
	// get data from bmkg
	resp, err := http.Get(config.BMKG)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read n unmarshall response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respData models.Response

	err = json.Unmarshal(body, &respData)
	if err != nil {
		return err
	}

	// print info
	log.Println(respData)

	return nil
}

func Calculate() {

}
