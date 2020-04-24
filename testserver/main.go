package main

import "fmt"
import "net"
import "github.com/davidcallanan/gomcp/javasock"

func main() {
	server := javasock.NewServer()
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

		server.AddConnection(connection, func() {
			connection.Close()
		})
	}
}
