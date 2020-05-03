package main

import "fmt"
import "net"
import "github.com/davidcallanan/gomcp/javaserver"
import "github.com/google/uuid"

func main() {
	const maxPlayers = 20
	const onlinePlayers = 2
	const version = "1.14-1.15"

	listener, err := net.Listen("tcp4", "localhost:25565")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Test server is now listening...")

	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("Accepted a connection!")

		var conn *javaserver.Connection
		var guid uuid.UUID

		conn = javaserver.NewConnection(connection, func() {
			connection.Close()
		}, javaserver.EventHandlers {
			OnStatusRequestV1: func() javaserver.StatusResponseV1 {
				return javaserver.StatusResponseV1 {
					Description: "Hello, World!",
					MaxPlayers: maxPlayers,
					OnlinePlayers: onlinePlayers,
				}
			},
		
			OnStatusRequestV2: func() javaserver.StatusResponseV2 {
				return javaserver.StatusResponseV2 {
					IsClientSupported: false,
					Version: version,
					Description: "§e§lHello, World!",
					MaxPlayers: maxPlayers,
					OnlinePlayers: onlinePlayers,
				}
			},
		
			OnStatusRequestV3: func() javaserver.StatusResponseV3 {
				return javaserver.StatusResponseV3 {
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
			},
		
			OnPlayerJoinRequest: func(data javaserver.PlayerJoinRequest) javaserver.PlayerJoinResponse {
				fmt.Printf("Player %s has requested to join the game.\n", data.ClientsideUsername)
				guid = uuid.New()
				return javaserver.PlayerJoinResponse {
					Uuid: guid,
				}
			},
		
			OnPlayerJoin: func() {
				fmt.Println("Player of whom I forget their username has joined the game.")

				conn.AddPlayerInfo([]javaserver.PlayerInfoToAdd {
					{ Uuid: guid, Username: "JohnDoe", Ping: 0 },
					{ Uuid: uuid.New(), Username: "CatsEyebrows", Ping: 5 },
					{ Uuid: uuid.New(), Username: "ElepantNostrel23", Ping: 500 },
				})

				conn.SpawnPlayer(javaserver.PlayerToSpawn {
					EntityId: 123,
					Uuid: guid,
					X: 0,
					Y: 70,
					Z: 0,
					Yaw: 0,
					Pitch: 0,
				})
			},
		})
	}
}
