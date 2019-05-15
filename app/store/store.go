// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package store

import (
	"io"

	"github.com/va-slyusarev/pinky/app/config"
	"github.com/va-slyusarev/pinky/app/types"
)

// Storer interface.
type Storer interface {
	Open(cfg *config.StoreConfig) error
	Close() error
	Put(item *types.Item) error
	Get(item *types.Item) (bool, error)
	Backup(w io.Writer) error
}

// Instance Storer.
func Instance(cfg *config.StoreConfig) (Storer, error) {
	return NewBoltDB(cfg)
}
