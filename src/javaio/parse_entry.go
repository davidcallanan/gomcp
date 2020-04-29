package javaio

import "fmt"
import "bufio"

/**  Clientbound entry is not implemented.  **/

// Serverbound

func ParseServerboundPacketUncompressed(data *bufio.Reader, state State) (result interface{}, err error) {
	if state == StateDeterminingProtocol {
		isNetty, err_ := data.ReadByte()
		
		if err_ != nil {
			err = &MalformedPacketError { "First packet ended immediately with no data" }
			return
		}
		
		// This detection mechanism is not working correctly
		if isNetty == 0 || true {
			_ = data.UnreadByte()
			result = ProtocolDetermined {
				NextState: StateHandshaking,
			}
		} else {
			result = ProtocolDetermined {
				NextState: StatePreNetty,
			}
		}
		
		return
	}
	
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
	case StatePreNetty:
		println("Ignoring prenetty stream")
	default:
		panic("State does not match one of non-invalid predefined enum values")
	}

	return
}

func ParseServerboundPacketCompressed(data *bufio.Reader, state State) (result interface{}, err error) {
	panic("ParseServerboundPacketCompressed not implemented")
}
