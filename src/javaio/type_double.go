package javaio

import "math"
import "bufio"

func ReadDouble(stream *bufio.Reader) (result float64, err error) {
	const len = 8
	var bytes [len]byte
	read, _ := stream.Read(bytes[:])

	if read < len {
		err = MalformedPacketError { "Double ended abruptly" }
		return
	}

	result = math.Float64frombits(
		uint64(bytes[0]) << 56 |
		uint64(bytes[1]) << 48 |
		uint64(bytes[2]) << 40 |
		uint64(bytes[3]) << 32 |
		uint64(bytes[4]) << 24 |
		uint64(bytes[5]) << 16 |
		uint64(bytes[6]) <<  8 |
		uint64(bytes[7]),
	)
	return
}

func WriteDouble(value float64, stream *bufio.Writer) {
	n := math.Float64bits(value)
	stream.Write([]byte {
		byte(n >> 56),
		byte(n >> 48),
		byte(n >> 40),
		byte(n >> 32),
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})
}
