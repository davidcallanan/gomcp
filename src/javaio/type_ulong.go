package javaio

import "bufio"

func WriteULong(value uint64, stream *bufio.Writer) {
	stream.Write([]byte {
		byte(value >> 56),
		byte(value >> 48),
		byte(value >> 40),
		byte(value >> 32),
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
}
