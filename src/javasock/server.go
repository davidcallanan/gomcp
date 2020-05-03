package javasock

import "io"
import "time"
import "bufio"
import "github.com/davidcallanan/gomcp/javaio"
import "github.com/google/uuid"

type client struct {
	id int
	ctx javaio.ClientContext
	input *bufio.Reader
	output *bufio.Writer
	closeCallback func()
	isClosed bool
}

func (client *client) close() {
	client.closeCallback()
	client.isClosed = true
}

type Server struct {
	handleStatusRequestV1 func(id int) StatusResponseV1
	handleStatusRequestV2 func(id int) StatusResponseV2
	handleStatusRequestV3 func(id int) StatusResponseV3
	handlePlayerJoinRequest func(data PlayerJoinRequest) PlayerJoinResponse
	handlePlayerJoin func(id int)
}

func NewServer() Server {
	return Server {
	}
}

type StatusResponseV1 struct {
	// Color-coding is not supported.
	// Description is treated as plain-text.
	// Section character must not be used, otherwise there will be undefined behaviour.
	PreventResponse bool
	Description string
	MaxPlayers int
	OnlinePlayers int
}

type StatusResponseV2 struct {
	PreventResponse bool
	IsClientSupported bool
	Version string
	Description string
	MaxPlayers int
	OnlinePlayers int
}

type StatusResponseV3 struct {
	PreventResponse bool
	IsClientSupported bool
	Version string
	Description string
	FaviconPng []byte
	MaxPlayers int
	OnlinePlayers int
	PlayerSample []string
}

func (server *Server) OnStatusRequestV1(handler func(id int) StatusResponseV1) {
	server.handleStatusRequestV1 = handler
}

func (server *Server) OnStatusRequestV2(handler func(id int) StatusResponseV2) {
	server.handleStatusRequestV2 = handler
}

func (server *Server) OnStatusRequestV3(handler func(id int) StatusResponseV3) {
	server.handleStatusRequestV3 = handler
}

type PlayerJoinRequest struct {
	ClientsideUsername string
}

type PlayerJoinResponse struct {
	PreventResponse bool
	Uuid uuid.UUID
}

func (server *Server) OnPlayerJoinRequest(handler func(data PlayerJoinRequest) PlayerJoinResponse) {
	server.handlePlayerJoinRequest = handler
}

func (server *Server) OnPlayerJoin(handler func(id int)) {
	server.handlePlayerJoin = handler
}

func (server *Server) AddConnection(id int, input io.Reader, output io.Writer, closeCallback func()) {
	client := &client {
		id: id,
		ctx: javaio.InitialClientContext,
		input: bufio.NewReader(input),
		output: bufio.NewWriter(output),
		closeCallback: closeCallback,
	}

	go func() {
		for !client.isClosed {
			server.handleReceive(client)
		}
	}()

	go func() {
		timer := time.NewTicker(time.Second * 20)

		for now := range timer.C {
			if client.isClosed {
				break
			}
			if client.ctx.State != javaio.StatePlay {
				continue
			}

			client.SendPacket(javaio.KeepAlive {
				Payload: now.Unix(),
			})
		}
	}()
}

func (client *client) SendPacket(packet interface{}) {
	javaio.EmitClientboundPacketUncompressed(packet, client.ctx, client.output)
}

func (server *Server) handleReceive(client *client) {
	packet, err := javaio.ParseServerboundPacketUncompressed(client.input, client.ctx.State)

	if err != nil {
		switch err.(type) {
		case javaio.UnsupportedPayloadError:
			println("Unsupported payload from client")
			return
		case javaio.MalformedPacketError:
			println("Malformed packet from client.. closing connection")
			client.close()
			return
		default:
			panic(err)
		}
	}

	switch packet := packet.(type) {
		// Determining Protocol
	case javaio.ProtocolDetermined:
		server.ProcessProtocolDetermined(client, packet)

		// Handshaking
	case javaio.Handshake:
		server.ProcessHandshake(client, packet)

		// Status
	case javaio.Packet_0051_StatusRequest:
		server.ProcessStatusRequest(client, packet)
	case javaio.Packet_0051_Ping:
		server.ProcessPing(client, packet)

		// Login
	case javaio.LoginStart:
		server.ProcessLoginStart(client, packet)

		// Pre-Netty
	case javaio.Packet_002E_StatusRequest:
		server.ProcessLegacyStatusRequest(client, packet)

		// Very Pre-Netty
	case javaio.VeryLegacyStatusRequest:
		server.ProcessVeryLegacyStatusRequest(client, packet)

		// Default
	default:
		// println("Unrecognized packet type")
	}
}

func (server *Server) ProcessProtocolDetermined(client *client, data javaio.ProtocolDetermined) {
	client.ctx.State = data.NextState
}

func (server *Server) ProcessHandshake(client *client, handshake javaio.Handshake) {
	client.ctx.Protocol = javaio.DecodePostNettyVersion(handshake.Protocol)
	client.ctx.State = handshake.NextState
}

