package varbinary

import "errors"
import "fmt"
import "io"

// variable length binary encoding/decoding of uint64 values.
// the length of the encoding carries information, so some binary encodings are redundant (ones high in the 8-byte range) they don’t have a decoding, due to all possible values already have an encoding thanks to that extra information.
// no gaps, uses all permutations of bytes up to the redundant values.
type Uint64 uint64

// string representation is the hexadecimal of the encoding, (implementing Stringer)
func (x Uint64) String() string {
	b := make([]byte, 8, 8)
	n, _ := x.Read(b)
	return fmt.Sprintf("% X", b[:n])
}

var bufErr error = errors.New("The Uint64's variable-length binary encoding can not be put in supplied Buffer.")

// write into a Uint64 from a variable length encoding provided in a []byte, (implementing io.Writer.)
func (u *Uint64) Write(b []byte) (int, error) {
	*u = GetUint64(b...)
	if len(b) > 7 {
		if b[7] == 0xff {
			return 8, encErr
		}
		return 8, io.EOF
	}
	return len(b), nil
}

var encErr error = errors.New("Buffer does not represent a valid binary encoding.")

// read a Uint64's encoding into a byte[]. (implementing io.Reader)
func (u Uint64) Read(b []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = bufErr
		}
	}()
	return PutUint64(b, u), io.EOF
}

// return the Uint64 represented by the []byte
func GetUint64(b ...byte) Uint64 {
	switch len(b) {
	case 0:
		return Uint64(0)
	case 1:
		return Uint64(0x1 + uint64(b[0]))
	case 2:
		return Uint64(0x0101 + uint64(b[0]) + uint64(b[1])<<8)
	case 3:
		return Uint64(0x010101 + uint64(b[0]) + uint64(b[1])<<8 + uint64(b[2])<<16)
	case 4:
		return Uint64(0x01010101 + uint64(b[0]) + uint64(b[1])<<8 + uint64(b[2])<<16 + uint64(b[3])<<24)
	case 5:
		return Uint64(0x0101010101 + uint64(b[0]) + uint64(b[1])<<8 + uint64(b[2])<<16 + uint64(b[3])<<24 + uint64(b[4])<<32)
	case 6:
		return Uint64(0x010101010101 + uint64(b[0]) + uint64(b[1])<<8 + uint64(b[2])<<16 + uint64(b[3])<<24 + uint64(b[4])<<32 + uint64(b[5])<<40)
	case 7:
		return Uint64(0x01010101010101 + uint64(b[0]) + uint64(b[1])<<8 + uint64(b[2])<<16 + uint64(b[3])<<24 + uint64(b[4])<<32 + uint64(b[5])<<40 + uint64(b[6])<<48)
	default:
		return Uint64(0x0101010101010101 + uint64(b[0]) + uint64(b[1])<<8 + uint64(b[2])<<16 + uint64(b[3])<<24 + uint64(b[4])<<32 + uint64(b[5])<<40 + uint64(b[6])<<48 + uint64(b[7])<<56)
	}
}

// put the representation of a Uint64 into the provided []byte
func PutUint64(b []byte, x Uint64) int {
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

// return new Uint64 whos representation is as the source Uint64 but with added byte(s)
func (x Uint64) Append(b ...byte) Uint64 {
	buf := make([]byte, 8, 8)
	n := PutUint64(buf, x)
	copy(buf[n:], b)
	return GetUint64(buf[:n+len(b)]...)
}

// return new Uint64 whos representation is as the source Uint64 but with removed byte(s)
func (x Uint64) Truncate(c int) Uint64 {
	buf := make([]byte, 8, 8)
	n := PutUint64(buf, x)
	return GetUint64(buf[:n-c]...)
}

