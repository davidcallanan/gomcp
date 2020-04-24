package javaio

import "bufio"
import "fmt"
import "unicode/utf8"

/**  Clientbound packet parsing will not be available for the foreseeable future unless contributed by others.  **/

///////////////////////////////////////
// Parser entry
///////////////////////////////////////

func ParseServerboundPacketUncompressed(data *bufio.Reader, state int) (result interface{}, err error) {
	length, err := ParseVarInt(data)
	if err != nil {
		return
	}

	// Ensure one does not read past the length of the packet
	data = newReaderSlice(data, int(length))

	packetId, err := ParseVarInt(data)
	if err != nil {
		return
	}

	switch state {
	case StateHandshaking:
		switch packetId {
		case 0:
			result, err = ParseHandshake(data)
		default:
			err = &UnsupportedPayloadError { fmt.Sprintf("Unrecognized packet id %d", packetId) }
		}
	case StateStatus:
		switch packetId {
		case 0:
			result, err = ParseStatusRequest(data)
		case 1:
			result, err = ParsePing(data)
		default:
			err = &UnsupportedPayloadError { fmt.Sprintf("Unrecognized packet id %d", packetId) }
		}
	case StateLogin:
		panic("Not implemented")
	case StatePlay:
		panic("Not implemented")
	default:
		panic("State does not match one of non-invalid predefined enum values")
	}

	return
}

func ParseServerboundPacketCompressed(data []byte) (result interface{}, bytesProcessed int, err error) {
	panic("ParseServerboundPacketCompressed not implemented")
}

///////////////////////////////////////
// Packets for handshake state
///////////////////////////////////////

func ParseHandshake(data *bufio.Reader) (result Handshake, err error) {
	protocolVersion, err := ParseVarInt(data)
	if err != nil {
		return
	}

	serverAddress, err := ParseString(data, 256)
	if err != nil {
		return
	}

	serverPort, err := ParseUnsignedShort(data)
	if err != nil {
		return
	}

	nextStateId, err := ParseVarInt(data)
	if err != nil {
		return
	}

	var nextState int

	if (nextStateId == 1) {
		nextState = StateStatus
	} else if (nextStateId == 2) {
		nextState = StateLogin
	} else {
		err = &MalformedPacketError { fmt.Sprintf("Unrecognized next state id %d in handshake", nextStateId) }
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

///////////////////////////////////////
// Packets for status state
///////////////////////////////////////

func ParseStatusRequest(data *bufio.Reader) (result StatusRequest, err error) {
	result = StatusRequest{}
	return
}

func ParsePing(data *bufio.Reader) (result Ping, err error) {
	payload, err := ParseLong(data)
	if err != nil {
		return
	}

	result = Ping {
		Payload: payload,
	}

	return
}

///////////////////////////////////////
// Packets for play state
///////////////////////////////////////

///////////////////////////////////////
// Basic types
///////////////////////////////////////

func ParseVarInt(data *bufio.Reader) (result int32, err error) {
	maxLength := 5
	idx := 0

	for {
		byte_, readErr := data.ReadByte()

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

func ParseLong(data *bufio.Reader) (result int64, err error) {
	const size = 4

	var buf [size]byte
	n, _ := data.Read(buf[:])

	if n < size {
		err = &MalformedPacketError { "Long ended abruptly" }
		return
	}

	result = int64(buf[3]) + 256 * int64(buf[2]) + 65536 * int64(buf[1]) + 4294967296 * int64(buf[0])
	return
}

func ParseUnsignedShort(data *bufio.Reader) (result uint16, err error) {
	const size = 2

	var buf [size]byte
	n, _ := data.Read(buf[:])

	if n < size {
		err = &MalformedPacketError { "Unsigned short ended abruptly" }
		return
	}

	result = uint16(buf[1]) + 256 * uint16(buf[0])
	return
}

func ParseString(data *bufio.Reader, maxRuneCount int) (result string, err error) {
	maxStrLength := maxRuneCount * 4
	strLength, err := ParseVarInt(data)

	if err != nil {
		return
	}

	if int(strLength) > maxStrLength {
		err = &MalformedPacketError { "String exceeded max rune count" } //*
		return
	}

	buf := make([]byte, strLength)
	n, _ := data.Read(buf)
	
	if (n != int(strLength)) {
		err = &MalformedPacketError { "String ended abruptly" }
		return
	}

	result = string(buf)

	if utf8.RuneCountInString(result) > maxRuneCount {
		err = &MalformedPacketError { "String exceeded max rune count" }
		return
	}

	return
}
