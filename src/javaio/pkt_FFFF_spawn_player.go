package javaio

import "bufio"
import "github.com/google/uuid"

type Packet_SpawnPlayer struct {
	EntityId int32
	Uuid uuid.UUID
	X float64
	Y float64
	Z float64
	Yaw uint8
	Pitch uint8
}

func PacketId_SpawnPlayer(protocol uint) int32 {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x05
	} else {
		// 1.14
		return 0x04
	}
	// todo: older versions not supported
	return -1
}

func Write_SpawnPlayer(data Packet_SpawnPlayer, ctx ClientContext, stream *bufio.Writer) {
	WriteVarInt(data.EntityId, stream)
	WriteUuidBin(data.Uuid, stream)
	WriteDouble(data.X, stream)
	WriteDouble(data.Y, stream)
	WriteDouble(data.Z, stream)
	WriteUByte(data.Yaw, stream)
	WriteUByte(data.Pitch, stream)

	if ctx.Protocol < 0x0286 {
		// 1.14 approximation
		// entity metadata must be sent in this version
		WriteUByte(0xff, stream) // end of entity metadata; no metadata sent for now
	}
}
