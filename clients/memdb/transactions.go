package memdb

import (
	"log"

	"github.com/hashicorp/go-memdb"
)

func Insert[T any](table string, db *memdb.MemDB, rec ...T) func() error {
	return func() error {
		tx := db.Txn(true)
		for _, r := range rec {
			if err := tx.Insert(table, r); err != nil {
				return err
			}
		}
		tx.Commit()
		return nil
	}
}

func Select[T any](table, idx string, db *memdb.MemDB, out *[]T) func() error {
	return func() error {
		tx := db.Txn(false)
		defer tx.Abort()

		it, err := tx.Get(table, idx)
		if err != nil {
			return err
		}
		for obj := it.Next(); obj != nil; obj = it.Next() {
			*out = append(*out, obj.(T))
		}
		log.Printf("select results: %v", len(*out))
		return nil
	}
}