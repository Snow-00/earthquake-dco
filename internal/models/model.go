package models

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

type RespGempa struct {
	Infogempa struct {
		Gempa Gempa `json:"gempa"`
	} `json:"Infogempa"`
}

type Message struct {
	ChatID  string `json:"chat_id"`
	Photo   string `json:"photo"`
	Caption string `json:"caption"`
}

type RespMessage struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}
