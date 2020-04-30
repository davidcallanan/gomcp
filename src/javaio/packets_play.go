package javaio

// Clientbound

type ChunkData struct {
	X int32
	Z int32
	IsNew bool
	Sections [][]uint32
}

// Serverbound
