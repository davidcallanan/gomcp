package javaio

import "bufio"
import "bytes"

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

func EmitCompassPosition(compassPosition CompassPosition, result *bufio.Writer) {
	EmitBlockPosition(compassPosition.Location, result)
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

func EmitChunkData(chunk ChunkData, result *bufio.Writer) {
	EmitInt(chunk.X, result)
	EmitInt(chunk.Z, result)
	EmitBoolean(chunk.IsNew, result)
	EmitVarInt(int32(chunk.SegmentMask), result) // TODO: enforce correct section mask
	
	// Hacked in NBT for heightmaps
	// Maybe we have to send the length as well?
	result.Write([]byte {
		10, // Compound start
		0,  // Length of compound name (1/2)
		0,  // Length of compound name (2/2)
		12, // Start long array
		0,  // Length of array name (1/2)
		15, // Length of array name (2/2)
	})
	result.Write([]byte("MOTION_BLOCKING")) // Array name
	EmitInt(36, result) // Length of array
	for i := 0; i < 288; i++ {
		result.WriteByte(0xFE) // arbitrary value for the heightmap
	}
	result.WriteByte(0) // Compound end

	if chunk.IsNew {
		// Set biome to void for the time being
		for i := 0; i < 1024; i++ {
			EmitInt(127, result)
		}
	}

	var dataBuf bytes.Buffer
	dataWriter := bufio.NewWriter(&dataBuf)
	
	for _, segment := range chunk.Segments {
		EmitChunkSegmentData(segment, dataWriter)
	}

	dataWriter.Flush()
	EmitVarInt(int32(dataBuf.Len()), result) // potentially unsafe cast
	result.Write(dataBuf.Bytes())

	EmitVarInt(0, result) // no block entities
}

func EmitChunkSegmentData(blocks []uint32, result *bufio.Writer) {
	const bitsPerBlock = 14

	EmitShort(4096, result) // block count
	EmitUnsignedByte(bitsPerBlock, result)
	// No palette because bits per block is full (can be optimized in future)

	if len(blocks) != 4096 {
		panic("There must be exactly 4096 blocks in each chunk segment")
	}

	bitLength := len(blocks) * bitsPerBlock
	length := (bitLength + bitLength % 64) / 64

	EmitVarInt(int32(length), result)

	currLong := uint64(0)
	currLongBit := uint64(0)

	for _, block := range blocks {
		// take each bit from the block
		for i := 0; i < bitsPerBlock; i++ {
			// extract bit
			bit := (block >> i) & 1
			// save bit into long
			currLong |= (uint64(bit) << currLongBit)
			currLongBit++
			// when long is filled, emit and reset
			if currLongBit == 64 {
				EmitLong(int64(currLong), result) // potentially unsafe cast
				currLong = 0
				currLongBit = 0
			}
		}
	}

	// flush any leftover that still needs to be emitted
	if currLongBit > 0 {
		EmitLong(int64(currLong), result) // potentially unsafe cast
	}
}
