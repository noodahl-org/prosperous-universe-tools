package fio

import (
	"prosperous-universe-tools/clients/fio/models"
	"prosperous-universe-tools/config"

	"github.com/go-resty/resty/v2"
)

type FIOClient struct {
	conf   *config.Conf
	client *resty.Client
	auth   *models.AuthResponse
	market []models.MarketData
}

type FIOClientOpts struct{}

var FIOClientOptions FIOClientOpts

type FIOClientOpt func(f *FIOClient)

// constructor
func NewFIOClient(opts ...FIOClientOpt) *FIOClient {
	fio := &FIOClient{
		market: []models.MarketData{},
	}
	for _, opt := range opts {
		opt(fio)
	}
	if fio.client == nil {
		fio.client = resty.New()
	}
	return fio
}

// options
func (FIOClientOpts) Client(h *resty.Client) FIOClientOpt {
	return func(f *FIOClient) {
		f.client = h
	}
}

func (FIOClientOpts) Config(c *config.Conf) FIOClientOpt {
	return func(f *FIOClient) {
		f.conf = c
	}
}

func (f *FIOClient) Market() []models.MarketData {
	return f.market
}
