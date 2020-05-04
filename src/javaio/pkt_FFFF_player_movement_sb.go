package javaio

import "bufio"

type Packet_PlayerPosSb struct {
	X float64
	Y float64
	Z float64
	OnGround bool
}

type Packet_PlayerLookSb struct {
	Yaw float32
	Pitch float32
	OnGround bool
}

type Packet_PlayerPosAndLookSb struct {
	X float64
	Y float64
	Z float64
	Yaw float32
	Pitch float32
	OnGround bool
}

// TODO: player actions: onground, sprint, sneak, etc.

func PacketId_PlayerPosSb(protocol uint) int {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x11
	} else {
		// 1.14
		return 0x11
	}
	// older versions not supported
	return -1
}

func PacketId_PlayerLookSb(protocol uint) int {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x13
	} else {
		// 1.14
		return 0x13
	}
	// older versions not supported
	return -1
}

func PacketId_PlayerPosAndLookSb(protocol uint) int {
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

func Read_PlayerPosSb(stream *bufio.Reader) (result Packet_PlayerPosSb, err error) {
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

	onGround, err := ReadBool(stream)
	if err != nil {
		return
	}
	
	result = Packet_PlayerPosSb {
		X: x,
		Y: y,
		Z: z,
		OnGround: onGround,
	}
	return
}

func Read_PlayerLookSb(stream *bufio.Reader) (result Packet_PlayerLookSb, err error) {
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
	
	result = Packet_PlayerLookSb {
		Yaw: yaw,
		Pitch: pitch,
		OnGround: onGround,
	}
	return
}

func Read_PlayerPosAndLookSb(stream *bufio.Reader) (result Packet_PlayerPosAndLookSb, err error) {
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
	
	result = Packet_PlayerPosAndLookSb {
		X: x,
		Y: y,
		Z: z,
		Yaw: yaw,
		Pitch: pitch,
		OnGround: onGround,
	}
	return
}
