package javaio

import "bufio"

func WriteBool(value bool, stream *bufio.Writer) {
	if (value == true) {
		WriteUByte(0x01, stream)
	} else {
		WriteUByte(0x00, stream)
	}
}
