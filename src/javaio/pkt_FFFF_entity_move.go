package javaio

import "bufio"

type Packet_EntityTranslate struct {
	EntityId int32
	DeltaX int16
	DeltaY int16
	DeltaZ int16
	Yaw uint8
	Pitch uint8
	OnGround bool
}

func PacketId_EntityTranslate(protocol uint) int {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x2A
	} else {
		// 1.14
		return 0x29
	}
	// todo: older versions
	return -1
}

func Write_EntityTranslate(data Packet_EntityTranslate, stream *bufio.Writer) {
	WriteVarInt(data.EntityId, stream)
	WriteShort(data.DeltaX, stream)
	WriteShort(data.DeltaY, stream)
	WriteShort(data.DeltaZ, stream)
	WriteUByte(data.Yaw, stream)
	WriteUByte(data.Pitch, stream)
	WriteBool(data.OnGround, stream)
}

type Packet_EntityVelocity struct {
	EntityId int32
	X int16
	Y int16
	Z int16
}

func PacketId_EntityVelocity(protocol uint) int {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x46
	} else {
		// 1.14
		return 0x45
	}
	// todo: older versions
	return -1
}

func Write_EntityVelocity(data Packet_EntityVelocity, stream *bufio.Writer) {
	WriteVarInt(data.EntityId, stream)
	WriteShort(data.X, stream)
	WriteShort(data.Y, stream)
	WriteShort(data.Z, stream)
}
