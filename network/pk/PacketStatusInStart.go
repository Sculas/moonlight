package pk

import (
	"github.com/sculas/moonlight/network/serde"
	"github.com/sculas/moonlight/server/client/state"
)

const IDStatusInStart = 0x00

type StatusInStart struct{}

func init() {
	RegisterPacket(state.Status, IDStatusInStart, func() Packet {
		return &StatusInStart{}
	})
}

func (StatusInStart) ID() uint32 {
	return IDHandshakingInSetProtocol
}

func (p *StatusInStart) Decode(*serde.ByteBuf) error {
	return nil
}

func (p *StatusInStart) Encode(*serde.ByteBuf) error {
	return Unimplemented
}
