package serde

import "errors"

var (
	ErrInvalidVarInt     = errors.New("invalid varint")
	ErrInvalidVarLong    = errors.New("invalid varlong")
	ErrInvalidStringSize = errors.New("invalid string size")
)
