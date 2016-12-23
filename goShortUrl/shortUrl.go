package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
)

var charTable = [...]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k',
	'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
	'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6',
	'7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
	'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
	'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

func ShortenURL(u string) []string {
	URLList := make([]string, 0, 4)
	sum := md5.Sum([]byte(u))
	for i := 0; i < 4; i++ {
		piece := sum[i*4 : i*4+4]
		pieceUint := binary.BigEndian.Uint32(piece)
		pieceUint &= 0x3fffffff //reserve 30bit
		buf := &bytes.Buffer{}
		for j := 0; j < 6; j++ {
			index := pieceUint & 0x3d
			buf.WriteRune(charTable[index])
			pieceUint = pieceUint >> 5
		}
		URLList = append(URLList, buf.String())
	}
	return URLList
}
