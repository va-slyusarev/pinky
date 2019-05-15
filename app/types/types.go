// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package types

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
)

var validURI = regexp.MustCompile(`^[a-zA-Z]+://`)

// Item type
type Item struct {
	ID  string
	URI string
}

// SetID set id
func (i *Item) SetID(id string) {
	i.ID = id
}

// SetURI set uri
func (i *Item) SetURI(uri string) {
	proto := ""
	if !validURI.MatchString(uri) {
		proto = "http://"
	}
	i.URI = fmt.Sprintf("%s%s", proto, uri)
}

// Marshal Item using gob
func (i *Item) Marshal() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(buf).Encode(i)
	if err != nil {
		return buf.Bytes(), fmt.Errorf("marshal item is broken: %s", err)
	}
	return buf.Bytes(), nil
}

// Unmarshal Item using gob
func (i *Item) Unmarshal(binary []byte) error {
	buf := bytes.NewBuffer(binary)
	err := gob.NewDecoder(buf).Decode(i)
	if err != nil {
		return fmt.Errorf("unmarshal item is broken: %s", err)
	}
	return nil
}
