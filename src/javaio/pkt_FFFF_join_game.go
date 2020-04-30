package javaio

import "bufio"

// TODO: figure out prot version

type JoinGame struct {
	EntityId int32
	Gamemode Gamemode
	Hardcore bool
	Dimension Dimension
	ViewDistance int32
	ReducedDebugInfo bool
	EnableRespawnScreen bool
}

func PacketId_JoinGame(protocol uint) int {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x26
	} else {
		// 1.14
		return 0x25
	}
	// older versions not supported
	return -1
}

func WriteJoinGame(data JoinGame, stream *bufio.Writer) {
	WriteInt(data.EntityId, stream)

	var gamemode byte
	switch data.Gamemode {
	case GamemodeSurvival:
		gamemode = 0
	case GamemodeCreative:
		gamemode = 1
	case GamemodeAdventure:
		gamemode = 2
	case GamemodeSpectator:
		gamemode = 3
	default:
		panic("Gamemode does not match one of non-invalid predefined enum types")
	}

	if data.Hardcore {
		// Enable hardcore flag
		gamemode |= 0x8
	}	

	WriteUByte(gamemode, stream)

	var dimension int32
	switch data.Dimension {
	case DimensionOverworld:
		dimension = 0
	case DimensionNether:
		dimension = -1
	case DimensionEnd:
		dimension = 1
	default:
		panic("Dimension does not match one of non-invalid predefined enum types")
	}

	WriteInt(dimension, stream)

	var hashedSeed int64 = 0 // seems kind of useless
	WriteLong(hashedSeed, stream)
	
	var maxPlayers byte = 0 // no longer utilized by client
	WriteUByte(maxPlayers, stream)

	var levelType string = "default" // seems kind of useless
	WriteString(levelType, stream)

	if data.ViewDistance > 32 {
		panic("View distance must not be greater than 32")
	}

	WriteVarInt(data.ViewDistance, stream)
	WriteBool(data.ReducedDebugInfo, stream)
	WriteBool(data.EnableRespawnScreen, stream)
}
