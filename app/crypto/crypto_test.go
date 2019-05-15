// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package crypto

import (
	"testing"

	"github.com/va-slyusarev/pinky/app/config"
)

var (
	cfg   = &config.CryptoConfig{Salt: 0x9876543210}
	cr, _ = Instance(cfg)
)

func BenchmarkCryptographer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		cr.Hash(string(n))
	}
}

func TestCryptographer(t *testing.T) {
	tests := []struct {
		value    string
		wantHash string
	}{
		{"one", "53c2WJg3bZK"},
		{"longlonglonglonglognlonglonglonglonglonglognlonglonglonglonglonglognlonglonglonglonglonglognlonglonglonglonglonglognlong", "YkcZ3BDt9KX"},
	}
	for i, tt := range tests {
		t.Run(string(i), func(t *testing.T) {
			gotHash := cr.Hash(tt.value)
			if gotHash != tt.wantHash {
				t.Errorf("value=%s, wantHash=%s | gotHash=%s", tt.value, tt.wantHash, gotHash)
			}
		})
	}
}
