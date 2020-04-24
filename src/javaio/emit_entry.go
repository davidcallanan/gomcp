package javaio

import "bytes"
import "bufio"

/**  Serverbound packet emission is not implemented.  **/

// Clientbound

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

func EmitClientboundPacketCompressed(packet interface{}, state int, output *bufio.Writer) (err error) {
	panic("EmicClientboundPacketCompressed not implemented")
}
