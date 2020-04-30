package javaio

// Clientbound

type KeepAlive struct {
	Payload int64
}

type ChunkData struct {
	X int32
	Z int32
	IsNew bool
	Sections [][]uint32
}

// Serverbound
