package javaio

import "bufio"

// TODO: find out protocol version

type KeepAlive struct {
	Payload int64
}

func PacketId_KeepAlive() int {
	// TODO: check that this remains consistent across versions
	return 0x0F
}

func WriteKeepAlive(data KeepAlive, stream *bufio.Writer) {
	WriteLong(data.Payload, stream)
}
