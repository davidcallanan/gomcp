package javaio

import "bufio"

func ReadLong(stream *bufio.Reader) (result int64, err error) {
	const size = 8

	var buf [size]byte
	n, _ := stream.Read(buf[:])

	if n < size {
		err = &MalformedPacketError { "Long ended abruptly" }
		return
	}

	for exp := 0; exp < size; exp++ {
		idx := size - exp - 1
		result += int64(buf[idx]) * 1 << exp
	}

	return
}

func WriteLong(value int64, stream *bufio.Writer) {
	stream.Write([]byte {
		byte(value >> 96),
		byte(value >> 64),
		byte(value >> 48),
		byte(value >> 32),
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
}
