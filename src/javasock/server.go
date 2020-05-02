package javasock

import "io"
import "time"
import "bufio"
import "github.com/davidcallanan/gomcp/javaio"
import "github.com/google/uuid"

type client struct {
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
	handleStatusRequestV1 func() StatusResponseV1
	handleStatusRequestV2 func() StatusResponseV2
	handleStatusRequestV3 func() StatusResponseV3
	handlePlayerJoin func(uuid uint32, clientsideUsername string)
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

func (server *Server) OnStatusRequestV1(callback func() StatusResponseV1) {
	server.handleStatusRequestV1 = callback
}

func (server *Server) OnStatusRequestV2(callback func() StatusResponseV2) {
	server.handleStatusRequestV2 = callback
}

func (server *Server) OnStatusRequestV3(callback func() StatusResponseV3) {
	server.handleStatusRequestV3 = callback
}

func (server *Server) OnPlayerJoin(callback func(uuid uint32, clientsideUsername string)) {
	server.handlePlayerJoin = callback
}

func (server *Server) AddConnection(input io.Reader, output io.Writer, closeCallback func()) {
	client := &client {
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

			client.SendPacket(&javaio.KeepAlive {
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
		case *javaio.UnsupportedPayloadError:
			println("Unsupported payload from client")
			return
		case *javaio.MalformedPacketError:
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
	case javaio.StatusRequest:
		server.ProcessStatusRequest(client, packet)
	case javaio.Ping:
		server.ProcessPing(client, packet)

		// Login
	case javaio.LoginStart:
		server.ProcessLoginStart(client, packet)

		// Pre-Netty
	case javaio.T_002E_StatusRequest:
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

func (server *Server) ProcessStatusRequest(client *client, _ javaio.StatusRequest) {
	res := server.handleStatusRequestV3()

	if (res.PreventResponse) {
		return
	}

	protocol := int32(0)
	if (res.IsClientSupported) {
		protocol = javaio.EncodePostNettyVersion(client.ctx.Protocol)
	}

	playerSample := make([]javaio.StatusResponsePlayer, len(res.PlayerSample), len(res.PlayerSample))

	for i, text := range res.PlayerSample {
		playerSample[i] = javaio.StatusResponsePlayer {
			Name: text,
			Uuid: "65bd239f-89f2-4cc7-ae8b-bb625525904e",
		}
	}

	client.SendPacket(&javaio.StatusResponse {
		Protocol: protocol,
		Version: res.Version,
		Description: res.Description,
		MaxPlayers: res.MaxPlayers,
		OnlinePlayers: res.OnlinePlayers,
		PlayerSample: playerSample,
	})
}

func (server *Server) ProcessLegacyStatusRequest(client *client, _ javaio.T_002E_StatusRequest) {
	res := server.handleStatusRequestV2()

	if (res.PreventResponse) {
		return
	}

	protocol := 0
	if (res.IsClientSupported) {
		protocol = int(client.ctx.Protocol)
	}

	client.SendPacket(&javaio.T_002E_StatusResponse {
		Protocol: protocol,
		Version: res.Version,
		Description: res.Description,
		MaxPlayers: res.MaxPlayers,
		OnlinePlayers: res.OnlinePlayers,
	})
}

func (server *Server) ProcessVeryLegacyStatusRequest(client *client, _ javaio.VeryLegacyStatusRequest) {
	res := server.handleStatusRequestV1()

	if (res.PreventResponse) {
		return
	}

	client.SendPacket(&javaio.VeryLegacyStatusResponse {
		Description: res.Description,
		MaxPlayers: res.MaxPlayers,
		OnlinePlayers: res.OnlinePlayers,
	})
}

func (server *Server) ProcessPing(client *client, ping javaio.Ping) {
	client.SendPacket(&javaio.Pong {
		Payload: ping.Payload,
	})
}

func (server *Server) ProcessLoginStart(client *client, data javaio.LoginStart) {
	println(data.ClientsideUsername)
	playerUuid := uuid.New()

	client.SendPacket(&javaio.LoginSuccess {
		Uuid: playerUuid,
		Username: data.ClientsideUsername,
	})

	client.ctx.State = javaio.StatePlay

	client.SendPacket(&javaio.JoinGame {
		EntityId: 0,
		Gamemode: javaio.GamemodeCreative,
		Hardcore: false,
		Dimension: javaio.DimensionOverworld,
		ViewDistance: 1,
		ReducedDebugInfo: false,
		EnableRespawnScreen: false,
	})

	client.SendPacket(&javaio.CompassPosition {
		Location: javaio.BlockPosition { X: 0, Y: 64, Z: 0 },
	})

	client.SendPacket(&javaio.PlayerPositionAndLook {
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
			client.SendPacket(&javaio.ChunkData {
				X: int32(x), Z: int32(z), IsNew: true,
				Sections: [][]uint32 { nil, blocksA[:], blocksB[:], blocksC[:] },
			})
		}
	}

	if server.handlePlayerJoin != nil {
		server.handlePlayerJoin(playerUuid.ID(), data.ClientsideUsername)
	}

	client.SendPacket(&javaio.PlayerInfoAdd {
		Players: []javaio.PlayerInfo {
			{ Uuid: uuid.New(), Username: "JohnDoe", Ping: 0 },
			{ Uuid: uuid.New(), Username: "CatsEyebrows", Ping: 5 },
			{ Uuid: uuid.New(), Username: "ElepantNostrel23", Ping: 500 },
		},
	})
}

func (server *Server) SpawnPlayer() {
	
}
