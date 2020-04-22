package javaio

import "fmt"
import "unicode/utf8"

// Packets

// TODO: this system needs to be switched to be stream-oriented rather than slice-oriented

func ParsePacketUncompressed(data []byte) (result interface{}, bytesProcessed int, err error) {
	length, bytes, err := ParseVarInt(data[bytesProcessed:])
	_ = length
	bytesProcessed += bytes

	if err != nil {
		return
	}

	packetId, bytes, err := ParseVarInt(data[bytesProcessed:])
	bytesProcessed += bytes

	if err != nil {
		return
	}

	if packetId == 0 {
		if len(data) == 0 {
			panic("Not yet implemented!!")
			// result, bytes, err := ParseRequest(data[bytesProcessed:])
			// bytesProcessed += bytes
		} else {
			result_, bytes, err_ := ParseHandshake(data[bytesProcessed:])
			result = result_
			err = err_
			bytesProcessed += bytes
		}
	} else {
		err = &UnsupportedPayloadError { fmt.Sprintf("Unrecognized packet id %d", packetId) }
	}

	return
}

func ParsePacketCompressed(data []byte) (result interface{}, bytesProcessed int, err error) {
	panic("ParsePacketCompressed not implemented")
}

func ParseHandshake(data []byte) (result Handshake, bytesProcessed int, err error) {
	protocolVersion, bytes, err := ParseVarInt(data[bytesProcessed:])
	bytesProcessed += bytes
	if err != nil {
		return
	}

	serverAddress, bytes, err := ParseString(data[bytesProcessed:], 32767)
	bytesProcessed += bytes
	if err != nil {
		return
	}

	serverPort, bytes, err := ParseUnsignedShort(data[bytesProcessed:])
	bytesProcessed += bytes
	if err != nil {
		return
	}

	nextState, bytes, err := ParseVarInt(data[bytesProcessed:])
	bytesProcessed += bytes
	if err != nil {
		return
	}

	result = Handshake {
		ProtocolVersion: protocolVersion,
		ServerAddress: serverAddress,
		ServerPort: serverPort,
		NextState: nextState,
	}

	return
}

// Basic components

func ParseVarInt(data []byte) (result int, bytesProcessed int, err error) {
	maxLength := 5
	idx := 0

	for {
		if idx >= len(data) {
			err = &MalformedPacketError { "VarInt ended abruptly" }
			return
		} else if idx >= maxLength {
			err = &MalformedPacketError { "VarInt exceeded max length "}
			return
		}

		byte_ := data[idx]
		value := byte_ & 0b01111111
		result |= int(uint(value) << uint(7 * idx))
		idx++

		if byte_ & 0b10000000 == 0 {
			break
		}
	}

	bytesProcessed = idx
	return
}

func ParseUnsignedShort(data []byte) (result uint16, bytesProcessed int, err error) {
	if len(data) < 2 {
		err = &MalformedPacketError { "Unsigned short ended abruptly" }
		return
	}

	result = uint16(data[1]) + 256 * uint16(data[0])
	bytesProcessed = 2
	return
}

func ParseString(data []byte, maxRuneCount int) (result string, bytesProcessed int, err error) {
	maxStrLength := maxRuneCount * 4
	strLength, varIntLength, err := ParseVarInt(data)
	packetLength := strLength + varIntLength

	if err != nil {
		return
	}

	if strLength > maxStrLength {
		err = &MalformedPacketError { "String exceeded max rune count" }
		return
	}

	result = string(data[varIntLength:packetLength])

	if len(result) != strLength {
		err = &MalformedPacketError { "String ended abruptly" }
		return
	}

	if utf8.RuneCountInString(result) > maxRuneCount {
		err = &MalformedPacketError { "String exceeded max rune count" }
		return
	}

	bytesProcessed = packetLength
	return
}
