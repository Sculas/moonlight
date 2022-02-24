package codec

import (
	"errors"
	"fmt"
	"github.com/sculas/moonlight/network/pk"
	"github.com/sculas/moonlight/network/serde"
	"github.com/sculas/moonlight/server/client/state"
	"io"
)

var (
	ErrInvalidPacketLength = errors.New("invalid packet length")
	ErrPacketGroupTooLarge = errors.New("grouped packets length wider than 21-bit")
)

const (
	packetHeaderSize = 2
)

type packetHandler = func(state state.ClientState, packet pk.Packet) bool

func Decode(buf *serde.ByteBuf, state *state.ClientState, f packetHandler) error {
	if !buf.Readable() {
		return io.EOF
	}
	indexes, err := splitPackets(buf)
	if err != nil {
		return err
	}

	for _, index := range indexes {
		if !buf.SeekTo(index) {
			return errors.New("failed to seek to packet")
		}
		packetId, err := buf.ReadVarInt()
		if err != nil {
			return err
		}
		packet, ok := pk.GetPacket(*state, int(packetId))
		if !ok {
			return errors.New(fmt.Sprintf("unknown packet id %d for state %d", packetId, *state))
		}
		err = packet.Decode(buf)
		if err != nil {
			return err
		}
		if !f(*state, packet) {
			return nil
		}
		// TODO: need to come back to this later sometime, if it's actually needed..
		/*if buf.ReadableBytes()-index > 0 {
			return packets, true, errors.New(fmt.Sprintf("packet %d is longer than expected: %d bytes left", packetId, buf.ReadableBytes()))
		}*/
	}
	return nil
}

// splitPackets splits a serde.ByteBuf into multiple packets, if any.
// It returns an array of found indexes of each start of a packet.
func splitPackets(buf *serde.ByteBuf) (indexes []int, err error) {
	ri := buf.SeekIndex() // get the original seek index
	defer func() {
		buf.SeekTo(ri) // restore the original seek index
	}()
	// Notchian servers read at max 3 combined packets, it seems?
	for i := 0; i < 3; i++ {
		if !buf.Readable() || buf.ReadableBytes() < packetHeaderSize {
			return
		}

		ci := buf.SeekIndex()     // get the current seek index
		read, _ := buf.ReadByte() // read the first byte, error can never happen here due to Readable() check

		if read >= 0 {
			buf.SeekTo(ci) // restore the current seek index
			length, e := buf.ReadVarInt()
			if e != nil {
				return nil, e
			}
			if length <= 0 {
				return nil, ErrInvalidPacketLength
			}

			if buf.Len() < int(length) {
				// FIXME: This is weird, shouldn't we return an error?
				//        Notchian servers don't seem to do that, so let's just return.
				return
			}

			index := buf.SeekIndex()
			_, e = buf.Seek(int(length)) // seek to the end of the packet
			if e != nil {
				return nil, e
			}

			indexes = append(indexes, index) // found a packet start!
		}
	}
	return nil, ErrPacketGroupTooLarge // TODO
}
