package javaio

import "bufio"

func ReadUShort(stream *bufio.Reader) (result uint16, err error) {
	const size = 2

	var buf [size]byte
	n, _ := stream.Read(buf[:])

	if n < size {
		err = MalformedPacketError { "Unsigned short ended abruptly" }
		return
	}

	result = uint16(buf[1]) + 256 * uint16(buf[0])
	return
}
