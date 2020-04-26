package javaio

// Clientbound

type KeepAlive struct {
	Payload int64
}

type JoinGame struct {
	Eid int32
	Gamemode Gamemode
	Hardcore bool
	Dimension Dimension
	ViewDistance int32
	ReducedDebugInfo bool
	EnableRespawnScreen bool
}

type CompassPosition struct {
	Location BlockPosition
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
	SegmentMask byte
	Segments [][]uint32
}

// Serverbound
