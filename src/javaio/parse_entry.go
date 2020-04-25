package javaio

import "fmt"
import "bufio"

/**  Clientbound entry is not implemented.  **/

// Serverbound

func ParseServerboundPacketUncompressed(data *bufio.Reader, state int) (result interface{}, err error) {
	length, err := ParseVarInt(data)
	if err != nil {
		return
	}

	// Ensure one does not read past the length of the packet
	data = newReaderSlice(data, int(length))

	packetId, err := ParseVarInt(data)
	if err != nil {
		return
	}

	switch state {
	case StateHandshaking:
		switch packetId {
		case 0:
			result, err = ParseHandshake(data)
		default:
			err = &UnsupportedPayloadError { fmt.Sprintf("Unrecognized packet id %d", packetId) }
		}
	case StateStatus:
		switch packetId {
		case 0:
			result, err = ParseStatusRequest(data)
		case 1:
			result, err = ParsePing(data)
		default:
			err = &UnsupportedPayloadError { fmt.Sprintf("Unrecognized packet id %d", packetId) }
		}
	case StateLogin:
		switch packetId {
		case 0:
			result, err = ParseLoginStart(data)
		default:
			err = &UnsupportedPayloadError { fmt.Sprintf("Unrecognized packet id %d", packetId) }
		}
	case StatePlay:
		// panic("Not implemented")
	default:
		panic("State does not match one of non-invalid predefined enum values")
	}

	return
}

func ParseServerboundPacketCompressed(data []byte) (result interface{}, bytesProcessed int, err error) {
	panic("ParseServerboundPacketCompressed not implemented")
}
