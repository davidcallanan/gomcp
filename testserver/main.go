package main

import "fmt"
import "net"
import "github.com/davidcallanan/gomcp/javasock"
import "github.com/google/uuid"

func main() {
	const maxPlayers = 20
	const onlinePlayers = 2
	const version = "1.14-1.15"

	server := javasock.NewServer()
	listener, err := net.Listen("tcp4", "localhost:25565")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	server.OnStatusRequestV1(func(_id int) javasock.StatusResponseV1 {
		return javasock.StatusResponseV1 {
			Description: "Hello, World!",
			MaxPlayers: maxPlayers,
			OnlinePlayers: onlinePlayers,
		}
	})

	server.OnStatusRequestV2(func(_id int) javasock.StatusResponseV2 {
		return javasock.StatusResponseV2 {
			IsClientSupported: false,
			Version: version,
			Description: "§e§lHello, World!",
			MaxPlayers: maxPlayers,
			OnlinePlayers: onlinePlayers,
		}
	})

	server.OnStatusRequestV3(func(_id int) javasock.StatusResponseV3 {
		return javasock.StatusResponseV3 {
			IsClientSupported: true,
			Version: version,
			Description: "§e§lHello, World!\n§r§aWelcome to this amazing server",
			MaxPlayers: maxPlayers,
			OnlinePlayers: onlinePlayers,
			PlayerSample: []string {
				"§aThis is",
				"§cthe most",
				"§8amazing thing",
				"§9§lever!!!",
			},
		}
	})

	server.OnPlayerJoinRequest(func(data javasock.PlayerJoinRequest) javasock.PlayerJoinResponse {
		fmt.Printf("Player %s has requested to join the game.\n", data.ClientsideUsername)
		return javasock.PlayerJoinResponse {
			Uuid: uuid.New(),
		}
	})

	server.OnPlayerJoin(func(id int) {
		fmt.Printf("Player with id %d has joined the game.\n", id)
	})

	fmt.Println("Test server is now listening...")

	nextId := 0

	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		id := nextId
		nextId++
		fmt.Println("Accepted a connection!")

		server.AddConnection(id, connection, connection, func() {
			connection.Close()
		})
	}
}
