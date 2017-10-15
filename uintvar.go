// variable length binary Encoding/Marshaling/ReadWriting of integer values.
// differs from other techniques in that it uses all permutations of bytes. (up to the redundant states.)
// an encodings length is not carried in the binary data itself, so to decode the length, in bytes, is required to be known by some other means.
// some binary states are never produced as an encoding, here they are redundant and don’t have a decoding.
// the redundant states occur due to the extra information carried in the variable length.
package varbinary

import "errors"
import "fmt"
import "io"

// Uint64 enables variable length binary encoding/decoding of uint64 values.
// implements: encoding.BinaryMarshaler,encoding.BinaryUnmarshaler and io.ReadWriter.
type Uint64 uint64

var DecErr error = errors.New("Bytes do not represent an encoding.")
var BufErr error = errors.New("Encoding can not be put in supplied Buffer.")

// string representation is the hexadecimal of the encoding, (implementing Stringer)
func (x Uint64) String() string {
	b := make([]byte, 8, 8)
	n, _ := x.Read(b)
	return fmt.Sprintf("% X", b[:n])
}

// read an Uint64's encoding into a provided byte[]. (implementing io.Reader)
// BufErr returned if buffer size not big enough to contain the encoding.
func (u *Uint64) Read(b []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = BufErr
		}
	}()
	return Uint64Put(*u, b), io.EOF
}

// write an Uint64 from an encoding provided in a []byte.  (implementing io.Writer.)
// DecErr returned if value is an unused, redundant, encoding
func (u *Uint64) Write(b []byte) (int, error) {
	*u = Uint64Decoder(b...)
	if len(b) > 7 {
		if *u > Uint64(0xFEFEFEFEFEFEFEFE) {
			return 0, DecErr
		}
		return 8, io.EOF
	}
	return len(b), nil
}

func (u *Uint64) MarshalBinary() (data []byte, err error) {
	return Uint64Encoder(*u), nil
}

// return the binary encoding of a Uint64
func Uint64Encoder(u Uint64) (b []byte) {
	b=make([]byte,8,8)
	return b[:Uint64Put(u,b)] 
}

func (u *Uint64) UnmarshalBinary(data []byte) error {
	switch len(data) {
	case 0:
		*u = Uint64(0)
	case 1:
		*u = Uint64(0x1 + uint64(data[0]))
	case 2:
		*u = Uint64(0x0101 + uint64(data[0]) + uint64(data[1])<<8)
	case 3:
		*u = Uint64(0x010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16)
	case 4:
		*u = Uint64(0x01010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24)
	case 5:
		*u = Uint64(0x0101010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24 + uint64(data[4])<<32)
	case 6:
		*u = Uint64(0x010101010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24 + uint64(data[4])<<32 + uint64(data[5])<<40)
	case 7:
		*u = Uint64(0x01010101010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24 + uint64(data[4])<<32 + uint64(data[5])<<40 + uint64(data[6])<<48)
	case 8:
		*u = Uint64(0x0101010101010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24 + uint64(data[4])<<32 + uint64(data[5])<<40 + uint64(data[6])<<48 + uint64(data[7])<<56)
		if *u < Uint64(0x0101010101010101) {  // must have wrapped
			return DecErr
		}
	default:
		return DecErr
	}
	return nil
}

// return the Uint64 represented by some bytes, ignores decoding error
func Uint64Decoder(b ...byte) (u Uint64) {
	(*Uint64).UnmarshalBinary(&u, b)
	return
}

// put the encoding of a Uint64 into the provided []byte, will panic if buffer not large enough.
func Uint64Put(x Uint64, b []byte) int {
	if x == 0 {
		return 0
	}
	if x < 0x0101 {
		x -= 1
		b[0] = uint8(x)
		return 1
	}
	if x < 0x010101 {
		x -= 0x0101
		b[0] = uint8(x)
		b[1] = uint8(x >> 8)
		return 2
	}
	if x < 0x01010101 {
		x -= 0x010101
		b[0] = uint8(x)
		b[1] = uint8(x >> 8)
		b[2] = uint8(x >> 16)
		return 3
	}
	if x < 0x0101010101 {
		x -= 0x01010101
		b[0] = uint8(x)
		b[1] = uint8(x >> 8)
		b[2] = uint8(x >> 16)
		b[3] = uint8(x >> 24)
		return 4
	}
	if x < 0x010101010101 {
		x -= 0x0101010101
		b[0] = uint8(x)
		b[1] = uint8(x >> 8)
		b[2] = uint8(x >> 16)
		b[3] = uint8(x >> 24)
		b[4] = uint8(x >> 32)
		return 5
	}
	if x < 0x01010101010101 {
		x -= 0x010101010101
		b[0] = uint8(x)
		b[1] = uint8(x >> 8)
		b[2] = uint8(x >> 16)
		b[3] = uint8(x >> 24)
		b[4] = uint8(x >> 32)
		b[5] = uint8(x >> 40)
		return 6
	}
	if x < 0x0101010101010101 {
		x -= 0x01010101010101
		b[0] = uint8(x)
		b[1] = uint8(x >> 8)
		b[2] = uint8(x >> 16)
		b[3] = uint8(x >> 24)
		b[4] = uint8(x >> 32)
		b[5] = uint8(x >> 40)
		b[6] = uint8(x >> 48)
		return 7
	}
	x -= 0x0101010101010101
	b[0] = uint8(x)
	b[1] = uint8(x >> 8)
	b[2] = uint8(x >> 16)
	b[3] = uint8(x >> 24)
	b[4] = uint8(x >> 32)
	b[5] = uint8(x >> 40)
	b[6] = uint8(x >> 48)
	b[7] = uint8(x >> 56)
	return 8
}

