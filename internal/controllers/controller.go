package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Snow-00/earthquake-dco/internal/config"
	"github.com/Snow-00/earthquake-dco/internal/models"
)

var EQ_POINT = [2]float64{-7.21, 107.66}

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

func CompareDist(dcPoint []float64, lat, long float64) bool {
	// calculate distance
	z := math.Sin(dcPoint[0])*math.Sin(lat) + math.Cos(dcPoint[0])*math.Cos(lat)*math.Cos(long-dcPoint[1])
	dist := math.Acos(z) * 6371 // earth radius

	return dist < config.MAX_DIST
}

func SendMessage(respGempa *models.RespGempa) (*models.RespMessage, error) {
	// prepare message
	imgUrl := strings.ReplaceAll(respGempa.Infogempa.Gempa.Shakemap, ".", "%2E")
	teleUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", config.ENV.BOT_TOKEN)
	bmkgImg := fmt.Sprintf("https://data.bmkg.go.id/DataMKG/TEWS/%s?000", imgUrl)

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
		ChatID:  config.ENV.CHAT_ID,
		Photo:   bmkgImg,
		Caption: text,
	}

	// marshal message
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

func SendGempa(recheck bool) (new, ok bool, err error) {
	// get gempa info
	respGempa, err := GetGempa()
	if err != nil {
		return false, false, err
	}

	// convert string to coordinate
	coordinate := strings.Split(respGempa.Infogempa.Gempa.Coordinates, ",")
	lat, _ := strconv.ParseFloat(coordinate[0], 64)
	long, _ := strconv.ParseFloat(coordinate[1], 64)

	// check for new eq info
	if EQ_POINT[0] == lat && EQ_POINT[1] == long && !recheck {
		return false, false, nil
	}

	EQ_POINT[0] = lat
	EQ_POINT[1] = long

	// convert to radian
	lat = lat * math.Pi / 180
	long = long * math.Pi / 180

	// compare distance
	if !CompareDist(config.DC_COORDS[0], lat, long) && !CompareDist(config.DC_COORDS[1], lat, long) && !CompareDist(config.DC_COORDS[2], lat, long) && !CompareDist(config.DC_COORDS[3], lat, long) {
		return true, false, nil
	}

	// send message
	respMsg, err := SendMessage(respGempa)
	if err != nil {
		return true, false, err
	}

	if !respMsg.Ok {
		err = fmt.Errorf("status: %d; description: %s", respMsg.ErrorCode, respMsg.Description)
		return true, false, err
	}

	return true, true, nil
}

func AlertErr(typeMsg string) error {
	// prepare message
	var text string
	teleUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.ENV.BOT_TOKEN)

	if typeMsg == "info" {
		text = fmt.Sprint("Recheck: ", EQ_POINT)
	} else {
		text = "Service gempa bumi gagal"
	}

	msg := struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{ChatID: config.ENV.ERR_CHAT_ID, Text: text}

	// marshall message
	reqJSON, _ := json.Marshal(msg)
	reqBody := bytes.NewReader(reqJSON)

	// send alert to tele
	_, err := http.Post(teleUrl, "application/json", reqBody)

	return err
}
