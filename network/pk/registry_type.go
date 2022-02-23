package pk

import (
	"errors"
	"github.com/sculas/moonlight/network/pk/direction"
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

type PacketSupplier func() Packet

var (
	packets = map[direction.Direction]map[int]PacketSupplier{
		direction.Clientbound: {},
		direction.Serverbound: {},
	}
)

func RegisterPacket(dir direction.Direction, id int, packet PacketSupplier) {
	packets[dir][id] = packet
}

func GetPacket(dir direction.Direction, id int) (Packet, bool) {
	p, ok := packets[dir][id]
	if !ok {
		return nil, false
	}
	return p(), true
}
