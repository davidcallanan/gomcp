package javaio

// TODO: figure out prot version

import "bufio"

type PlayerPositionAndLook struct {
	X float64
	Y float64
	Z float64
	Yaw float32
	Pitch float32
	IsRelX bool
	IsRelY bool
	IsRelZ bool
	IsRelYaw bool
	IsRelPitch bool
}

func PacketId_PlayerPositionAndLook(protocol uint) int {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x36
	} else {
		// 1.14
		return 0x35
	}
	// older versions not supported
	return -1
}

func WritePlayerPositionAndLook(data PlayerPositionAndLook, stream *bufio.Writer) {
	WriteDouble(data.X, stream)
	WriteDouble(data.Y, stream)
	WriteDouble(data.Z, stream)
	WriteFloat(data.Yaw, stream)
	WriteFloat(data.Pitch, stream)

	var flags byte

	if data.IsRelX {
		flags |= 0x01
	}
	if data.IsRelY {
		flags |= 0x02
	}
	if data.IsRelZ {
		flags |= 0x04
	}
	if data.IsRelYaw {
		flags |= 0x10
	}
	if data.IsRelPitch {
		flags |= 0x08
	}

	WriteUByte(flags, stream)

	// Seems pointless for now.
	// Probably useful for interpolation, etc.
	var teleportId int32 = 0
	WriteVarInt(teleportId, stream)
}
