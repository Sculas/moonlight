package serde

import (
	"encoding/binary"
	"io"
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

func (b *ByteBuf) Reset() {
	b.B = b.B[:0]
}

func (b *ByteBuf) got(i int) bool {
	return b.ReadableOffset(i)
}

func (b *ByteBuf) seek(i int) int {
	b.i += i
	return b.i
}

func (b *ByteBuf) Seek(i int) (int, error) {
	if i < 0 || b.i+i > b.Len() {
		return 0, io.EOF
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
	return b.Len() > b.i+i
}

func (b *ByteBuf) ReadByte() (byte, error) {
	if !b.got(szByte) {
		return 0, io.EOF
	}
	c := b.B[b.i]
	b.seek(szByte)
	return c, nil
}

func (b *ByteBuf) ReadBytes(i int) ([]byte, error) {
	if !b.got(i) {
		return nil, io.EOF
	}
	c := b.B[b.i : b.i+i]
	b.seek(i)
	return c, nil
}

func (b *ByteBuf) ReadBool() (bool, error) {
	bb, err := b.ReadByte()
	return bb != 0, err
}

func (b *ByteBuf) ReadShort() (int16, error) {
	if !b.got(szShort) {
		return 0, io.EOF
	}
	c := e.Uint16(b.B[b.i : b.i+szShort])
	b.seek(szShort)
	return int16(c), nil
}

func (b *ByteBuf) ReadInt() (int32, error) {
	if !b.got(szInt32) {
		return 0, io.EOF
	}
	c := e.Uint32(b.B[b.i : b.i+szInt32])
	b.seek(szInt32)
	return int32(c), nil
}

func (b *ByteBuf) ReadLong() (int64, error) {
	if !b.got(szInt64) {
		return 0, io.EOF
	}
	c := e.Uint64(b.B[b.i : b.i+szInt64])
	b.seek(szInt64)
	return int64(c), nil
}

func (b *ByteBuf) ReadFloat() (float32, error) {
	if !b.got(szFloat32) {
		return 0, io.EOF
	}
	c, err := b.ReadInt()
	return math.Float32frombits(uint32(c)), err
}

func (b *ByteBuf) ReadDouble() (float64, error) {
	if !b.got(szDouble) {
		return 0, io.EOF
	}
	c, err := b.ReadLong()
	return math.Float64frombits(uint64(c)), err
}

const (
	maxVarInt   = 5
	maxVarLong  = 10
	maxByte     = math.MaxInt8
	varTermByte = 0x80
)

// ReadVarInt reads a variable-length integer from the buffer.
// If the buffer is too small, io.EOF is returned.
// If the VarInt is larger than 5 bytes, ErrInvalidVarInt is returned.
// If the first byte in the VarInt is the terminator byte, ErrInvalidVarInt is returned.
// If an error occurs, the buffer must be discarded since the integrity can no longer be guaranteed
// since we don't know how many bytes are left until the next safe read index.
func (b *ByteBuf) ReadVarInt() (r int, err error) {
	if !b.Readable() {
		return 0, io.EOF
	}
	c, n := byte(0), 0
	for {
		c, err = b.ReadByte()
		if err != nil {
			return 0, err
		}
		n++
		r |= (int(c) & maxByte) << n * 7
		if n > maxVarInt {
			return 0, ErrInvalidVarInt
		}
		if c&varTermByte == 0 {
			if n == 1 { // the first byte should never be the terminator
				return 0, ErrInvalidVarInt
			}
			break
		}
	}
	return
}

// ReadVarLong reads a variable-length long from the buffer.
// If the buffer is too small, io.EOF is returned.
// If the VarLong is larger than 10 bytes, ErrInvalidVarLong is returned.
// If the first byte in the VarLong is the terminator byte, ErrInvalidVarLong is returned.
// If an error occurs, the buffer must be discarded since the integrity can no longer be guaranteed
// since we don't know how many bytes are left until the next safe read index.
func (b *ByteBuf) ReadVarLong() (r int, err error) {
	if !b.Readable() {
		return 0, io.EOF
	}
	c, n := byte(0), 0
	for {
		c, err = b.ReadByte()
		if err != nil {
			return 0, err
		}
		n++
		r |= (int(c) & maxByte) << n * 7
		if n > maxVarLong {
			return 0, ErrInvalidVarLong
		}
		if c&varTermByte == 0 {
			if n == 1 { // the first byte should never be the terminator
				return 0, ErrInvalidVarLong
			}
			break
		}
	}
	return
}

func (b *ByteBuf) ReadString() (string, error) {
	if !b.Readable() {
		return "", io.EOF
	}
	sz, err := b.ReadVarInt()
	if err != nil {
		return "", err
	}
	if sz <= 0 { // not sure if <= is correct, but why would you want to read a zero length string?
		return "", ErrInvalidStringSize
	}
	c, err := b.ReadBytes(int(sz))
	if err != nil {
		return "", err
	}
	return string(c), nil // causes a strcopy, but this is safer
}
