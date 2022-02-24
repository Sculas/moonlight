package codec

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/network/pk"
	"github.com/sculas/moonlight/network/serde"
)

func PacketEncoder(packet pk.Packet, c gnet.Conn) error {
	buf := serde.Get()
	defer serde.Put(buf)

	// TODO
	if err := packet.Encode(buf); err != nil {
		return err
	}

	return nil
}
