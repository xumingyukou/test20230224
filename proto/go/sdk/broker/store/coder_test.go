package store

import (
	"bytes"
	"math"
	"testing"
)

func TestPack(t *testing.T) {
	var enc IndexEncoder
	var dec IndexDecoder

	//基本功能测试
	enc.PackUint16(65535).PackUint16(1).PackUint16(0)

	{
		dec.FromBytes(enc.Bytes(false))
		var a, b, c uint16
		dec.UnpackUint16(&a).UnpackUint16(&b).UnpackUint16(&c)

		if a != 65535 || b != 1 || c != 0 {
			t.Errorf("Unpack error %d %d %d", a, b, c)
		}
	}

	enc.Reset()
	enc.PackRaw([]byte("t_")).PackInt32(int32(2147483647)).PackInt32(-2147483648).PackInt32(0)

	{
		dec.FromBytes(enc.Bytes(false))
		hdr := make([]byte, 2, 10)
		var a, b, c int32
		dec.UnpackRaw(&hdr).UnpackInt32(&a).UnpackInt32(&b).UnpackInt32(&c)
		if string(hdr[:]) != "t_" {
			t.Error("Unpack hdr error " + string(hdr[:]))
		}
		if a != 2147483647 || b != -2147483648 || c != 0 {
			t.Errorf("Unpack error %d %d %d", a, b, c)
		}
	}

	enc.Reset()

	enc.PackRaw([]byte("i__")).PackInt16(-32768).PackInt32(-2147483645).PackInt64(math.MinInt64).PackString("abcdefgh").
		PackUint16(65535).PackUint32(2147483642).PackUint64(0xffffffffffffffff)

	{
		hdr := make([]byte, 3, 10)
		var a int16
		var b int32
		var c int64
		var d string
		var e uint16
		var f uint32
		var g uint64

		dec.FromBytes(enc.Bytes(false)).UnpackRaw(&hdr).UnpackInt16(&a).UnpackInt32(&b).UnpackInt64(&c).
			UnpackString(&d).UnpackUint16(&e).UnpackUint32(&f).UnpackUint64(&g)

		if string(hdr[:]) != "i__" {
			t.Error("Unpack hdr error " + string(hdr[:]))
		}

		if a != -32768 || b != -2147483645 || c != math.MinInt64 || d != "abcdefgh" {
			t.Errorf("Unpack error %d %d %d %s", a, b, c, d)
		}

		if e != 65535 || f != 2147483642 || g != 0xffffffffffffffff {
			t.Errorf("Unpack error %d %d %d", e, f, g)
		}
	}

	// 比较编码后的值是否符合编码前的顺序
	enc.Reset()
	//4bytes Raw + int32 + 2bytes raw + string + int64
	v1 := enc.PackRaw([]byte("prx1")).PackInt32(0).PackRaw([]byte("i_")).PackString("abc").PackInt64(0x7fffffffffffffff).Bytes(true)
	enc.Reset()
	v2 := enc.PackRaw([]byte("prx1")).PackInt32(0).PackRaw([]byte("i_")).PackString("abcdefg").PackInt64(0x0).Bytes(true)

	if bytes.Compare(v1, v2) >= 0 {
		t.Error("v1 should less than v2")
	}

	enc.Reset()
	v2 = enc.PackRaw([]byte("prx1")).PackInt32(0).PackRaw([]byte("i_")).PackString("abc").PackInt64(0x0).Bytes(true)

	if bytes.Compare(v1, v2) <= 0 {
		t.Error("v1 should greater than v2")
	}

	enc.Reset()
	v2 = enc.PackRaw([]byte("prx1")).PackInt32(0).PackRaw([]byte("j_")).PackString("abc").PackInt64(0x7fffffffffffffff).Bytes(true)

	if bytes.Compare(v1, v2) >= 0 {
		t.Error("v1 should less than v2")
	}

	enc.Reset()

	v2 = enc.PackRaw([]byte("prx1")).PackInt32(0).PackRaw([]byte("i_")).PackString("abcdefghijk").PackInt64(0).Bytes(true)

	if bytes.Compare(v1, v2) >= 0 {
		t.Error("v1 should less than v2")
	}
	enc.Reset()

	v2 = enc.PackRaw([]byte("prx1")).PackInt32(-1).PackRaw([]byte("i_")).PackString("abcdefg").PackInt64(0x0).Bytes(true)
	if bytes.Compare(v1, v2) <= 0 {
		t.Error("v1 should greater than v2")
	}

	enc.Reset()
	v2 = enc.PackRaw([]byte("prx2")).PackInt32(0).PackRaw([]byte("i_")).PackString("abc").PackInt64(0x7fffffffffffffff).Bytes(true)
	if bytes.Compare(v1, v2) >= 0 {
		t.Error("v1 should less than v2")
	}

	enc.Reset()
	enc.PackRaw([]byte("prx2")).PackInt32(0).PackRaw([]byte("abc"))

	{
		dec.FromBytes(enc.Bytes(false))
		a := make([]byte, 4)
		var b int32
		c := make([]byte, 10)

		dec.UnpackRaw(&a).UnpackInt32(&b).UnpackRaw(&c)

		if string(a) != "prx2" || b != 0 || string(c) != "abc" {
			t.Errorf("Unpack error %s %d %s", string(a), b, string(c))
		}
	}
}
