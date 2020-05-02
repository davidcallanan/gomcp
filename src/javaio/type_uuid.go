package javaio

import "bufio"
import "github.com/google/uuid"

func WriteUuidBin(uuid uuid.UUID, stream *bufio.Writer) {
	data, err := uuid.MarshalBinary()

	if err != nil {
		panic(err)
	}
	
	// TODO: not sure which order to use here
	stream.Write(data)

	// for i := range data {
	// 	stream.WriteByte(data[len(data) - 1 - i])	
	// }
}
