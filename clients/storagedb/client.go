package storagedb

import (
	"log"

	"github.com/tidwall/buntdb"
)

type BuntDBClient struct {
	db *buntdb.DB
}

type BuntDBClientOpts struct{}

var BuntDBClientOptions BuntDBClientOpts

type BuntDBClientOpt func(b *BuntDBClient)

func NewBuntDBClient(opts ...BuntDBClientOpt) *BuntDBClient {
	bunt := &BuntDBClient{}
	for _, opt := range opts {
		opt(bunt)
	}
	d, err := buntdb.Open("market.db")
	if err != nil {
		log.Panic(err)
	}
	bunt.db = d
	return bunt
}

func (b *BuntDBClient) DB() *buntdb.DB {
	return b.db
}
