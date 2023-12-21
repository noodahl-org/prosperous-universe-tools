package memdb

import (
	"log"

	memdb "github.com/hashicorp/go-memdb"
)

type MemDBClient struct {
	schema *memdb.DBSchema
	db     *memdb.MemDB
}

type MemDBClientOpts struct{}

var MemDBClientOptions MemDBClientOpts

type MemDBClientOpt func(m *MemDBClient)

func NewMemDBClient(opts ...MemDBClientOpt) *MemDBClient {
	mdb := &MemDBClient{}
	for _, opt := range opts {
		opt(mdb)
	}
	return mdb
}

func (m *MemDBClient) DB() *memdb.MemDB {
	return m.db
}

func (MemDBClientOpts) Schema(s *memdb.DBSchema) MemDBClientOpt {
	return func(m *MemDBClient) {
		m.schema = s
	}
}

func (MemDBClientOpts) Database(db *memdb.MemDB) MemDBClientOpt {
	return func(m *MemDBClient) {
		m.db = db
	}
}

func DefaultDatabase(schema *memdb.DBSchema) *memdb.MemDB {
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Panic(err)
	}
	return db
}

func TableSchema() *memdb.DBSchema {
	tables := map[string]*memdb.TableSchema{
		"market": {
			Name: "market",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "ID"},
				},
				"company": {
					Name:    "company",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "CompanyID"},
				},
				"ticker": {
					Name:    "ticker",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "MaterialTicker"},
				},
				"exchange_ticker": {
					Name:    "exchange_ticker",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "ExchangeTicker"},
				},
			},
		},
		"ticker": {
			Name: "ticker",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "ID"},
				},
				"ticker": {
					Name:    "ticker",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "Ticker"},
				},
			},
		},
		"material": {
			Name: "material",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "ID"},
				},
				"ticker": {
					Name:    "ticker",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Ticker"},
				},
				"name": {
					Name:    "name",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Name"},
				},
				"category_name": {
					Name:    "category_name",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "CategoryName"},
				},
			},
		},
	}
	return &memdb.DBSchema{
		Tables: tables,
	}
}
