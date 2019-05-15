// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package crypto

import (
	"strconv"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		value uint64
	}{
		{0}, {1}, {1111}, {45645645645},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			e := StdEncoding.Encode(tt.value)
			d, err := StdEncoding.Decode(e)

			if err != nil || tt.value != d {
				t.Errorf("decode value - %d, got value - %d", d, tt.value)
				return
			}
		})
	}
}
