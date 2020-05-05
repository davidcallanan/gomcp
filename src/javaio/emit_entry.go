package javaio

import "bytes"
import "bufio"

/**  Serverbound packet emission is not implemented.  **/

// Clientbound

func EmitClientboundPacketUncompressed(packet interface{}, ctx ClientContext, output *bufio.Writer) {
	if ctx.State == StatePreNetty {
		Write_002E_StatusResponse(packet.(Packet_002E_StatusResponse), output)
		output.Flush()
		return
	} else if ctx.State == StateVeryPreNetty {
		WriteVeryLegacyStatusResponse(packet.(VeryLegacyStatusResponse), output)
		output.Flush()
		return
	}
	
	var packetId int32 = -1
	var packetIdBuf bytes.Buffer
	var dataBuf bytes.Buffer
	packetIdWriter := bufio.NewWriter(&packetIdBuf)
	dataWriter := bufio.NewWriter(&dataBuf)

	switch ctx.State {
	case StateHandshaking:
		panic("Packet cannot be emitted in handshaking state")
	case StateStatus:
		switch packet := packet.(type) {
		case Packet_0051_StatusResponse:
			packetId = 0x00
			Write_0051_StatusResponse(packet, dataWriter)
		case Packet_0051_Pong:
			packetId = 0x01
			Write_0051_Pong(packet, dataWriter)
		default:
			panic("Packet cannot be emitted in status state")
		}
	case StateLogin:
		switch packet := packet.(type) {
		case LoginSuccess:
			packetId = 0x02
			EmitLoginSuccess(packet, dataWriter)	
		default:
			panic("Packet cannot be emitted in login state")
		}
	case StatePlay:
		switch packet := packet.(type) {
		case KeepAlive:
			packetId = int32(PacketId_KeepAlive(ctx.Protocol))
			WriteKeepAlive(packet, dataWriter)
		case JoinGame:
			packetId = int32(PacketId_JoinGame(ctx.Protocol))
			WriteJoinGame(packet, ctx, dataWriter)	
		case CompassPosition:
			packetId = int32(PacketId_CompassPosition(ctx.Protocol))
			WriteCompassPosition(packet, dataWriter)
		case PlayerPositionAndLook:
			packetId = int32(PacketId_PlayerPositionAndLook(ctx.Protocol))
			WritePlayerPositionAndLook(packet, dataWriter)
		case ChunkData:
			packetId = int32(PacketId_ChunkData(ctx.Protocol))
			WriteChunkData(packet, ctx, dataWriter)
		case PlayerInfoAdd:
			packetId = int32(PacketId_PlayerInfo(ctx.Protocol))
			WritePlayerInfoAdd(packet, dataWriter)
		case Packet_SpawnPlayer:
			packetId = int32(PacketId_SpawnPlayer(ctx.Protocol))
			Write_SpawnPlayer(packet, ctx, dataWriter)
		case Packet_EntityTranslate:
			packetId = int32(PacketId_EntityTranslate(ctx.Protocol))
			Write_EntityTranslate(packet, dataWriter)
		case Packet_EntityVelocity:
			packetId = int32(PacketId_EntityVelocity(ctx.Protocol))
			Write_EntityVelocity(packet, dataWriter)
		default:
			panic("Packet cannot be emitted in play state (likely because not implemented)")
		}
	default:
		panic("State does not match one of non-invalid predefined enum types")
	}
	
	if packetId == -1 {
		panic("Internal package bug: packet id was not set while preparing to emit a packet")
	}

	WriteVarInt(packetId, packetIdWriter)
	dataWriter.Flush()
	packetIdWriter.Flush()
	length := packetIdBuf.Len() + dataBuf.Len()
	lengthInt32 := int32(length)

	if length > int(lengthInt32) {
		panic("Emitted packet data was too large to hold its size in VarInt")
	}

	WriteVarInt(lengthInt32, output)
	output.Write(packetIdBuf.Bytes())
	output.Write(dataBuf.Bytes())
	output.Flush()
}

func EmitClientboundPacketCompressed(packet interface{}, ctx ClientContext, output *bufio.Writer) {
	panic("EmicClientboundPacketCompressed not implemented")
}
