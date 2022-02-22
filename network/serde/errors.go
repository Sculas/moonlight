package serde

import "errors"

var (
	ErrVarInt            = errors.New("VarInt too big")
	ErrVarLong           = errors.New("VarLong too big")
	ErrInvalidStringSize = errors.New("invalid string size")
)
