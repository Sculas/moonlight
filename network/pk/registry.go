package pk

import (
	"github.com/sculas/moonlight/network/pk/direction"
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
