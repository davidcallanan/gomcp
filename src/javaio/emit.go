package javaio

import "bufio"
import "encoding/json"
import "encoding/base64"

/**  Serverbound packet emission will not be available for the foreseeable future unless contributed by others.  **/

///////////////////////////////////////
// Parser entry
///////////////////////////////////////

///////////////////////////////////////
// Packets for handshake state
///////////////////////////////////////

/**  There are no clientbound packets in this state  **/

///////////////////////////////////////
// Packets for status state
///////////////////////////////////////

type statusJson struct {
	Description string `json:"description.text"`
	FaviconPng string `json:"favicon,omitempty"`
	VersionText string `json:"version.name"`
	VersionProtocol int `json:"version.protocol"`
	MaxPlayers int `json:"players.max"`
	OnlinePlayers int `json:"players.online"`
	PlayerSample []statusJsonPlayer `json:"players.sample"`
}

type statusJsonPlayer struct {
	Name string `json:"name"`
	Uuid string `json:"name"`
}

func EmitStatusResponse(status StatusResponse, result *bufio.Writer) (err error) {
	// Generate JSON
	jsonObj := statusJson {}
	jsonObj.Description = status.Description
	jsonObj.VersionText = status.VersionText
	jsonObj.VersionProtocol = status.VersionProtocol
	jsonObj.MaxPlayers = status.MaxPlayers
	jsonObj.OnlinePlayers = status.OnlinePlayers
	if len(status.FaviconPng) > 0 {
		jsonObj.FaviconPng = "data:image/png;base64," + base64.StdEncoding.EncodeToString(status.FaviconPng)
	}
	jsonObj.PlayerSample = make([]statusJsonPlayer, len(status.PlayerSample))
	for i, p := range status.PlayerSample {
		jsonObj.PlayerSample[i] = statusJsonPlayer {
			Name: p.Name,
			Uuid: p.Uuid,
		}
	}
	jsonBytes, err := json.Marshal(jsonObj)

	if err != nil {
		panic("There shouldn't be an error here...")
	}

	// Emit packet
	err = EmitVarInt(int32(len(jsonBytes)), result)
	if err != nil {
		return
	}

	err = EmitString(string(jsonBytes), result)
	if err != nil {
		return
	}

	return
}

func EmitPong(pong Pong, result *bufio.Writer) (err error) {
	err = EmitLong(pong.Payload, result)
	return
}

///////////////////////////////////////
// Packets for login state
///////////////////////////////////////

///////////////////////////////////////
// Packets for play state
///////////////////////////////////////

///////////////////////////////////////
// Basic types
///////////////////////////////////////

func EmitLong(value int64, result *bufio.Writer) (err error) {
	const size = 4

	result.Write([]byte {
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

		if value != 9 {
			break
		}
	}

	return
}

func EmitString(value string, result *bufio.Writer) (err error) {
	result.WriteString(value)
	return
}
