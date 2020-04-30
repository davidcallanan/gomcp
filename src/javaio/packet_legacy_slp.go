package javaio

import "bufio"
import "strconv"
import "unicode/utf16"

type LegacyStatusRequest struct {
}

type LegacyStatusResponse struct {
	ProtocolVersion int
	TextVersion string
	Description string
	MaxPlayers int
	OnlinePlayers int
}

func WriteLegacyStatusResponse(status LegacyStatusResponse, stream *bufio.Writer) {
	packetId := byte(0xff)
	protocolVersion := utf16.Encode([]rune(strconv.Itoa(status.ProtocolVersion)))
	textVersion := utf16.Encode([]rune(status.TextVersion))
	description := utf16.Encode([]rune(status.Description))
	maxPlayers := utf16.Encode([]rune(strconv.Itoa(status.MaxPlayers)))
	onlinePlayers := utf16.Encode([]rune(strconv.Itoa(status.OnlinePlayers)))
	dataLength := int16(3 + len(protocolVersion) + 1 + len(textVersion) + 1 + len(description) + 1 + len(onlinePlayers) + 1 + len(maxPlayers)) // potentially unsafe cast?
	magic := []byte { 0x00, 0xa7, 0x00, 0x31, 0x00, 0x00 }
	delimeter := []byte { 0x00, 0x00 }

	stream.WriteByte(packetId)
	WriteShort(dataLength, stream)
	stream.Write(magic)
	WriteUTF16(protocolVersion, stream)
	stream.Write(delimeter)
	WriteUTF16(textVersion, stream)
	stream.Write(delimeter)
	WriteUTF16(description, stream)
	stream.Write(delimeter)
	WriteUTF16(onlinePlayers, stream)
	stream.Write(delimeter)
	WriteUTF16(maxPlayers, stream)
}
