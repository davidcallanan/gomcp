package javaio

import "bufio"
import "fmt"
import "unicode/utf8"

// Clientbound packets

/**  Will not be available for the foreseeable future unless contributed by others.  **/

// Serverbound packets

func ParseServerboundPacketUncompressed(data *bufio.Reader) (result interface{}, err error) {
	length, err := ParseVarInt(data)
	_ = length
	if err != nil {
		return
	}

	packetId, err := ParseVarInt(data)
	if err != nil {
		return
	}

	if packetId == 0 {
		if _, err := data.Peek(1); err == nil {
			// No more data in this packet
			panic("Not yet implemented!!")
		} else {
			result, err = ParseHandshake(data)
		}
	} else {
		err = &UnsupportedPayloadError { fmt.Sprintf("Unrecognized packet id %d", packetId) }
	}

	return
}

func ParseServerboundPacketCompressed(data []byte) (result interface{}, bytesProcessed int, err error) {
	panic("ParseServerboundPacketCompressed not implemented")
}

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

	nextState, err := ParseVarInt(data)
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

// Building blocks

func ParseVarInt(data *bufio.Reader) (result int, err error) {
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
		result |= int(uint(value) << uint(7 * idx))
		idx++

		if byte_ & 0b10000000 == 0 {
			break
		}
	}

	return
}

func ParseUnsignedShort(data *bufio.Reader) (result uint16, err error) {
	var buf [2]byte
	n, _ := data.Read(buf[:])

	if n < 2 {
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

	if strLength > maxStrLength {
		err = &MalformedPacketError { "String exceeded max rune count" } //*
		return
	}

	buf := make([]byte, strLength)
	n, _ := data.Read(buf)
	
	if (n != strLength) {
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
