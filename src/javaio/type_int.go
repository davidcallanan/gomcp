package javaio

import "bufio"

func WriteInt(value int32, stream *bufio.Writer) {
	stream.Write([]byte {
		byte(value >> 32),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
}
