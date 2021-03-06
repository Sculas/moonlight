package serde

import (
	"encoding/binary"
	"math"
)

var e = binary.BigEndian

type ByteBuf struct {
	B []byte
	i int // reader index
}

func NewByteBuf() *ByteBuf {
	return &ByteBuf{
		B: make([]byte, 0), // TODO: default size
	}
}

func (b *ByteBuf) Len() int {
	return len(b.B)
}

func (b *ByteBuf) Bytes() []byte {
	return b.B
}

func (b *ByteBuf) Write(p []byte) {
	b.B = p
	b.i = 0
}

func (b *ByteBuf) Reset() {
	b.B = b.B[:0]
	b.i = 0
}

func (b *ByteBuf) got(i int) bool {
	return b.ReadableOffset(i)
}

func (b *ByteBuf) seek(i int) int {
	b.i += i
	return b.i
}

func (b *ByteBuf) SeekIndex() int {
	return b.i
}

func (b *ByteBuf) Seek(i int) (int, error) {
	if i < 0 || b.i+i > b.Len() {
		return 0, wrapEOF("Seek")
	}
	return b.seek(i), nil
}

func (b *ByteBuf) SeekTo(i int) bool {
	if i < 0 || i > b.Len() {
		return false
	}
	b.i = i
	return true
}

func (b *ByteBuf) Index() int {
	return b.i
}

func (b *ByteBuf) Readable() bool {
	return b.Len() > b.i
}

func (b *ByteBuf) ReadableOffset(i int) bool {
	return b.Len() >= b.i+i
}

func (b *ByteBuf) ReadableBytes() int {
	return b.Len() - b.i
}

// ReadByte reads a byte from the buffer.
// It returns a signed byte and an error when the buffer is empty.
func (b *ByteBuf) ReadByte() (int8, error) {
	if !b.got(szByte) {
		return 0, wrapEOF("Byte")
	}
	c := b.B[b.i]
	b.seek(szByte)
	return int8(c), nil
}

// FIXME: should we convert the 2 functions below to []int8 too?

func (b *ByteBuf) ReadBytes(i int) ([]byte, error) {
	if !b.got(i) {
		return nil, wrapEOF("Bytes")
	}
	c := b.B[b.i : b.i+i]
	b.seek(i)
	return c, nil
}

func (b *ByteBuf) ReadAllBytes() ([]byte, error) {
	if !b.Readable() {
		return nil, wrapEOF("AllBytes")
	}
	c := b.B[b.i:]
	b.seek(b.Len())
	return c, nil
}

func (b *ByteBuf) ReadBool() (bool, error) {
	bb, err := b.ReadByte()
	return bb != 0, err
}

func (b *ByteBuf) ReadShort() (int16, error) {
	if !b.got(szShort) {
		return 0, wrapEOF("Short")
	}
	c := e.Uint16(b.B[b.i : b.i+szShort])
	b.seek(szShort)
	return int16(c), nil
}

func (b *ByteBuf) ReadInt() (int32, error) {
	if !b.got(szInt32) {
		return 0, wrapEOF("Int")
	}
	c := e.Uint32(b.B[b.i : b.i+szInt32])
	b.seek(szInt32)
	return int32(c), nil
}

func (b *ByteBuf) ReadLong() (int64, error) {
	if !b.got(szInt64) {
		return 0, wrapEOF("Long")
	}
	c := e.Uint64(b.B[b.i : b.i+szInt64])
	b.seek(szInt64)
	return int64(c), nil
}

func (b *ByteBuf) ReadFloat() (float32, error) {
	if !b.got(szFloat32) {
		return 0, wrapEOF("Float")
	}
	c, err := b.ReadInt()
	return math.Float32frombits(uint32(c)), err
}

func (b *ByteBuf) ReadDouble() (float64, error) {
	if !b.got(szDouble) {
		return 0, wrapEOF("Double")
	}
	c, err := b.ReadLong()
	return math.Float64frombits(uint64(c)), err
}

const (
	maxVarInt        = 5
	maxVarLong       = 10
	varTermByte int8 = -128 // TIL that in Java, 0x80 byte is actually an underflow byte
)

// ReadVarInt reads a variable-length integer from the buffer.
// If the buffer is too small, EOF is returned.
// If the VarInt is larger than 5 bytes, ErrInvalidVarInt is returned.
// If an error occurs, the buffer must be discarded since the integrity can no longer be guaranteed
// since we don't know how many bytes are left until the next safe read index.
func (b *ByteBuf) ReadVarInt() (int32, error) {
	if !b.Readable() {
		return 0, wrapEOF("VarInt")
	}
	var i uint32
	maxRead := int(math.Min(maxVarInt, float64(b.ReadableBytes())))
	for j := 0; j < maxRead; j++ {
		k, err := b.ReadByte()
		if err != nil {
			return 0, wrap("VarInt", err)
		}
		i |= uint32(k&0x7F) << (7 * j)
		if (k & varTermByte) != varTermByte {
			return int32(i), nil
		}
	}
	return 0, ErrInvalidVarInt
}

// ReadVarLong reads a variable-length long from the buffer.
// If the buffer is too small, EOF is returned.
// If the VarLong is larger than 10 bytes, ErrInvalidVarLong is returned.
// If an error occurs, the buffer must be discarded since the integrity can no longer be guaranteed
// since we don't know how many bytes are left until the next safe read index.
func (b *ByteBuf) ReadVarLong() (int64, error) {
	if !b.Readable() {
		return 0, wrapEOF("VarLong")
	}
	var i uint64
	maxRead := int(math.Min(maxVarLong, float64(b.ReadableBytes())))
	for j := 0; j < maxRead; j++ {
		k, err := b.ReadByte()
		if err != nil {
			return 0, wrap("VarLong", err)
		}
		i |= uint64(k&0x7F) << (7 * j)
		if (k & varTermByte) != varTermByte {
			return int64(i), nil
		}
	}
	return 0, ErrInvalidVarLong
}

func (b *ByteBuf) ReadString(max int) (string, error) {
	if !b.Readable() {
		return "", wrapEOF("String")
	}
	sz32, err := b.ReadVarInt()
	if err != nil {
		return "", wrap("String", err)
	}
	sz := int(sz32)
	if sz < 0 || sz > max*4 {
		return "", ErrInvalidStringSize
	}
	c, err := b.ReadBytes(sz)
	if err != nil {
		return "", wrap("String", err)
	}
	if len(c) > max*4 {
		return "", ErrInvalidStringSize
	}
	return string(c), nil // causes a strcopy, but this is safer
}
