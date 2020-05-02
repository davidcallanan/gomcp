package javaio

import "bufio"
import "encoding/json"
import "encoding/base64"

// Types

type Packet_0051_StatusRequest struct {
}

type Packet_0051_StatusResponse struct {
	Protocol int32
	Version string
	Description string
	FaviconPng []byte
	MaxPlayers int
	OnlinePlayers int
	PlayerSample []Packet_0051_StatusResponse_Player
}

type Packet_0051_StatusResponse_Player struct {
	Name string
	Uuid string
}

type Packet_0051_Ping struct {
	Payload int64
}

type Packet_0051_Pong struct {
	Payload int64
}

// Read

func Read_0051_StatusRequest(stream *bufio.Reader) (result Packet_0051_StatusRequest, err error) {
	result = Packet_0051_StatusRequest{}
	return
}

func Read_0051_Ping(stream *bufio.Reader) (result Packet_0051_Ping, err error) {
	payload, err := ReadLong(stream)
	if err != nil {
		return
	}

	result = Packet_0051_Ping {
		Payload: payload,
	}
	return
}

// Write

type statusJson struct {
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	FaviconPng string `json:"favicon,omitempty"`
	Version struct {
		Name string `json:"name"`
		Protocol int32 `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max int `json:"max"`
		Online int `json:"online"`
		Sample []statusJsonPlayer `json:"sample"`
	} `json:"players"`
}

type statusJsonPlayer struct {
	Name string `json:"name"`
	Uuid string `json:"id"`
}

func Write_0051_StatusResponse(status Packet_0051_StatusResponse, stream *bufio.Writer) {
	// Generate JSON
	jsonObj := statusJson {}
	jsonObj.Description.Text = status.Description
	jsonObj.Version.Name = status.Version
	jsonObj.Version.Protocol = status.Protocol
	jsonObj.Players.Max = status.MaxPlayers
	jsonObj.Players.Online = status.OnlinePlayers
	if len(status.FaviconPng) > 0 {
		jsonObj.FaviconPng = "data:image/png;base64," + base64.StdEncoding.EncodeToString(status.FaviconPng)
	}
	jsonObj.Players.Sample = make([]statusJsonPlayer, len(status.PlayerSample))
	for i, p := range status.PlayerSample {
		jsonObj.Players.Sample[i] = statusJsonPlayer {
			Name: p.Name,
			Uuid: p.Uuid,
		}
	}
	jsonBytes, err := json.Marshal(jsonObj)

	if err != nil {
		panic(err)
	}

	// Emit packet
	WriteString(string(jsonBytes), stream)
}

func Write_0051_Pong(pong Packet_0051_Pong, stream *bufio.Writer) {
	WriteLong(pong.Payload, stream)
}
