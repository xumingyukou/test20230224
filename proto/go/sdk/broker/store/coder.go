package store

import (
	"bytes"
	"encoding/binary"
	"runtime"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

/*
	用于创建可内存中按照字节序可比较key的辅助函数
	目前支持普通bytes，整型和字符串类型的编解码
	使用方法：
	整形直接使用PackInt/Uint系列接口
	固定长度字符串直接使用PackRaw接口即可
	变长字符串如果在最末尾可以直接使用PackRaw接口，如果在中间需要使用PackString接口
*/

type IndexEncoder struct {
	buf bytes.Buffer
}

type IndexDecoder struct {
	buf []byte
}

func (e *IndexEncoder) Reset() *IndexEncoder {
	e.buf.Reset()
	return e
}

// 返回的值只在下次修改pack前有效
func (e *IndexEncoder) Bytes(dup bool) []byte {
	if dup {
		tmp := make([]byte, len(e.buf.Bytes()))
		copy(tmp, e.buf.Bytes())
		return tmp
	} else {
		return e.buf.Bytes()
	}
}

// 写入原始bytes，不进行pack处理
func (e *IndexEncoder) PackRaw(b []byte) *IndexEncoder {
	e.buf.Write(b)
	return e
}

func (e *IndexEncoder) PackInt16(i int16) *IndexEncoder {
	var tmp [2]byte
	u := uint16(i) + 0x8000
	binary.BigEndian.PutUint16(tmp[:], u)
	e.buf.Write(tmp[:])

	return e
}

func (e *IndexEncoder) PackUint16(i uint16) *IndexEncoder {
	var tmp [2]byte
	binary.BigEndian.PutUint16(tmp[:], i)
	e.buf.Write(tmp[:])

	return e
}

func (e *IndexEncoder) PackInt32(i int32) *IndexEncoder {
	var tmp [4]byte
	u := uint32(i) + 0x80000000
	binary.BigEndian.PutUint32(tmp[:], u)
	e.buf.Write(tmp[:])

	return e
}

func (e *IndexEncoder) PackUint32(i uint32) *IndexEncoder {
	var tmp [4]byte
	binary.BigEndian.PutUint32(tmp[:], i)
	e.buf.Write(tmp[:])

	return e
}

func (e *IndexEncoder) PackInt64(i int64) *IndexEncoder {
	var tmp [8]byte
	u := uint64(i) + 0x8000000000000000
	binary.BigEndian.PutUint64(tmp[:], u)
	e.buf.Write(tmp[:])

	return e
}

func (e *IndexEncoder) PackUint64(i uint64) *IndexEncoder {
	var tmp [8]byte
	binary.BigEndian.PutUint64(tmp[:], i)
	e.buf.Write(tmp[:])

	return e
}

func (e *IndexEncoder) PackString(s string) *IndexEncoder {
	dLen := len(s)
	reallocSize := (dLen/encGroupSize + 1) * (encGroupSize + 1)
	if reallocSize > e.buf.Cap() {
		e.buf.Grow(reallocSize - e.buf.Cap())
	}

	for idx := 0; idx <= dLen; idx += encGroupSize {
		remain := dLen - idx
		padCount := 0
		if remain >= encGroupSize {
			e.buf.WriteString(s[idx : idx+encGroupSize])
		} else {
			padCount = encGroupSize - remain
			e.buf.WriteString(s[idx:])
			e.buf.Write(pads[:padCount])
		}

		marker := encMarker - byte(padCount)
		e.buf.WriteByte(marker)
	}

	return e
}

func (d *IndexDecoder) Reset() *IndexDecoder {
	d.buf = nil
	return d
}

func (d *IndexDecoder) FromBytes(b []byte) *IndexDecoder {
	d.buf = b
	return d
}

func (d *IndexDecoder) RemainLen() int {
	return len(d.buf)
}

func (d *IndexDecoder) UnpackRaw(b *[]byte) *IndexDecoder {
	n := copy(*b, d.buf)
	*b = (*b)[:n]

	d.buf = d.buf[n:]

	return d
}

func (d *IndexDecoder) UnpackInt16(i *int16) *IndexDecoder {
	if len(d.buf) < 2 {
		panic("Unpack int16 failed")
	}

	u := binary.BigEndian.Uint16(d.buf)
	*i = int16(u - 0x8000)
	d.buf = d.buf[2:]

	return d
}

func (d *IndexDecoder) UnpackUint16(i *uint16) *IndexDecoder {
	if len(d.buf) < 2 {
		panic("Unpack uint16 failed")
	}

	*i = binary.BigEndian.Uint16(d.buf)
	d.buf = d.buf[2:]

	return d
}

func (d *IndexDecoder) UnpackInt32(i *int32) *IndexDecoder {
	if len(d.buf) < 4 {
		panic("Unpack int32 failed")
	}

	u := binary.BigEndian.Uint32(d.buf)
	*i = int32(u - 0x80000000)
	d.buf = d.buf[4:]

	return d
}

func (d *IndexDecoder) UnpackUint32(i *uint32) *IndexDecoder {
	if len(d.buf) < 4 {
		panic("Unpack uint32 failed")
	}

	*i = binary.BigEndian.Uint32(d.buf)

	d.buf = d.buf[4:]

	return d
}

func (d *IndexDecoder) UnpackInt64(i *int64) *IndexDecoder {
	if len(d.buf) < 8 {
		panic("Unpack int64 failed")
	}

	u := binary.BigEndian.Uint64(d.buf)
	*i = int64(u - 0x8000000000000000)

	d.buf = d.buf[8:]

	return d
}

func (d *IndexDecoder) UnpackUint64(i *uint64) *IndexDecoder {
	if len(d.buf) < 8 {
		panic("Unpack uint64 failed")
	}

	*i = binary.BigEndian.Uint64(d.buf)
	d.buf = d.buf[8:]

	return d
}

func (d *IndexDecoder) UnpackString(s *string) *IndexDecoder {
	var groupBytes []byte
	var sb strings.Builder

	for {
		if len(d.buf) < encGroupSize+1 {
			panic("Unpack string failed, insufficient bytes to decode value")
		}
		groupBytes = d.buf[:encGroupSize+1]

		group := groupBytes[:encGroupSize]
		marker := groupBytes[encGroupSize]

		padCount := encMarker - marker
		if padCount > encGroupSize {
			panic("Unpack string failed, invalid marker byte")
		}

		realGroupSize := encGroupSize - padCount
		sb.Write(group[:realGroupSize])

		d.buf = d.buf[encGroupSize+1:]

		if padCount != 0 {
			var padByte = encPad
			// Check validity of padding bytes.
			for _, v := range group[realGroupSize:] {
				if v != padByte {
					panic("Unpack string failed, invalid padding byte")
				}
			}
			break
		}
	}

	*s = sb.String()

	return d
}

// string encoding，from tidb
const (
	encGroupSize = 8
	encMarker    = byte(0xFF)
	encPad       = byte(0x0)
)

var (
	pads = make([]byte, encGroupSize)
)

// EncodeBytes guarantees the encoded value is in ascending order for comparison,
// encoding with the following rule:
//  [group1][marker1]...[groupN][markerN]
//  group is 8 bytes slice which is padding with 0.
//  marker is `0xFF - padding 0 count`
// For example:
//   [] -> [0, 0, 0, 0, 0, 0, 0, 0, 247]
//   [1, 2, 3] -> [1, 2, 3, 0, 0, 0, 0, 0, 250]
//   [1, 2, 3, 0] -> [1, 2, 3, 0, 0, 0, 0, 0, 251]
//   [1, 2, 3, 4, 5, 6, 7, 8] -> [1, 2, 3, 4, 5, 6, 7, 8, 255, 0, 0, 0, 0, 0, 0, 0, 0, 247]
// Refer: https://github.com/facebook/mysql-5.6/wiki/MyRocks-record-format#memcomparable-format
func EncodeBytes(b []byte, data []byte) []byte {
	// Allocate more space to avoid unnecessary slice growing.
	// Assume that the byte slice size is about `(len(data) / encGroupSize + 1) * (encGroupSize + 1)` bytes,
	// that is `(len(data) / 8 + 1) * 9` in our implement.
	dLen := len(data)
	reallocSize := (dLen/encGroupSize + 1) * (encGroupSize + 1)
	result := reallocBytes(b, reallocSize)
	for idx := 0; idx <= dLen; idx += encGroupSize {
		remain := dLen - idx
		padCount := 0
		if remain >= encGroupSize {
			result = append(result, data[idx:idx+encGroupSize]...)
		} else {
			padCount = encGroupSize - remain
			result = append(result, data[idx:]...)
			result = append(result, pads[:padCount]...)
		}

		marker := encMarker - byte(padCount)
		result = append(result, marker)
	}

	return result
}

// EncodedBytesLength returns the length of data after encoded
func EncodedBytesLength(dataLen int) int {
	mod := dataLen % encGroupSize
	padCount := encGroupSize - mod
	return dataLen + padCount + 1 + dataLen/encGroupSize
}

func decodeBytes(b []byte, buf []byte, reverse bool) ([]byte, []byte, error) {
	if buf == nil {
		buf = make([]byte, 0, len(b))
	}
	buf = buf[:0]
	for {
		if len(b) < encGroupSize+1 {
			return nil, nil, errors.New("insufficient bytes to decode value")
		}

		groupBytes := b[:encGroupSize+1]

		group := groupBytes[:encGroupSize]
		marker := groupBytes[encGroupSize]

		var padCount byte
		if reverse {
			padCount = marker
		} else {
			padCount = encMarker - marker
		}
		if padCount > encGroupSize {
			return nil, nil, errors.Errorf("invalid marker byte, group bytes %q", groupBytes)
		}

		realGroupSize := encGroupSize - padCount
		buf = append(buf, group[:realGroupSize]...)
		b = b[encGroupSize+1:]

		if padCount != 0 {
			var padByte = encPad
			if reverse {
				padByte = encMarker
			}
			// Check validity of padding bytes.
			for _, v := range group[realGroupSize:] {
				if v != padByte {
					return nil, nil, errors.Errorf("invalid padding byte, group bytes %q", groupBytes)
				}
			}
			break
		}
	}
	if reverse {
		reverseBytes(buf)
	}
	return b, buf, nil
}

// DecodeBytes decodes bytes which is encoded by EncodeBytes before,
// returns the leftover bytes and decoded value if no error.
// `buf` is used to buffer data to avoid the cost of makeslice in decodeBytes when DecodeBytes is called by Decoder.DecodeOne.
func DecodeBytes(b []byte, buf []byte) ([]byte, []byte, error) {
	return decodeBytes(b, buf, false)
}

func reallocBytes(b []byte, n int) []byte {
	newSize := len(b) + n
	if cap(b) < newSize {
		bs := make([]byte, len(b), newSize)
		copy(bs, b)
		return bs
	}

	return b
}

// See https://golang.org/src/crypto/cipher/xor.go
const wordSize = int(unsafe.Sizeof(uintptr(0)))
const supportsUnaligned = runtime.GOARCH == "386" || runtime.GOARCH == "amd64"

func fastReverseBytes(b []byte) {
	n := len(b)
	w := n / wordSize
	if w > 0 {
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			bw[i] = ^bw[i]
		}
	}

	for i := w * wordSize; i < n; i++ {
		b[i] = ^b[i]
	}
}

func safeReverseBytes(b []byte) {
	for i := range b {
		b[i] = ^b[i]
	}
}

func reverseBytes(b []byte) {
	if supportsUnaligned {
		fastReverseBytes(b)
		return
	}

	safeReverseBytes(b)
}
