package javaio

import "bufio"
import "encoding/json"
import "encoding/base64"

/**  Serverbound packet emission will not be available for the foreseeable future unless contributed by others.  **/
/**  All clientbound packets are implemented.  **/

// Clientbound

type statusJson struct {
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	FaviconPng string `json:"favicon,omitempty"`
	Version struct {
		Text string `json:"name"`
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

func EmitStatusResponse(status StatusResponse, result *bufio.Writer) {
	// Generate JSON
	jsonObj := statusJson {}
	jsonObj.Description.Text = status.Description
	jsonObj.Version.Text = status.VersionText
	jsonObj.Version.Protocol = status.VersionProtocol
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
	WriteString(string(jsonBytes), result)
}

func EmitPong(pong Pong, result *bufio.Writer) {
	WriteLong(pong.Payload, result)
}
