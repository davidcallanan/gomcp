package javaio

import "bufio"

func WriteShort(value int16, stream *bufio.Writer) {
	stream.Write([]byte {
		byte(value >> 8),
		byte(value),
	})
}
