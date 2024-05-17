package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/Snow-00/earthquake-dco/internal/config"
	"github.com/Snow-00/earthquake-dco/internal/models"
)

func GetGempa() (*models.RespGempa, error) {
	// get data from bmkg
	resp, err := http.Get(config.BMKG)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read n unmarshall response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respData models.RespGempa

	err = json.Unmarshal(body, &respData)

	return &respData, err
}

func Distance(dcPoint []float64, lat, long float64) bool {
	// convert to radian
	dcPoint[0] = dcPoint[0] * math.Pi / 180
	dcPoint[1] = dcPoint[1] * math.Pi / 180
	lat = lat * math.Pi / 180
	long = long * math.Pi / 180

	// calculate distance
	z := math.Sin(dcPoint[0])*math.Sin(lat) + math.Cos(dcPoint[0])*math.Cos(lat)*math.Cos(long-dcPoint[1])
	dist := math.Acos(z) * 6371 // earth radius

	return dist < config.MAX_DIST
}

func SendMessage(respGempa *models.RespGempa) (*models.RespMessage, error) {
	// prepare message
	teleUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.ENV.BOT_TOKEN)

	text := fmt.Sprintf(
		`Dear All,
		Berikut kami informasikan gempa terbaru berdasarkan data BMKG:
		
		%s | %s
		Wilayah: %s
		Magnitude: %s SR
		Kedalaman: %s
		Potensi: %s`,
		respGempa.Infogempa.Gempa.Tanggal,
		respGempa.Infogempa.Gempa.Jam,
		respGempa.Infogempa.Gempa.Wilayah,
		respGempa.Infogempa.Gempa.Magnitude,
		respGempa.Infogempa.Gempa.Kedalaman,
		respGempa.Infogempa.Gempa.Potensi,
	)

	msg := &models.Message{
		ChatID: config.ENV.CHAT_ID,
		Text:   text,
	}

	reqJSON, _ := json.Marshal(msg)
	reqBody := bytes.NewReader(reqJSON)

	// send message
	resp, err := http.Post(teleUrl, "application/json", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respObj models.RespMessage

	err = json.Unmarshal(bodyBytes, &respObj)

	return &respObj, err
}

func SendGempa() error {
	// get gempa info
	respGempa, err := GetGempa()
	if err != nil {
		return err
	}

	// convert string to coordinate
	// coordinate := strings.Split(respGempa.Infogempa.Gempa.Coordinates, ",")
	// lat, _ := strconv.ParseFloat(coordinate[0], 64)
	// long, _ := strconv.ParseFloat(coordinate[1], 64)

	// // compare distance
	// dist := Distance(config.ENV.MBCA, lat, long)

	// send message
	_, err = SendMessage(respGempa)

	return err
}
