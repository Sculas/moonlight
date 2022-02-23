package pk

import (
	"github.com/sculas/moonlight/network/pk/direction"
	"github.com/sculas/moonlight/network/serde"
	"math"
)

const IDHandshake = 0x00

type Handshake struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      int16
	NextState       int32
}

func init() {
	RegisterPacket(direction.Serverbound, IDHandshake, func() Packet {
		return &Handshake{}
	})
}

func (Handshake) ID() uint32 {
	return IDHandshake
}

func (h *Handshake) Decode(buf *serde.ByteBuf) (err error) {
	h.ProtocolVersion, err = buf.ReadVarInt()
	if err != nil {
		return
	}
	h.ServerAddress, err = buf.ReadString(math.MaxInt16)
	if err != nil {
		return
	}
	h.ServerPort, err = buf.ReadShort()
	if err != nil {
		return
	}
	h.NextState, err = buf.ReadVarInt()
	if err != nil {
		return
	}
	return
}

func (h *Handshake) Encode(*serde.ByteBuf) error {
	return Unimplemented
}
