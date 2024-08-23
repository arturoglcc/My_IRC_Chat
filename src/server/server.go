package main

import (
	"fmt"
	"log"
	"net"
)

// StartServer initializes the server on the given port
func StartServer(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("could not start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server is listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Handling new connection from %s", conn.RemoteAddr().String())
}
