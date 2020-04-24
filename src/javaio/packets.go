package javaio

// Handshake packets

type Handshake struct {
	ProtocolVersion int
	ServerAddress string
	ServerPort uint16
	NextState int
}

type LegacyStatusRequest struct {}

// Status packets

type StatusRequest struct {}

type Ping struct {
	Payload int64
}

type Pong struct {
	Payload int64
}

// Login packets

// Play packets
