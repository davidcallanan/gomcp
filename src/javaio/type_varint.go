package javaio

import "bufio"

func ReadVarInt(stream *bufio.Reader) (result int32, err error) {
	maxLength := 5
	idx := 0

	for {
		byte_, readErr := stream.ReadByte()

		if readErr != nil {
			err = &MalformedPacketError { "VarInt ended abruptly" }
			return
		} else if idx >= maxLength {
			err = &MalformedPacketError { "VarInt exceeded max length "}
			return
		}

		value := byte_ & 0b01111111
		result |= int32(uint(value) << uint(7 * idx))
		idx++

		if byte_ & 0b10000000 == 0 {
			break
		}
	}

	return
}

func WriteVarInt(value int32, stream *bufio.Writer) {
	for {
		byte_ := byte(value & 0b01111111)
		value >>= 7
		if value != 0 {
			byte_ |= 0b10000000
		}

		stream.WriteByte(byte_)

		if value == 0 {
			break
		}
	}
}
