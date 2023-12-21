package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"prosperous-universe-tools/clients/analysis"
	"prosperous-universe-tools/clients/charm"
	"prosperous-universe-tools/clients/fio"
	"prosperous-universe-tools/clients/fio/models"
	"prosperous-universe-tools/clients/memdb"
	mdb "prosperous-universe-tools/clients/memdb"
	"prosperous-universe-tools/clients/storagedb"
	"prosperous-universe-tools/config"
	"prosperous-universe-tools/utils"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	conf *config.Conf
	rc   *resty.Client
	fc   *fio.FIOClient
	db   *mdb.MemDBClient
	sdb  *storagedb.BuntDBClient
	ui   *charm.CharmUI
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

	sdb = storagedb.NewBuntDBClient()

	ui = charm.NewCharmUI(
		charm.CharmUIOptions.MemDB(db),
	)

}

func main() {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context) {
		for {
			//data collection
			utils.Handle([]func() error{
				fc.Login,
				fc.Bids,
				fc.Orders,
			}...)
			utils.Handle(mdb.Insert("material", db.DB(), models.StaticMaterialList()...))
			utils.Handle(mdb.Insert("market", db.DB(), fc.Market()...))
			utils.Handle(analysis.MarketSummary(db.DB()))
			utils.Handle(storagedb.Insert(
				time.Now().UTC().Format("2006-01-02:03"),
				sdb.DB(),
				fc.Market(),
			))
			wg.Done()
			time.Sleep(30 * time.Minute)
			wg.Add(1)
		}
	}(ctx)
	wg.Wait()
	ui.Start()
}

func print[T any](data []T) func() error {
	return func() error {
		bytes, err := json.MarshalIndent(data, "", "    ")
		log.Print(string(bytes))
		return err
	}
}
