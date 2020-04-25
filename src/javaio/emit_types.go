package javaio

import "bufio"

/**  Not all types are implemented at this time.  **/

func EmitLong(value int64, result *bufio.Writer) (err error) {
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

	return
}

func EmitVarInt(value int32, result *bufio.Writer) (err error) {
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

	return
}

func EmitString(value string, result *bufio.Writer) (err error) {
	// TODO: int32 cast potentially unsafe?
	err = EmitVarInt(int32(len(value)), result)
	if err != nil {
		return
	}

	result.WriteString(value)
	return
}
