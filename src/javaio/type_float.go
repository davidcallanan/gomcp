package javaio

import "math"
import "bufio"

func ReadFloat(stream *bufio.Reader) (result float32, err error) {
	const len = 4
	var bytes [len]byte
	read, _ := stream.Read(bytes[:])

	if read < len {
		err = MalformedPacketError { "Float ended abruptly" }
		return
	}

	result = math.Float32frombits(
		uint32(bytes[0]) << 24 |
		uint32(bytes[1]) << 16 |
		uint32(bytes[2]) <<  8 |
		uint32(bytes[3]),
	)
	return
}

func WriteFloat(value float32, stream *bufio.Writer) {
	n := math.Float32bits(value)
	stream.Write([]byte {
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})
}
