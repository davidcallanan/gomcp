package javaio

// Handshake packets

type Handshake struct {
	ProtocolVersion int32
	ServerAddress string
	ServerPort uint16
	NextState int
}

type LegacyStatusRequest struct {}

// Status packets

type StatusRequest struct {}

type StatusResponse struct {
	Description string
	FaviconPng []byte
	VersionText string
	VersionProtocol int32
	MaxPlayers int
	OnlinePlayers int
	PlayerSample []StatusResponsePlayer
}

type StatusResponsePlayer struct {
	Name string
	Uuid string
}

type Ping struct {
	Payload int64
}

type Pong struct {
	Payload int64
}

// Login packets

// Play packets
