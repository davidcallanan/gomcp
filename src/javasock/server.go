package javasock

import "io"
import "time"
import "bufio"
import "math/rand"
import "github.com/davidcallanan/gomcp/javaio"
import "github.com/google/uuid"

type client struct {
	state javaio.State
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
	handlePlayerJoin func(uuid uint32, clientsideUsername string)
}

func NewServer() Server {
	return Server {
	}
}

func (server *Server) OnPlayerJoin(callback func(uuid uint32, clientsideUsername string)) {
	server.handlePlayerJoin = callback
}

func (server *Server) AddConnection(input io.Reader, output io.Writer, closeCallback func()) {
	client := &client {
		state: javaio.StateDeterminingProtocol,
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
			if client.state != javaio.StatePlay {
				continue
			}

			javaio.EmitClientboundPacketUncompressed(&javaio.KeepAlive {
				Payload: now.Unix(),
			}, client.state, client.output)
		}
	}()
}

func (server *Server) handleReceive(client *client) {
	packet, err := javaio.ParseServerboundPacketUncompressed(client.input, client.state)

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

		// Default
	default:
		println("Unrecognized packet type")
	}
}

func (server *Server) ProcessProtocolDetermined(client *client, data javaio.ProtocolDetermined) {
	println(data.NextState)
	client.state = data.NextState
}

func (server *Server) ProcessHandshake(client *client, handshake javaio.Handshake) {
	client.state = handshake.NextState
}

func (server *Server) ProcessStatusRequest(client *client, _ javaio.StatusRequest) {
	javaio.EmitClientboundPacketUncompressed(&javaio.StatusResponse {
		Description: "§e§lHello, World!\n§rWelcome to this §kk§repic§kk§r server",
		VersionText: "Outdated Minecraft",
		VersionProtocol: 578,
		MaxPlayers: 20,
		OnlinePlayers: 2,
		PlayerSample: []javaio.StatusResponsePlayer {
			{ Name: "§aThis is", Uuid: "65bd239f-89f2-4cc7-ae8b-bb625525904e" },
			{ Name: "§cthe most", Uuid: "65bd239f-89f2-4cc7-ae8b-bb625525904e" },
			{ Name: "§8amazing thing", Uuid: "65bd239f-89f2-4cc7-ae8b-bb625525904e" },
			{ Name: "§9§lever!!!", Uuid: "65bd239f-89f2-4cc7-ae8b-bb625525904e" },
		},
	}, client.state, client.output)
}

func (server *Server) ProcessPing(client *client, ping javaio.Ping) {
	javaio.EmitClientboundPacketUncompressed(&javaio.Pong {
		Payload: ping.Payload,
	}, client.state, client.output)
}

func (server *Server) ProcessLoginStart(client *client, data javaio.LoginStart) {
	println(data.ClientsideUsername)
	playerUuid := uuid.New()

	javaio.EmitClientboundPacketUncompressed(&javaio.LoginSuccess {
		Uuid: playerUuid,
		Username: data.ClientsideUsername,
	}, client.state, client.output)

	client.state = javaio.StatePlay

	javaio.EmitClientboundPacketUncompressed(&javaio.JoinGame {
		Eid: 0,
		Gamemode: javaio.GamemodeCreative,
		Hardcore: false,
		Dimension: javaio.DimensionOverworld,
		ViewDistance: 1,
		ReducedDebugInfo: false,
		EnableRespawnScreen: false,
	}, client.state, client.output)

	javaio.EmitClientboundPacketUncompressed(&javaio.CompassPosition {
		Location: javaio.BlockPosition { X: 0, Y: 64, Z: 0 },
	}, client.state, client.output)

	javaio.EmitClientboundPacketUncompressed(&javaio.PlayerPositionAndLook {
		X: 0, Y: 64, Z: 0, Yaw: 0, Pitch: 0,
	}, client.state, client.output)

	var blocks [4096]uint32 // initialized to 0 (I hope that corresponds to stone)

	for i := range blocks {
		blocks[i] = uint32(rand.Intn(100))
		blocks[i] = 1
	}

	for x := -1; x <= 1; x++ {
		for z := -1; z <= 1; z++ {
			javaio.EmitClientboundPacketUncompressed(&javaio.ChunkData {
				X: int32(x), Z: int32(z), IsNew: true,
				Sections: [][]uint32 { nil, nil, nil, blocks[:] },
			}, client.state, client.output)
		}
	}

	if server.handlePlayerJoin != nil {
		server.handlePlayerJoin(playerUuid.ID(), data.ClientsideUsername)
	}
}

func (server *Server) SpawnPlayer() {
	
}
