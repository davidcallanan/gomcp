package javaio

import "bufio"
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
	protocolVersion := utf16.Encode([]rune(string(status.ProtocolVersion)))
	textVersion := utf16.Encode([]rune(status.TextVersion))
	description := utf16.Encode([]rune(status.Description))
	maxPlayers := utf16.Encode([]rune(string(status.MaxPlayers)))
	onlinePlayers := utf16.Encode([]rune(string(status.OnlinePlayers)))
	dataLength := int16(6 + len(protocolVersion) * 2 + 2 + len(textVersion) * 2 + 2 + len(description) * 2 + 2 + len(onlinePlayers) * 2 + 2 + len(maxPlayers) * 2) // potentially unsafe cast?
	magic := []byte { 0x00, 0xa7, 0x00, 0x31, 0x00, 0x00 }
	delimeter := []byte {0x00, 0x00}

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
