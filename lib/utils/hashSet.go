package utils

import "encoding/binary"

type Hashable interface {
	Hash() uint64
}

func Uint64ToBytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b
}
