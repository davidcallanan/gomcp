package javaio

/**  All types for this state are available.  **/

// Clientbound

type StatusResponse struct {
	Protocol int32
	Version string
	Description string
	FaviconPng []byte
	MaxPlayers int
	OnlinePlayers int
	PlayerSample []StatusResponsePlayer
}

type StatusResponsePlayer struct {
	Name string
	Uuid string
}

type Pong struct {
	Payload int64
}

// Serverbound

type StatusRequest struct {}

type Ping struct {
	Payload int64
}
