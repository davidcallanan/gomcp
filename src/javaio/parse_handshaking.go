package javaio

import "fmt"
import "bufio"

/**  There are no clientbound packets in this state.  **/
/**  All serverbound packets in this state are implemented.  **/

// Serverbound

func ParseHandshake(data *bufio.Reader) (result Handshake, err error) {
	protocol, err := ReadVarInt(data)
	if err != nil {
		return
	}

	serverAddress, err := ReadString(data, 256)
	if err != nil {
		return
	}

	serverPort, err := ReadUShort(data)
	if err != nil {
		return
	}

	nextStateId, err := ReadVarInt(data)
	if err != nil {
		return
	}

	var nextState State

	if (nextStateId == 1) {
		nextState = StateStatus
	} else if (nextStateId == 2) {
		nextState = StateLogin
	} else {
		err = MalformedPacketError { fmt.Sprintf("Unrecognized next state id %d in handshake", nextStateId) }
		return
	}

	result = Handshake {
		Protocol: protocol,
		ServerAddress: serverAddress,
		ServerPort: serverPort,
		NextState: nextState,
	}

	return
}
