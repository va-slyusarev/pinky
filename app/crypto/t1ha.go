// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package crypto

import (
	"sync"

	"github.com/dgryski/go-t1ha"
)

type digest struct {
	mu   sync.Mutex
	seed uint64
	hash uint64
}

// NewT1HA t1ha implementation.
func NewT1HA(seed uint64) Hasher {
	return &digest{seed: seed}
}

// Hash return hash base58 format of string value.
func (d *digest) Hash(value string) string {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.hash = t1ha.Sum64([]byte(value), d.seed)
	return StdEncoding.Encode(d.hash)
}
