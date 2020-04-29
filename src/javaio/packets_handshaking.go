package javaio

/**  There are no clientbound packets in this state.  **/
/**  All serverbound packets in this state are implemented.  **/

// Serverbound

type Handshake struct {
	ProtocolVersion int32
	ServerAddress string
	ServerPort uint16
	NextState State
}
