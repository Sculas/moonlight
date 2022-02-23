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
	ErrUnknownPacket       = errors.New("unknown packet id")
)

func PacketDecoder(buf *serde.ByteBuf) (pk.Packet, error) {
	if !buf.Readable() {
		return nil, io.EOF
	}
	length, err := buf.ReadVarInt()
	if err != nil {
		return nil, err
	}
	if length <= 0 {
		return nil, ErrInvalidPacketLength
	}
	fmt.Printf("buf length: %d\n", buf.Len())
	fmt.Printf("length: %d\n", length)
	if buf.Len() < length {
		return nil, ErrInvalidPacketLength
	}
	packetId, err := buf.ReadVarInt()
	packet, ok := pk.GetPacket(direction.Serverbound, packetId)
	if !ok {
		return nil, ErrUnknownPacket
	}
	err = packet.Decode(buf)
	if err != nil {
		return nil, err
	}
	return packet, nil
}
