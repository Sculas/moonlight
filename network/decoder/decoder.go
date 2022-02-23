package decoder

import (
	"errors"
	"fmt"
	"github.com/sculas/moonlight/network/pk"
	"github.com/sculas/moonlight/network/pk/direction"
	"github.com/sculas/moonlight/network/serde"
	"io"
)

var (
	ErrInvalidPacketLength = errors.New("invalid packet length")
)

func PacketDecoder(buf *serde.ByteBuf) (pk.Packet, bool, error) {
	// FIXME: add packet splitter implementation
	if !buf.Readable() {
		return nil, false, io.EOF
	}
	length, err := buf.ReadVarInt()
	if err != nil {
		return nil, false, err
	}
	if length <= 0 {
		return nil, false, ErrInvalidPacketLength
	}
	if buf.Len() < int(length) {
		return nil, false, ErrInvalidPacketLength
	}
	packetId, err := buf.ReadVarInt()
	// FIXME: I completely forgot about that different states use the same packet ids, so I need to sort that out
	fmt.Printf("packet id: %d\n", packetId)
	packet, ok := pk.GetPacket(direction.Serverbound, int(packetId))
	if !ok {
		return nil, true, errors.New(fmt.Sprintf("unknown packet id %d", packetId))
	}
	fmt.Printf("assuming it's packet: %T\n", packet) // I was right, this returns pk.Handshake while it's actually pk.StatusRequest
	err = packet.Decode(buf)
	if err != nil {
		return nil, false, err
	}
	/*if buf.Readable() {
		return packet, true, errors.New(fmt.Sprintf("packet %d is longer than expected: %d bytes left", packetId, buf.ReadableBytes()))
	}*/
	return packet, false, nil
}
