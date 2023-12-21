package storagedb

import (
	"encoding/json"

	"github.com/tidwall/buntdb"
)

func Insert[T any](key string, db *buntdb.DB, recs []T) func() error {
	return func() error {
		err := db.Update(func(tx *buntdb.Tx) error {
			data, err := json.Marshal(recs)
			if err != nil {
				return err
			}
			_, _, err = tx.Set(key, string(data), nil)
			return err
		})
		if err != nil {
			return err
		}
		return db.Shrink()
	}
}

func Select[T any](db *buntdb.DB, out *[]T) func() error {
	return func() error {
		return db.View(func(tx *buntdb.Tx) error {
			err := tx.Ascend("", func(key, value string) bool {
				t := []T{}
				if err := json.Unmarshal([]byte(value), &t); err != nil {
					return false
				}
				*out = append(*out, t...)
				return true // continue iteration
			})
			return err
		})
	}
}

func Keys(db *buntdb.DB, out *[]string) func() error {
	return func() error {
		return db.View(func(tx *buntdb.Tx) error {
			tx.AscendKeys("", func(key, value string) bool {
				*out = append(*out, key)
				return true
			})
			return nil
		})
	}
}
