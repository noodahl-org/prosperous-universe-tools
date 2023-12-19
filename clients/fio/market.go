package fio

import (
	"log"
	"prosperous-universe-tools/clients/fio/models"
	"strings"
)

func (f *FIOClient) Orders() error {
	csv := ""
	resp, err := f.client.R().
		SetResult(&csv).
		Get("/csv/orders")
	if err != nil {
		return err
	}
	csv = string(resp.Body())
	for i, row := range strings.Split(csv, "\r\n") {
		if i == 0 || row == "" {
			continue
		}
		data := models.MarketFromCSV(row, "ask")
		if data.ID != "" {
			f.market = append(f.market, data)
		}

	}
	log.Printf("orders: %v", resp.StatusCode())
	return nil
}

func (f *FIOClient) Bids() error {
	csv := ""
	resp, err := f.client.R().
		Get("/csv/bids")
	if err != nil {
		return err
	}
	csv = string(resp.Body())
	for i, row := range strings.Split(csv, "\r\n") {
		if i == 0 || row == "" {
			continue
		}

		data := models.MarketFromCSV(row, "bid")
		if data.ID != "" {
			f.market = append(f.market, data)
		}
	}
	log.Printf("bids: %v", resp.StatusCode())
	return nil
}
