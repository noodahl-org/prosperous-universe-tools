package main

import (
	"encoding/json"
	"log"
	"net/http"
	"prosperous-universe-tools/clients/analysis"
	"prosperous-universe-tools/clients/fio"
	"prosperous-universe-tools/clients/memdb"
	mdb "prosperous-universe-tools/clients/memdb"
	"prosperous-universe-tools/config"
	"prosperous-universe-tools/utils"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	conf *config.Conf
	rc   *resty.Client
	fc   *fio.FIOClient
	db   *mdb.MemDBClient
)

func init() {
	conf = config.NewConf()
	client := http.DefaultClient

	rc = resty.NewWithClient(client)
	rc.SetBaseURL(conf.BaseURL)
	rc.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"accept":       "application/json",
	})

	fc = fio.NewFIOClient(
		fio.FIOClientOptions.Config(conf),
		fio.FIOClientOptions.Client(rc),
	)

	db = mdb.NewMemDBClient(
		memdb.MemDBClientOptions.Database(
			memdb.DefaultDatabase(memdb.TableSchema()),
		),
	)
}

func main() {
	summary := []analysis.TickerSummary{}
	for {
		//data collection
		utils.Handle([]func() error{
			fc.Login,
			fc.Bids,
			fc.Orders,
		}...)
		utils.Handle(mdb.Insert("market", db.DB(), fc.Market()...))
		utils.Handle(analysis.MarketSummary(db.DB(), &summary))
		utils.Handle(print(summary))
		time.Sleep(30 * time.Minute)
	}

}

func print[T any](data []T) func() error {
	return func() error {
		bytes, err := json.MarshalIndent(data, "", "    ")
		log.Print(string(bytes))
		return err
	}
}
