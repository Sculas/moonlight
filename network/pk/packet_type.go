package pk

import (
	"errors"
	"github.com/sculas/moonlight/network/serde"
)

type Packet interface {
	// ID returns the ID of the packet.
	ID() uint32
	// Decode decodes the serde.ByteBuf into the internal data structure
	// of the packet.
	Decode(buf *serde.ByteBuf) error
	// Encode encodes the internal data structure of the packet into
	// the serde.ByteBuf.
	Encode(buf *serde.ByteBuf) error
}

var (
	Unimplemented = errors.New("this function is not implemented")
)
