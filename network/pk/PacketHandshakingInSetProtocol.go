package pk

import (
	"github.com/sculas/moonlight/network/serde"
	"github.com/sculas/moonlight/server/client/state"
	"math"
)

const IDHandshakingInSetProtocol = 0x00

type HandshakingInSetProtocol struct {
	ProtocolVersion int32
	ServerAddress   string
	ServerPort      int16
	NextState       int32
}

func init() {
	RegisterPacket(state.Handshaking, IDHandshakingInSetProtocol, func() Packet {
		return &HandshakingInSetProtocol{}
	})
}

func (HandshakingInSetProtocol) ID() uint32 {
	return IDHandshakingInSetProtocol
}

func (p *HandshakingInSetProtocol) Decode(buf *serde.ByteBuf) (err error) {
	p.ProtocolVersion, err = buf.ReadVarInt()
	if err != nil {
		return
	}
	p.ServerAddress, err = buf.ReadString(math.MaxInt16)
	if err != nil {
		return
	}
	p.ServerPort, err = buf.ReadShort()
	if err != nil {
		return
	}
	p.NextState, err = buf.ReadVarInt()
	if err != nil {
		return
	}
	return
}

func (p *HandshakingInSetProtocol) Encode(*serde.ByteBuf) error {
	return Unimplemented
}
