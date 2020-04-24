package javaio

import "bufio"

/**  Clientbound parsing will not be available for the foreseeable future unless contributed by others.  **/
/**  All serverbound packets in this state are implemented.  **/

// Serverbound

func ParseStatusRequest(data *bufio.Reader) (result StatusRequest, err error) {
	result = StatusRequest{}
	return
}

func ParsePing(data *bufio.Reader) (result Ping, err error) {
	payload, err := ParseLong(data)
	if err != nil {
		return
	}

	result = Ping {
		Payload: payload,
	}

	return
}
