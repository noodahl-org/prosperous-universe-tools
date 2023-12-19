package config

import (
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed config.json
var confBytes []byte

type Conf struct {
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewConf() *Conf {
	out := Conf{}
	if err := json.Unmarshal(confBytes, &out); err != nil {
		log.Panic(err)
	}
	return &out
}
