package main

import "fmt"
import "net"
import "github.com/davidcallanan/gomcp/javasock"

func main() {
	const maxPlayers = 20
	const onlinePlayers = 2
	const version = "1.14-1.15"

	var clients []uint32
	server := javasock.NewServer()
	listener, err := net.Listen("tcp4", "localhost:25565")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	server.OnStatusRequestV1(func() javasock.StatusResponseV1 {
		return javasock.StatusResponseV1 {
			Description: "Hello, World!",
			MaxPlayers: maxPlayers,
			OnlinePlayers: onlinePlayers,
		}
	})

	server.OnStatusRequestV2(func() javasock.StatusResponseV2 {
		return javasock.StatusResponseV2 {
			IsClientSupported: false,
			Version: version,
			Description: "§e§lHello, World!",
			MaxPlayers: maxPlayers,
			OnlinePlayers: onlinePlayers,
		}
	})

	server.OnStatusRequestV3(func() javasock.StatusResponseV3 {
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

	server.OnPlayerJoin(func(uuid uint32, _ string) {
		clients = append(clients, uuid)
	})

	fmt.Println("Test server is now listening...")

	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("Accepted a connection!")

		server.AddConnection(connection, connection, func() {
			connection.Close()
		})
	}
}
