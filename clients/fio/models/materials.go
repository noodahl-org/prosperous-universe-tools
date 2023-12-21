package models

import (
	_ "embed"
	"encoding/json"
)

//go:embed static/materials.json
var materialBytes []byte

type Material struct {
	ID           string  `json:"ID"`
	CategoryName string  `json:"CategoryName"`
	CategoryID   string  `json:"CategoryId"`
	Name         string  `json:"Name"`
	Ticker       string  `json:"Ticker"`
	Weight       float64 `json:"Weight"`
	Volume       float64 `json:"Volume"`
}

func StaticMaterialList() []Material {
	out := []Material{}
	json.Unmarshal(materialBytes, &out)
	return out
}
