package varbinary
// variable length binary encoding/decoding integer values.
// differs from other techniques in it uses all permutations of bytes. (up to the redundant states.)
// an encodings length, (in bytes not bits), carries information, it is not carried in the binary data itself, so to decode the length is required to be known.
// some binary states are never produced as an encoding, here they are redundant and don’t have a decoding.
// (redundant states occur due to the extra information carried in the variable length.)

import "errors"
import "fmt"
import "io"

// Uint64 unables variable length binary encoding/decoding of uint64 values.
// implements: io.MarshalBinary,io.UnmarshalBinary,io.Reader (provided buffer needs to be long enough for encoding)
type Uint64 uint64

var decErr error = errors.New("Bytes do not represent an encoding.")
var bufErr error = errors.New("Encoding can not be put in supplied Buffer.")

// string representation is the hexadecimal of the encoding, (implementing Stringer)
func (x Uint64) String() string {
	b := make([]byte, 8, 8)
	n, _ := x.Read(b)
	return fmt.Sprintf("% X", b[:n])
}

// encode a Uint64's into a provided byte[]. (implementing io.Reader)
// error returned if buffer size not big enough to contain the encoding.
func (u *Uint64) Read(b []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = bufErr
		}
	}()
	
	return Uint64Put(*u,b), io.EOF
}

// decode a Uint64 from the provided []byte.  (implementing io.Writer.)
// error returned if value is unused, redundant, encoding
func (u *Uint64) Write(b []byte) (int, error) {
	*u = Uint64Decoder(b...)
	if len(b) > 7 {
		if *u>Uint64(0x0101010101010101){return 0,decErr}
		return 8, io.EOF
	}
	return len(b), nil
}


func (u *Uint64) MarshalBinary() (data []byte, err error){
	return Uint64Encoder(*u),nil
}

// return the binary encoding of a Uint64
func Uint64Encoder(u Uint64) []byte{
	if u == 0 {
		return []byte{}
	}
	if u < 0x0101 {
		u -= 1
		return []byte{uint8(u)}
	}
	if u < 0x010101 {
		u -= 0x0101
		return []byte{uint8(u),uint8(u >> 8)}
	}
	if u < 0x01010101 {
		u -= 0x010101
		return []byte{uint8(u),uint8(u >> 8),uint8(u >> 16)}
	}
	if u < 0x0101010101 {
		u -= 0x01010101
		return []byte{uint8(u),uint8(u >> 8),uint8(u >> 16),uint8(u >> 24)}
	}
	if u < 0x010101010101 {
		u -= 0x0101010101
		return []byte{uint8(u),uint8(u >> 8),uint8(u >> 16),uint8(u >> 24),uint8(u >> 32)}
	}
	if u < 0x01010101010101 {
		u -= 0x010101010101
		return []byte{uint8(u),uint8(u >> 8),uint8(u >> 16),uint8(u >> 24),uint8(u >> 32), uint8(u >> 40)}
	}
	if u < 0x0101010101010101 {
		u -= 0x01010101010101
		return []byte{uint8(u),uint8(u >> 8),uint8(u >> 16),uint8(u >> 24),uint8(u >> 32), uint8(u >> 40),uint8(u >> 48)}
	}
	u -= 0x0101010101010101
	return []byte{uint8(u),uint8(u >> 8),uint8(u >> 16),uint8(u >> 24),uint8(u >> 32), uint8(u >> 40),uint8(u >> 48),uint8(u >> 56)}
}

func (u *Uint64) UnmarshalBinary(data []byte) error {
	switch len(data) {
	case 0:
		*u= Uint64(0)
	case 1:
		*u= Uint64(0x1 + uint64(data[0]))
	case 2:
		*u= Uint64(0x0101 + uint64(data[0]) + uint64(data[1])<<8)
	case 3:
		*u= Uint64(0x010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16)
	case 4:
		*u= Uint64(0x01010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24)
	case 5:
		*u= Uint64(0x0101010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24 + uint64(data[4])<<32)
	case 6:
		*u= Uint64(0x010101010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24 + uint64(data[4])<<32 + uint64(data[5])<<40)
	case 7:
		*u= Uint64(0x01010101010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24 + uint64(data[4])<<32 + uint64(data[5])<<40 + uint64(data[6])<<48)
	case 8:
		*u= Uint64(0x0101010101010101 + uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24 + uint64(data[4])<<32 + uint64(data[5])<<40 + uint64(data[6])<<48 + uint64(data[7])<<56)
		if *u<Uint64(0x0101010101010101){return decErr}
	default:
		return decErr
	}
	return nil
}


// return the Uint64 represented by some bytes
func Uint64Decoder(b ...byte) Uint64 {
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
func Uint64Put(x Uint64,b []byte) int {
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

