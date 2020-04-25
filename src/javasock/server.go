package javasock

import "io"
import "bufio"
import "github.com/davidcallanan/gomcp/javaio"
import "github.com/google/uuid"

type client struct {
	state int
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
}

func NewServer() Server {
	return Server {
	}
}

func (server *Server) AddConnection(input io.Reader, output io.Writer, closeCallback func()) {
	client := &client {
		state: javaio.StateHandshaking,
		input: bufio.NewReader(input),
		output: bufio.NewWriter(output),
		closeCallback: closeCallback,
	}

	go func() {
		for !client.isClosed {
			server.handleReceive(client)
		}
	}()
}

func (server *Server) handleReceive(client *client) {
	packet, err := javaio.ParseServerboundPacketUncompressed(client.input, client.state)

	if err != nil {
		switch err.(type) {
		case *javaio.UnsupportedPayloadError:
			return
		case *javaio.MalformedPacketError:
			client.close()
			return
		default:
			panic(err)
		}
	}

	switch packet := packet.(type) {
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
		panic("Unrecognized packet type")
	}
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
	javaio.EmitClientboundPacketUncompressed(&javaio.LoginSuccess {
		Uuid: uuid.New(),
		Username: data.ClientsideUsername,
	}, client.state, client.output)
}
