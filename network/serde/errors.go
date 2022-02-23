package serde

import "errors"

var (
	ErrInvalidVarInt     = errors.New("invalid varint")
	ErrInvalidVarLong    = errors.New("invalid varlong")
	ErrInvalidStringSize = errors.New("invalid string size")
)

func wrap(s string, err error) error {
	return errors.New(s + ": " + err.Error())
}

func wrapEOF(s string) error {
	return errors.New("EOF: " + s)
}
