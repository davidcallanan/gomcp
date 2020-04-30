package javaio

import "bufio"
import "strconv"
import "unicode/utf16"

// TODO: figure out the protocol version number when this packet was introduced.

type VeryLegacyStatusRequest struct {
}

type VeryLegacyStatusResponse struct {
	// Color-coding is not supported. Description is treated as plain-text. Section character must not be used, otherwise there will be undefined behaviour.
	Description string
	MaxPlayers int
	OnlinePlayers int
}

func WriteVeryLegacyStatusResponse(status VeryLegacyStatusResponse, stream *bufio.Writer) {
	packetId := byte(0xff)
	description := utf16.Encode([]rune(status.Description))
	maxPlayers := utf16.Encode([]rune(strconv.Itoa(status.MaxPlayers)))
	onlinePlayers := utf16.Encode([]rune(strconv.Itoa(status.OnlinePlayers)))
	dataLength := int16(len(description) + 1 + len(onlinePlayers) + 1 + len(maxPlayers)) // potentially unsafe cast?
	delimeter := []byte { 0x00, 0xa7 }

	stream.WriteByte(packetId)
	WriteShort(dataLength, stream)
	WriteUTF16(description, stream)
	stream.Write(delimeter)
	WriteUTF16(onlinePlayers, stream)
	stream.Write(delimeter)
	WriteUTF16(maxPlayers, stream)
}
