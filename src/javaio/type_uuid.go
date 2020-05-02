package javaio

import "bufio"
import "github.com/google/uuid"

func WriteUuidBin(uuid uuid.UUID, stream *bufio.Writer) {
	data, err := uuid.MarshalBinary()

	if err != nil {
		panic(err)
	}
	
	for i := range data {
		stream.WriteByte(data[len(data) - 1 - i])	
	}
}
