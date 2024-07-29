package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	BOT_TOKEN   string
	CHAT_ID     string
	ERR_CHAT_ID string
	DC_1        []float64
	DC_2        []float64
	DC_3        []float64
	DC_4        []float64
}

const BMKG = "https://data.bmkg.go.id/DataMKG/TEWS/autogempa.json?000"
const MAX_DIST = 200

var ENV Config

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		viper.BindEnv("BOT_TOKEN")
		viper.BindEnv("CHAT_ID")
		viper.BindEnv("ERR_CHAT_ID")
		viper.BindEnv("DC_1")
		viper.BindEnv("DC_2")
		viper.BindEnv("DC_3")
		viper.BindEnv("DC_4")
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		log.Fatal(err)
	}

	log.Println("Load config success")
}
