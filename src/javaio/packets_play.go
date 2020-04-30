package javaio

// Clientbound

type KeepAlive struct {
	Payload int64
}

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

type ChunkData struct {
	X int32
	Z int32
	IsNew bool
	Sections [][]uint32
}

// Serverbound
