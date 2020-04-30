package javaio

import "bufio"
import "strconv"
import "unicode/utf16"

type T_002E_StatusRequest struct {
}

type T_002E_StatusResponse struct {
	Protocol int
	Version string
	Description string
	MaxPlayers int
	OnlinePlayers int
}

func Write_002E_StatusResponse(status T_002E_StatusResponse, stream *bufio.Writer) {
	packetId := byte(0xff)
	protocol := utf16.Encode([]rune(strconv.Itoa(status.Protocol)))
	version := utf16.Encode([]rune(status.Version))
	description := utf16.Encode([]rune(status.Description))
	maxPlayers := utf16.Encode([]rune(strconv.Itoa(status.MaxPlayers)))
	onlinePlayers := utf16.Encode([]rune(strconv.Itoa(status.OnlinePlayers)))
	dataLength := int16(3 + len(protocol) + 1 + len(version) + 1 + len(description) + 1 + len(onlinePlayers) + 1 + len(maxPlayers)) // potentially unsafe cast?
	magic := []byte { 0x00, 0xa7, 0x00, 0x31, 0x00, 0x00 }
	delimeter := []byte { 0x00, 0x00 }

	stream.WriteByte(packetId)
	WriteShort(dataLength, stream)
	stream.Write(magic)
	WriteUTF16(protocol, stream)
	stream.Write(delimeter)
	WriteUTF16(version, stream)
	stream.Write(delimeter)
	WriteUTF16(description, stream)
	stream.Write(delimeter)
	WriteUTF16(onlinePlayers, stream)
	stream.Write(delimeter)
	WriteUTF16(maxPlayers, stream)
}
