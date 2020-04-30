package javaio

import "bufio"

// TODO: figure out prot version

type CompassPosition struct {
	Location BlockPosition
}

func WriteCompassPosition(compassPosition CompassPosition, stream *bufio.Writer) {
	WriteBlockPos(compassPosition.Location, stream)
}

func PacketId_CompassPosition(protocol uint) int {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x4e
	} else {
		// 1.14
		return 0x4d
	}
	// older versions not supported
	return -1
}
