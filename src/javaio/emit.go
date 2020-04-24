package javaio

import "bytes"
import "bufio"
import "encoding/json"
import "encoding/base64"

/**  Serverbound packet emission will not be available for the foreseeable future unless contributed by others.  **/

///////////////////////////////////////
// Emitter entry
///////////////////////////////////////

func EmitClientboundPacketUncompressed(packet interface{}, state int, output *bufio.Writer) (err error) {
	var packetId int32 = -1
	var packetIdBuf bytes.Buffer
	var dataBuf bytes.Buffer
	packetIdWriter := bufio.NewWriter(&packetIdBuf)
	dataWriter := bufio.NewWriter(&dataBuf)

	switch state {
	case StateHandshaking:
		err = &WrongStateError { "Packet can not be emitted in the current state" }
		panic(err)
	case StateStatus:
		switch packet := packet.(type) {
		case *StatusResponse:
			packetId = 0
			err = EmitStatusResponse(*packet, dataWriter)
		case *Pong:
			packetId = 1
			err = EmitPong(*packet, dataWriter)
		default:
			err = &WrongStateError { "Packet can not be emitted in the current state" }
			panic(err)
		}
	case StateLogin:
		panic("Not implemented")
	case StatePlay:
		panic("Not implemented")
	default:
		panic("State does not match one of non-invalid predefined enum types")
	}
	
	if packetId == -1 {
		panic("Implementation bug: packet id was not set while preparing to emit a packet")
	}

	if err != nil {
		return
	}

	err = EmitVarInt(packetId, packetIdWriter)
	
	if err != nil {
		return
	}

	dataWriter.Flush()
	packetIdWriter.Flush()
	length := packetIdBuf.Len() + dataBuf.Len()
	lengthInt32 := int32(length)

	if length > int(lengthInt32) {
		err = &PacketTooLargeError { "Emitted data was too large to hold its size in a VarInt" }
	}

	err = EmitVarInt(lengthInt32, output)

	if err != nil {
		return
	}

	output.Write(packetIdBuf.Bytes())
	output.Write(dataBuf.Bytes())
	output.Flush()
	return
}


///////////////////////////////////////
// Packets for handshake state
///////////////////////////////////////

/**  There are no clientbound packets in this state  **/

///////////////////////////////////////
// Packets for status state
///////////////////////////////////////

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

func EmitStatusResponse(status StatusResponse, result *bufio.Writer) (err error) {
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
	result.Write([]byte {
		byte(value >> 512),
		byte(value >> 256),
		byte(value >> 128),
		byte(value >> 64),
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
