package javasock

import "io"
import "bufio"
import "github.com/davidcallanan/gomcp/javaio"

type Server struct {
	state int
}

func NewServer() Server {
	return Server {
		state: javaio.StateHandshaking,
	}
}

func (s *Server) AddConnection(ior io.Reader, close func()) {
	r := bufio.NewReader(ior)

	go func() {
		for {
			packet, err := javaio.ParseServerboundPacketUncompressed(r, s.state)

			if err != nil {
				switch err.(type) {
				case *javaio.UnsupportedPayloadError:
					continue
				case *javaio.MalformedPacketError:
					close()
					return
				default:
					panic(err)
				}
			}

			switch p := packet.(type) {
			case javaio.Handshake:
				s.ProcessHandshake(p)
			default:
				panic("Unrecognized packet type")
			}
		}
	}()
}

func (s *Server) ProcessHandshake(handshake javaio.Handshake) {
	println("GOT A HANDSHAKE YAY!")
}
