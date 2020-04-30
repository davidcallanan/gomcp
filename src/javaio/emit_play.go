package javaio

import "bufio"
import "bytes"

/**  Serverbound packet emission is not implemented.  **/

// Clientbound

func EmitKeepAlive(data KeepAlive, result *bufio.Writer) {
	WriteLong(data.Payload, result)
}

func EmitPlayerPositionAndLook(data PlayerPositionAndLook, result *bufio.Writer) {
	WriteDouble(data.X, result)
	WriteDouble(data.Y, result)
	WriteDouble(data.Z, result)
	WriteFloat(data.Yaw, result)
	WriteFloat(data.Pitch, result)

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

	WriteUByte(flags, result)

	// Seems pointless for now.
	// Probably useful for interpolation, etc.
	var teleportId int32 = 0
	WriteVarInt(teleportId, result)
}

func EmitChunkData(chunk ChunkData, result *bufio.Writer) {
	sectionMask := int32(0)

	for i, section := range chunk.Sections {
		if i >= 8 {
			return
		}

		if len(section) != 0 {
			sectionMask |= 1 << i
		}
	}

	WriteInt(chunk.X, result)
	WriteInt(chunk.Z, result)
	WriteBool(chunk.IsNew, result)
	WriteVarInt(sectionMask, result)
	
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
	WriteInt(36, result) // Length of array
	for i := 0; i < 288 / 9; i++ {
		// arbitrary value for the heightmap
		result.WriteByte(0b00100000)
		result.WriteByte(0b00010000)
		result.WriteByte(0b00001000)
		result.WriteByte(0b00000100)
		result.WriteByte(0b00000010)
		result.WriteByte(0b00000001)
		result.WriteByte(0b00000000)
		result.WriteByte(0b10000000)
		result.WriteByte(0b01000000)
	}
	result.WriteByte(0) // Compound end

	if chunk.IsNew {
		// Set biome to void for the time being
		for i := 0; i < 1024; i++ {
			WriteInt(127, result)
		}
	}

	var dataBuf bytes.Buffer
	dataWriter := bufio.NewWriter(&dataBuf)
	
	for i, section := range chunk.Sections {
		if i >= 8 {
			break
		}

		if len(section) != 0 {
			EmitChunkSectionData(section, dataWriter)
		}
	}

	dataWriter.Flush()
	WriteVarInt(int32(dataBuf.Len()), result) // potentially unsafe cast
	result.Write(dataBuf.Bytes())

	WriteVarInt(0, result) // no block entities
}

func EmitChunkSectionData(blocks []uint32, result *bufio.Writer) {
	const bitsPerBlock = 14

	WriteShort(4096, result) // block count
	WriteUByte(bitsPerBlock, result)
	// No palette because bits per block is full (can be optimized in future)

	if len(blocks) != 4096 {
		panic("There must be exactly 4096 blocks in each chunk section")
	}

	bitLength := len(blocks) * bitsPerBlock
	length := (bitLength + bitLength % 64) / 64

	WriteVarInt(int32(length), result)

	currLong := uint64(0)
	start := uint64(0)

	for _, block := range blocks {
		var b uint64 = uint64(block) & ((1 << bitsPerBlock) - 1)
		currLong |= b << start

		if start + bitsPerBlock >= 64 {
			WriteLong(int64(currLong), result)
			currLong = 0
			currLong |= b >> (64 - start)
			start -= 64
		}

		start += bitsPerBlock
	}

	// currLong := uint64(0)
	// currLongBit := uint64(0)

	// for _, block := range blocks {
	// 	// take each bit from the block
	// 	for i := 0; i < bitsPerBlock; i++ {
	// 		// extract bit
	// 		bit := (block >> i) & 1
	// 		// save bit into long
	// 		currLong |= uint64(bit) << currLongBit
	// 		currLongBit++
	// 		// when long is filled, emit and reset
	// 		if currLongBit == 64 {
	// 			EmitLong(int64(currLong), result) // potentially unsafe cast
	// 			fmt.Printf("%064b\n", currLong)
	// 			currLong = 0
	// 			currLongBit = 0
	// 		}
	// 	}
	// }

	// // ensure no leftover bits
	// if currLongBit > 0 {
	// 	panic("Shouldn't reach this point")
	// }
}
