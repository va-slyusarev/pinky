// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package crypto

import (
	"errors"
	"fmt"
)

// An Encoding is a radix 58 encoding/decoding scheme, defined by a
// 58-character alphabet.
type Encoding struct {
	encode    [58]byte
	decodeMap [256]int
}

// StdEncoding is the standard base58 encoding
var StdEncoding, _ = NewEncoding(encodeStd)

var encodeStd = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// NewEncoding returns a new Encoding defined by the given alphabet, which
// must be a 58-byte string that does not contain the CR / LF characters.
func NewEncoding(encoder string) (*Encoding, error) {
	if len(encoder) != 58 {
		return nil, errors.New("base58: encoding alphabet is not 58-bytes")
	}
	for i := 0; i < len(encoder); i++ {
		if encoder[i] == '\n' || encoder[i] == '\r' {
			return nil, errors.New("base58: encoding alphabet contains newline character")
		}
	}
	e := new(Encoding)
	for i := range e.decodeMap {
		e.decodeMap[i] = -1
	}
	for i := range encoder {
		e.encode[i] = byte(encoder[i])
		e.decodeMap[e.encode[i]] = i
	}
	return e, nil
}

// Encode returns encoded uint64 value to Base58 string.
func (e *Encoding) Encode(value uint64) string {
	if value == 0 {
		return string(e.encode[:1])
	}
	bin := make([]byte, 0, 8)
	for value > 0 {
		bin = append(bin, e.encode[value%58])
		value /= 58
	}

	for i, j := 0, len(bin)-1; i < j; i, j = i+1, j-1 {
		bin[i], bin[j] = bin[j], bin[i]
	}
	return string(bin)
}

// Decode returns decoded uint64 by Base58 string.
func (e *Encoding) Decode(value string) (uint64, error) {
	if value == "" {
		return 0, errors.New("base58: value should not be empty")
	}
	var n uint64
	for i := range value {
		u := e.decodeMap[value[i]]
		if u < 0 {
			return 0, fmt.Errorf("base58: invalid character - %d:%s", i, string(value[i]))
		}
		n = n*58 + uint64(u)
	}
	return n, nil
}
