package javaio

import "bufio"

func WriteUByte(value byte, stream *bufio.Writer) {
	stream.WriteByte(value)
}
