package javaio

import "bufio"
import "github.com/google/uuid"

type PlayerInfoAdd struct {
	Players []PlayerInfo
}

type PlayerInfo struct {
	Uuid uuid.UUID
	Username string
	Ping int32
}

func PacketId_PlayerInfo(protocol uint) int {
	// TODO: this is an approximation
	if protocol >= 0x0286 {
		// 1.15
		return 0x34
	} else {
		// 1.14
		return 0x33
	}
	// todo older versions
}

func WritePlayerInfoAdd(data PlayerInfoAdd, stream *bufio.Writer) {
	// WARNING: My instinctial conclusion is that the number of players on the tab list
	// cannot exceed the number of players online? I haven't managed to get more than one player

	WriteVarInt(0, stream) // action 0: add players
	WriteVarInt(int32(len(data.Players)), stream) // potentially unsafe cast?

	for _, player := range data.Players {
		WriteUuidBin(player.Uuid, stream)
		WriteString(player.Username, stream)
		WriteVarInt(0, stream) // property count; no properties for now
		WriteVarInt(0, stream) // gamemode survival; not worried about this for now
		WriteVarInt(player.Ping, stream)
		WriteBool(false, stream) // has display name; false for now
	}
}
