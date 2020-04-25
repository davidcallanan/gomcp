package javaio

// Clientbound

type JoinGame struct {
	Eid int32
	Gamemode Gamemode
	Hardcore bool
	Dimension Dimension
	ViewDistance int32
	ReducedDebugInfo bool
	EnableRespawnScreen bool
}

// Serverbound
