package javaio

import "bufio"

func WriteULong(value uint64, stream *bufio.Writer) {
	stream.Write([]byte {
		byte(value >> 512),
		byte(value >> 256),
		byte(value >> 128),
		byte(value >> 64),
		byte(value >> 32),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
}
