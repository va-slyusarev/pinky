// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package types

import (
	"reflect"
	"strconv"
	"testing"
)

func TestItem_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		ID  string
		URI string
	}{
		{ID: "id", URI: "www.ya.ru"},
		{},
	}
	for count, tt := range tests {
		t.Run(strconv.Itoa(count), func(t *testing.T) {
			i := &Item{}
			i.SetID(tt.ID)
			i.SetURI(tt.URI)

			bytes, err := i.Marshal()
			if err != nil {
				t.Errorf("Item.Marshal() error = %v", err)
				return
			}
			j := &Item{}
			err = j.Unmarshal(bytes)
			if err != nil {
				t.Errorf("Item.Unmarshal() error = %v", err)
				return
			}
			if !reflect.DeepEqual(i, j) {
				t.Errorf("Marshal/Unmarshal is broken")
			}
		})
	}
}
