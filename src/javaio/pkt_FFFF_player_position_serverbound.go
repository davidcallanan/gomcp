package javaio

import "bufio"

type Packet_PlayerPositionAndLookServerbound struct {
	X float64
	Y float64
	Z float64
	Yaw float32
	Pitch float32
	OnGround bool
}

func PacketId_PlayerPositionAndLookServerbound(protocol uint) int {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x12
	} else {
		// 1.14
		return 0x12
	}
	// older versions not supported
	return -1
}

func Read_PlayerPositionAndLookServerbound(stream *bufio.Reader) (result Packet_PlayerPositionAndLookServerbound, err error) {
	x, err := ReadDouble(stream)
	if err != nil {
		return
	}
	
	y, err := ReadDouble(stream)
	if err != nil {
		return
	}

	z, err := ReadDouble(stream)
	if err != nil {
		return
	}

	yaw, err := ReadFloat(stream)
	if err != nil {
		return
	}

	pitch, err := ReadFloat(stream)
	if err != nil {
		return
	}
	
	onGround, err := ReadBool(stream)
	if err != nil {
		return
	}
	
	result = Packet_PlayerPositionAndLookServerbound {
		X: x,
		Y: y,
		Z: z,
		Yaw: yaw,
		Pitch: pitch,
		OnGround: onGround,
	}
	return
}
