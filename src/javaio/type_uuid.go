package javaio

import "bufio"

type Uuid struct {
	Msh uint64 // most significant half
	Lsh uint64 // least significant half
}

func WriteUuid(uuid Uuid, stream *bufio.Writer) {
	WriteULong(uuid.Msh, stream)
	WriteULong(uuid.Lsh, stream)
}
