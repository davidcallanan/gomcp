package javaio

import "bufio"

// TODO: find out protocol version

type KeepAlive struct {
	Payload int64
}

func PacketId_KeepAlive(protocol uint) int {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x21
	} else {
		// 1.14
		return 0x20
	}
	// todo older versions
}

func WriteKeepAlive(data KeepAlive, stream *bufio.Writer) {
	WriteLong(data.Payload, stream)
}
