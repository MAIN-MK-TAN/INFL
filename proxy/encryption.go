// This encryption module is NOT complete.
package main

import (
	"io"
)

var z85Encoder = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.-:+=^!/*?&<>()[]{}@%$#")
var z85Decoder [256]byte

func init() {
	for i := range z85Decoder {
		z85Decoder[i] = 0xFF
	}
	for i, b := range z85Encoder {
		z85Decoder[b] = byte(i)
	}
}

func Enc(data []byte) []byte {
	pad := (4 - len(data)%4) % 4
	data = append(data, make([]byte, pad)...)
	return Z85Encode(data)
}

func dEnc(data []byte) ([]byte, error) {
	return Z85Decode(data)
}

func Z85Encode(src []byte) []byte {
	if len(src)%4 != 0 {
		return nil
	}
	out := make([]byte, len(src)*5/4)
	var val uint32
	for i, j := 0, 0; i < len(src); i += 4 {
		val = uint32(src[i])<<24 | uint32(src[i+1])<<16 | uint32(src[i+2])<<8 | uint32(src[i+3])
		for k := 4; k >= 0; k-- {
			out[j+k] = z85Encoder[val%85]
			val /= 85
		}
		j += 5
	}
	return out
}

func Z85Decode(src []byte) ([]byte, error) {
	if len(src)%5 != 0 {
		return nil, io.ErrUnexpectedEOF
	}
	out := make([]byte, len(src)*4/5)
	var val uint32
	for i, j := 0, 0; i < len(src); i += 5 {
		val = 0
		for k := 0; k < 5; k++ {
			d := z85Decoder[src[i+k]]
			if d == 0xFF {
				return nil, io.ErrUnexpectedEOF
			}
			val = val*85 + uint32(d)
		}
		out[j+0] = byte(val >> 24)
		out[j+1] = byte(val >> 16)
		out[j+2] = byte(val >> 8)
		out[j+3] = byte(val)
		j += 4
	}
	return out, nil
}
