// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package store

import (
	"fmt"
	"io"

	"github.com/boltdb/bolt"

	"github.com/va-slyusarev/pinky/app/config"
	"github.com/va-slyusarev/pinky/app/types"
)

const (
	bucketItem = "item"
)

type boltDB struct {
	db *bolt.DB
}

// NewBoltDB create new boltdb.
func NewBoltDB(cfg *config.StoreConfig) (Storer, error) {
	bdb := new(boltDB)
	err := bdb.Open(cfg)
	return bdb, err
}

// Open new and init boltdb instance.
func (b *boltDB) Open(cfg *config.StoreConfig) error {
	db, err := bolt.Open(cfg.DBPath, 0600, &bolt.Options{Timeout: cfg.Timeout})
	if err != nil {
		return fmt.Errorf("boltdb: open storage is broken: %v", err)
	}
	if db == nil {
		return fmt.Errorf("boltdb: storage is nil")
	}

	b.db = db

	err = b.createBucket(bucketItem)
	if err != nil {
		return fmt.Errorf("boltdb: create item bucket: %v", err)
	}

	return nil
}

// Close boltdb instance.
func (b *boltDB) Close() error {
	err := b.db.Close()
	if err != nil {
		return fmt.Errorf("boltdb: close storage is broken: %s", err)
	}
	return nil
}

// Put Item in boltDB by item ID.
func (b *boltDB) Put(item *types.Item) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		marshal, err := item.Marshal()
		if err != nil {
			return fmt.Errorf("boltdb: store: %v", err)
		}
		err = tx.Bucket([]byte(bucketItem)).Put([]byte(item.ID), marshal)
		if err != nil {
			return fmt.Errorf("boltdb: store url bucket is broken: %v", err)
		}
		return nil
	})
}

// Get Item from boltdb by item ID.
func (b *boltDB) Get(item *types.Item) (bool, error) {
	binary, err := b.fromStore(bucketItem, item.ID)
	if err != nil {
		return false, fmt.Errorf("boltdb: get item bucket: %v", err)
	}
	err = item.Unmarshal(binary)
	if err != nil {
		return false, fmt.Errorf("boltdb: unmarshal item: %v", err)
	}
	return true, nil
}

// Backup boltdb instance.
func (b *boltDB) Backup(w io.Writer) error {
	err := b.db.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(w)
		return err
	})
	return err
}

// createBucket.
func (b *boltDB) createBucket(bucket string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
}

// fromStore.
func (b *boltDB) fromStore(bucket, key string) ([]byte, error) {
	var res []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		res = b.Get([]byte(key))
		if len(res) == 0 {
			return fmt.Errorf("not found")
		}
		return nil
	})
	return res, err
}
