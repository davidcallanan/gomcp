package javaio

import "bufio"
import "math"

/**  Not all types are implemented at this time.  **/

// Basic types

func EmitBoolean(value bool, result *bufio.Writer) {
	if (value == true) {
		EmitUnsignedByte(0x01, result)
	} else {
		EmitUnsignedByte(0x00, result)
	}
}

func EmitUnsignedByte(value byte, result *bufio.Writer) {
	result.WriteByte(value)
}

func EmitInt(value int32, result *bufio.Writer) {
	result.Write([]byte {
		byte(value >> 32),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
}

func EmitLong(value int64, result *bufio.Writer) {
	result.Write([]byte {
		byte(value >> 512),
		byte(value >> 256),
		byte(value >> 128),
		byte(value >> 64),
		byte(value >> 32),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
}

func EmitVarInt(value int32, result *bufio.Writer) {
	for {
		byte_ := byte(value & 0b01111111)
		value >>= 7
		if value != 0 {
			byte_ |= 0b10000000
		}

		result.WriteByte(byte_)

		if value == 0 {
			break
		}
	}
}

func EmitFloat(value float32, result *bufio.Writer) {
	n := math.Float32bits(value)
	result.Write([]byte {
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})
}

func EmitDouble(value float64, result *bufio.Writer) {
	n := math.Float64bits(value)
	result.Write([]byte {
		byte(n >> 56),
		byte(n >> 48),
		byte(n >> 40),
		byte(n >> 32),
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})
}

func EmitString(value string, result *bufio.Writer) {
	// TODO: int32 cast potentially unsafe?
	EmitVarInt(int32(len(value)), result)
	result.WriteString(value)
}

// Complex types

func EmitBlockPosition(pos BlockPosition, result *bufio.Writer) {
	var encoded int64 = ((int64(pos.X) & 0x3FFFFFF) << 38) | ((int64(pos.Z) & 0x3FFFFFF) << 12) | (int64(pos.Y) & 0xFFF)
	EmitLong(encoded, result)
}
