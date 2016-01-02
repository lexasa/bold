package model

import (
	"encoding/gob"
	"bytes"
)

const (
	StatusEnabled bool = true
	StatusEnabledTxt string = "enabled"

	StatusDisabled bool = false
	StatusDisabledTxt string = "disabled"

	ActionAllow bool = true
	ActionAllowTxt string = "allow"

	ActionDeny bool = false
	ActionDenyTxt string = "block"

)

func Marshal(v interface{}) ([]byte, error) {
	var stor bytes.Buffer
	enc := gob.NewEncoder(&stor)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return stor.Bytes(), nil
}

func Unmarshal(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
}
