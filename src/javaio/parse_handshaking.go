package javaio

import "fmt"
import "bufio"

/**  There are no clientbound packets in this state.  **/
/**  All serverbound packets in this state are implemented.  **/

// Serverbound

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

	nextStateId, err := ParseVarInt(data)
	if err != nil {
		return
	}

	var nextState int

	if (nextStateId == 1) {
		nextState = StateStatus
	} else if (nextStateId == 2) {
		nextState = StateLogin
	} else {
		err = &MalformedPacketError { fmt.Sprintf("Unrecognized next state id %d in handshake", nextStateId) }
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
