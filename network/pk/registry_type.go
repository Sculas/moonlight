package pk

import (
	"errors"
	"github.com/sculas/moonlight/network/serde"
	"github.com/sculas/moonlight/server/client/state"
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
	packets = map[state.ClientState]map[int]PacketSupplier{
		state.Handshaking: {},
		state.Play:        {},
		state.Status:      {},
		state.Login:       {},
	}
)

// RegisterPacket registers a packet with the given ID and state.
// Note that clientbound packets do not need to be registered, as they are encoded directly.
// However, serverbound packets must always be registered or else they won't be handled.
func RegisterPacket(state state.ClientState, id int, packet PacketSupplier) {
	packets[state][id] = packet
}

func GetPacket(state state.ClientState, id int) (Packet, bool) {
	p, ok := packets[state][id]
	if !ok {
		return nil, false
	}
	return p(), true
}
