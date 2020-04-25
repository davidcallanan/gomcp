package javaio

import "bufio"

/**  Serverbound packet emission is not implemented.  **/

// Clientbound

func EmitKeepAlive(data KeepAlive, result *bufio.Writer) {
	EmitLong(data.Payload, result)
}

func EmitJoinGame(data JoinGame, result *bufio.Writer) {
	EmitInt(data.Eid, result)

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

	EmitUnsignedByte(gamemode, result)

	var dimension int32
	switch data.Dimension {
	case DimensionOverworld:
		dimension = 0
	case DimensionNether:
		dimension = -1
	case DimensionEnd:
		dimension = 1
	default:
		panic("What sort of a dimension is that?")
	}

	EmitInt(dimension, result)

	var hashedSeed int64 = 0 // seems kind of useless
	EmitLong(hashedSeed, result)
	
	var maxPlayers byte = 0 // no longer utilized by client
	EmitUnsignedByte(maxPlayers, result)

	var levelType string = "default" // seems kind of useless
	EmitString(levelType, result)

	if data.ViewDistance > 32 {
		panic("Render distance must not be greater than 32")
	}

	EmitVarInt(data.ViewDistance, result)
	EmitBoolean(data.ReducedDebugInfo, result)
	EmitBoolean(data.EnableRespawnScreen, result)
}

func EmitPlayerPositionAndLook(data PlayerPositionAndLook, result *bufio.Writer) {
	EmitDouble(data.X, result)
	EmitDouble(data.Y, result)
	EmitDouble(data.Z, result)
	EmitFloat(data.Yaw, result)
	EmitFloat(data.Pitch, result)

	var flags byte

	if data.IsRelX {
		flags |= 0x01
	}
	if data.IsRelY {
		flags |= 0x02
	}
	if data.IsRelZ {
		flags |= 0x04
	}
	if data.IsRelYaw {
		flags |= 0x10
	}
	if data.IsRelPitch {
		flags |= 0x08
	}

	EmitUnsignedByte(flags, result)

	// Seems pointless for now.
	// Probably useful for interpolation, etc.
	var teleportId int32 = 0
	EmitVarInt(teleportId, result)
}

func EmitCompassPosition(compassPosition CompassPosition, result *bufio.Writer) {
	EmitBlockPosition(compassPosition.Location, result)
}
