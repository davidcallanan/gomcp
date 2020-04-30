package javaio

import "fmt"
import "bufio"

/**  Clientbound entry is not implemented.  **/

// Serverbound

func ParseServerboundPacketUncompressed(data *bufio.Reader, state State) (result interface{}, err error) {
	if state == StateDeterminingProtocol {
		// This pre-netty detection is pretty good.
		// It does however require that there are at least 2 bytes in the
		// first pre-netty packet, and that its second byte is not null.
		// This means a legacy SLP packet from Minecraft versions before
		// 1.4 is not detected, because its packet length is only 1 byte.
		// A hack is used to determine if only 1 byte is available.
		
		// This check breaks the intuitive purity contract.
		// This changes the contract to take into account the
		// number of buffered bytes. This new contract must be
		// tested within unit tests in the future.
		if data.Buffered() >= 2 {
			preview, _ := data.Peek(2)
		
			if len(preview) < 2 {
				err = &MalformedPacketError { "Stream ended abruptly" }
				return
			}

			if preview[1] == 0 {
				result = ProtocolDetermined {
					NextState: StateHandshaking,
				}
			} else {
				result = ProtocolDetermined {
					NextState: StatePreNetty,
				}
			}

			return
		} else {
			preview, _ := data.Peek(1)
			
			if len(preview) < 1 {
				err = &MalformedPacketError { "Stream ended abruptly" }
				return
			}

			// TODO: add various checks for other starting bytes in the legacy protocol
			if preview[0] == 0xfe {
				result = ProtocolDetermined {
					NextState: StateVeryPreNetty,
				}
			} else {
				result = ProtocolDetermined {
					NextState: StateHandshaking,
				}
			}

			return
		}
	} else if (state == StatePreNetty) {
		// Temporary hack-check
		b, _ := data.ReadByte()
		if b == 0xfe {
			result = LegacyStatusRequest {
			}
		}
		return
	} else if (state == StateVeryPreNetty) {
		// Temporary hack-check
		b, _ := data.ReadByte()
		if b == 0xfe {
			result = VeryLegacyStatusRequest {
			}
		}
		return
	}
	
	length, err := ReadVarInt(data)
	if err != nil {
		return
	}

	// Ensure one does not read past the length of the packet
	data = newReaderSlice(data, int(length))

	packetId, err := ReadVarInt(data)
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

func ParseServerboundPacketCompressed(data *bufio.Reader, state State) (result interface{}, err error) {
	panic("ParseServerboundPacketCompressed not implemented")
}