func (server *Server) ProcessStatusRequest(client *client, _ javaio.Packet_0051_StatusRequest) {
	if server.handleStatusRequestV3 == nil {
		return
	}

	res := server.handleStatusRequestV3(client.id)

	if (res.PreventResponse) {
		return
	}

	protocol := int32(0)
	if (res.IsClientSupported) {
		protocol = javaio.EncodePostNettyVersion(client.ctx.Protocol)
	}

	playerSample := make([]javaio.Packet_0051_StatusResponse_Player, len(res.PlayerSample), len(res.PlayerSample))

	for i, text := range res.PlayerSample {
		playerSample[i] = javaio.Packet_0051_StatusResponse_Player {
			Name: text,
			Uuid: "65bd239f-89f2-4cc7-ae8b-bb625525904e",
		}
	}

	client.SendPacket(javaio.Packet_0051_StatusResponse {
		Protocol: protocol,
		Version: res.Version,
		Description: res.Description,
		MaxPlayers: res.MaxPlayers,
		OnlinePlayers: res.OnlinePlayers,
		PlayerSample: playerSample,
	})
}

func (server *Server) ProcessLegacyStatusRequest(client *client, _ javaio.Packet_002E_StatusRequest) {
	if server.handleStatusRequestV2 == nil {
		return
	}

	res := server.handleStatusRequestV2(client.id)

	if (res.PreventResponse) {
		return
	}

	protocol := 0
	if (res.IsClientSupported) {
		protocol = int(client.ctx.Protocol)
	}

	client.SendPacket(javaio.Packet_002E_StatusResponse {
		Protocol: protocol,
		Version: res.Version,
		Description: res.Description,
		MaxPlayers: res.MaxPlayers,
		OnlinePlayers: res.OnlinePlayers,
	})
}

func (server *Server) ProcessVeryLegacyStatusRequest(client *client, _ javaio.VeryLegacyStatusRequest) {
	if server.handleStatusRequestV1 == nil {
		return
	}
	
	res := server.handleStatusRequestV1(client.id)

	if (res.PreventResponse) {
		return
	}

	client.SendPacket(javaio.VeryLegacyStatusResponse {
		Description: res.Description,
		MaxPlayers: res.MaxPlayers,
		OnlinePlayers: res.OnlinePlayers,
	})
}

func (server *Server) ProcessPing(client *client, ping javaio.Packet_0051_Ping) {
	client.SendPacket(javaio.Packet_0051_Pong {
		Payload: ping.Payload,
	})
}

func (server *Server) ProcessLoginStart(client *client, data javaio.LoginStart) {
	res := server.handlePlayerJoinRequest(PlayerJoinRequest {
		ClientsideUsername: data.ClientsideUsername,
	})

	if res.PreventResponse {
		return
	}
	
	playerUuid := res.Uuid

	client.SendPacket(javaio.LoginSuccess {
		Uuid: playerUuid,
		Username: data.ClientsideUsername,
	})

	client.ctx.State = javaio.StatePlay

	client.SendPacket(javaio.JoinGame {
		EntityId: 0,
		Gamemode: javaio.GamemodeCreative,
		Hardcore: false,
		Dimension: javaio.DimensionOverworld,
		ViewDistance: 1,
		ReducedDebugInfo: false,
		EnableRespawnScreen: false,
	})

	client.SendPacket(javaio.CompassPosition {
		Location: javaio.BlockPosition { X: 0, Y: 64, Z: 0 },
	})

	client.SendPacket(javaio.PlayerPositionAndLook {
		X: 0, Y: 64, Z: 0, Yaw: 0, Pitch: 0,
	})

	var blocksA [4096]uint32
	var blocksB [4096]uint32
	var blocksC [4096]uint32

	for i := range blocksA {
		if i < 256 {
			blocksA[i] = 33
		} else {
			blocksA[i] = 1
		}
	}

	for i := range blocksB {
		blocksB[i] = 1
	}

	for i := range blocksC {
		if i > 4095 - 256 {
			blocksC[i] = 9
		} else {
			blocksC[i] = 10
		}
	}

	for x := -3; x <= 3; x++ {
		for z := -3; z <= 3; z++ {
			client.SendPacket(javaio.ChunkData {
				X: int32(x), Z: int32(z), IsNew: true,
				Sections: [][]uint32 { nil, blocksA[:], blocksB[:], blocksC[:] },
			})
		}
	}

	if server.handlePlayerJoin != nil {
		server.handlePlayerJoin(client.id)
	}

	client.SendPacket(javaio.PlayerInfoAdd {
		Players: []javaio.PlayerInfo {
			{ Uuid: uuid.New(), Username: "JohnDoe", Ping: 0 },
			{ Uuid: uuid.New(), Username: "CatsEyebrows", Ping: 5 },
			{ Uuid: uuid.New(), Username: "ElepantNostrel23", Ping: 500 },
		},
	})
}

func (server *Server) SpawnPlayer() {
	
}
