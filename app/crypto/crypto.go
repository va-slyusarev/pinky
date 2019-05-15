// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package crypto

import (
	"github.com/va-slyusarev/pinky/app/config"
)

// Hasher interface.
type Hasher interface {
	Hash(value string) string
}

// Instance Hasher.
func Instance(cfg *config.CryptoConfig) (Hasher, error) {
	return NewT1HA(cfg.Salt), nil
}
