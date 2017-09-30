package varbinary

import "testing"
import "fmt"


func TestUint64(t *testing.T) {
	b:=make([]byte,8)
	l:=PutUint64(b,Uint64(0))
	if l!=0	{t.Errorf("0 Not len 0")}
	v:=GetUint64(b[:l]...)
	if v!=0 {t.Errorf("not 0 (%d)",v)}

	l=PutUint64(b,Uint64(255))
	if l!=1	{t.Errorf("255 Not len 1")}
	v=GetUint64(b[:l]...)
	if v!=255 {t.Errorf("not 255 (%d)",v)}

	l=PutUint64(b,Uint64(256))
	if l!=1	{t.Errorf("256 Not len 2")}
	v=GetUint64(b[:l]...)
	if v!=256 {t.Errorf("not 256 (%d)",v)}

	l=PutUint64(b,Uint64(257))
	if l!=2	{t.Errorf("257 Not len 2")}
	v=GetUint64(b[:l]...)
	if v!=257 {t.Errorf("not 257 (%d)",v)}

	l=PutUint64(b,Uint64(257))
	if l!=2	{t.Errorf("257 Not len 2")}
	v=GetUint64(b[:l]...)
	if v!=257 {t.Errorf("not 257 (%d)",v)}

	l=PutUint64(b,Uint64(123456))
	if l!=3	{t.Errorf("123456 Not len 3")}
	v=GetUint64(b[:l]...)
	if v!=123456 {t.Errorf("not 123456 (%d)",v)}

	l=PutUint64(b,Uint64(123456789))
	if l!=4	{t.Errorf("123456789 Not len 4")}
	v=GetUint64(b[:l]...)
	if v!=123456789 {t.Errorf("not 123456789 (%d)",v)}

	l=PutUint64(b,Uint64(123456789012))
	if l!=5	{t.Errorf("123456789012 Not len 5")}
	v=GetUint64(b[:l]...)
	if v!=123456789012 {t.Errorf("not 123456789012 (%d)",v)}

	l=PutUint64(b,Uint64(123456789012345678))
	if l!=8	{t.Errorf("123456789012345678 Not len 8")}
	v=GetUint64(b[:l]...)
	if v!=123456789012345678 {t.Errorf("not 123456789012345678 (%d)",v)}
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
	if GetUint64()!=0 {t.Errorf("empty slice not 0 (%d)",GetUint64())}
	if GetUint64(0x00)!=1 {t.Errorf("empty 0x00 not 1 (%d)",GetUint64(0x00))}
	if GetUint64(0x00,0x00)!=257 {t.Errorf("empty []byte{0x00,0x00} not 257 (%d)",GetUint64(0x00,0x00))}

	v:=GetUint64(0xfe,0xfe,0xfe,0xfe,0xfe,0xfe,0xfe,0xfe)
	if fmt.Sprintf("%d",v)!="18446744073709551615"{t.Errorf("not '18446744073709551615' (%d)",v)}
	v++
	if fmt.Sprintf("%d",v)!="0"{t.Errorf("not '0' (%d)",v)}
}

func TestUint64AppendTruncate(t *testing.T) {
	v:=Uint64.Append(Uint64(0),0x00)
	if fmt.Sprintf("%d",v)!="1"{t.Errorf("not '1' (%d)",v)}
	v=Uint64.Truncate(v,1)
	if fmt.Sprintf("%d",v)!="0"{t.Errorf("not '0' (%d)",v)}
}

