package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type MarketData struct {
	ID             string    `json:"id"`
	MaterialTicker string    `json:"material_ticker"`
	ExchangeCode   string    `json:"exchange_code"`
	CompanyID      string    `json:"company_id"`
	CompanyName    string    `json:"company_name"`
	CompanyCode    string    `json:"company_code"`
	ItemCount      int       `json:"item_count"`
	ItemCost       float64   `json:"item_cost"`
	Type           string    `json:"type"`
	CreatedAt      time.Time `json:"created_at"`
}

func MarketFromCSV(row, t string) MarketData {
	split := strings.Split(row, ",")
	if len(split) == 8 {
		split = append(split[:4], split[5:]...)
	}
	if len(split) != 7 {
		return MarketData{}
	}
	cnt, err := strconv.Atoi(split[5])
	if err != nil {
		return MarketData{}
	}

	price, err := strconv.ParseFloat(split[6], 64)
	if err != nil {
		return MarketData{}
	}

	s := sha256.New()
	hash := base64.StdEncoding.EncodeToString(s.Sum([]byte(
		fmt.Sprintf("%s%s%s%v%s", split[0], split[1], split[2], price, t),
	)))
	return MarketData{
		ID:             hash,
		MaterialTicker: split[0],
		ExchangeCode:   split[1],
		CompanyID:      split[2],
		CompanyName:    strings.Replace(split[3], "\\", "", -1),
		CompanyCode:    split[4],
		ItemCount:      cnt,
		ItemCost:       price,
		Type:           t,
		CreatedAt:      time.Now().UTC(),
	}
}
