package varbinary

import "testing"
import "fmt"


func TestUint64(t *testing.T) {
	var tests = []struct {
	  v Uint64
	  length int
	}{
	  {0, 0},
	  {255, 1},
	  {256, 1},
	  {257, 2},
	  {123456, 3},
	  {123456789, 4},
	  {123456789012, 5},
	  {123456789012345678, 8},
	}

	b:=make([]byte,8)
	for _, tt := range tests {
		l:=Uint64Put(tt.v,b)
		if l != tt.length {
			t.Errorf("%v length %v",tt,l)
		}
		var v Uint64
		(*Uint64).UnmarshalBinary(&v,b[:l])
		if v != tt.v {
			t.Errorf("%v encoded/decoded to %v",tt,l)
		}
	}		
}

func TestUint64Encode(t *testing.T) {
	if fmt.Sprintf("%s",Uint64(257))!="00 00"{t.Errorf("not '00 00' (%s)",Uint64(257))}
	n:=65793   // 1+1<<8+1<<16
	if fmt.Sprint(Uint64(n))!="00 00 00"{t.Errorf("not '00 00 00' (%s)",Uint64(n))}
	n=16843009
	if fmt.Sprint(Uint64(n))!="00 00 00 00"{t.Errorf("not '00 00 00 00' (%s)",Uint64(n))}
	n-=1
	if fmt.Sprint(Uint64(n))!="FF FF FF"{t.Errorf("not 'FF FF FF' (%s)",Uint64(n))}
	n=1103823438081
	if fmt.Sprint(Uint64(n))!="00 00 00 00 00 00"{t.Errorf("not '00 00 00 00 00 00' (%s)",Uint64(n))}
	n-=1
	if fmt.Sprint(Uint64(n))!="FF FF FF FF FF"{t.Errorf("not 'FF FF FF FF FF' (%s)",Uint64(n))}
	var d uint64 =0xffffffffffffffff
	if fmt.Sprint(Uint64(d))!="FE FE FE FE FE FE FE FE"{t.Errorf("not 'FE FE FE FE FE FE FE FE' (%s)",Uint64(d))}
}

func TestUint64Decode(t *testing.T) {
	if Uint64Decoder([]byte{}...)!=0 {t.Errorf("empty slice not 0 (%d)",Uint64Decoder([]byte{}...))}
	if Uint64Decoder([]byte{0x00}...)!=1 {t.Errorf("empty 0x00 not 1 (%d)",Uint64Decoder([]byte{0x00}...))}
	if Uint64Decoder([]byte{0x00,0x00}...)!=257 {t.Errorf("empty []byte{0x00,0x00} not 257 (%d)",Uint64Decoder([]byte{0x00,0x00}...))}

	v:=Uint64Decoder([]byte{0xfe,0xfe,0xfe,0xfe,0xfe,0xfe,0xfe,0xfe}...)
	if fmt.Sprintf("%d",v)!="18446744073709551615"{t.Errorf("not '18446744073709551615' (%d)",v)}
	v++
	if fmt.Sprintf("%d",v)!="0"{t.Errorf("not '0' (%d)",v)}
}

