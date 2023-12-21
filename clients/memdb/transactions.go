package memdb

import (
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

func Select[T any](table, idx string, args interface{}, db *memdb.MemDB, out *[]T) func() error {
	return func() error {
		tx := db.Txn(false)
		defer tx.Abort()
		var (
			it  memdb.ResultIterator
			err error
		)
		if args == nil {
			it, err = tx.Get(table, idx)
		} else {
			it, err = tx.Get(table, idx, args)
		}

		if err != nil {
			return err
		}
		for obj := it.Next(); obj != nil; obj = it.Next() {
			*out = append(*out, obj.(T))
		}
		return nil
	}
}

func SelectOne[T any](table, idx string, args interface{}, db *memdb.MemDB, out *T) func() error {
	return func() error {
		tx := db.Txn(false)

		it, err := tx.First(table, idx, args)
		if err != nil {
			return err
		}
		if it != nil {
			*out = it.(T)
		}

		return nil
	}
}
