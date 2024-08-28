// server.go
package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

// Server struct represents the server with necessary fields to manage connections and handle client communication.
type Server struct {
	Address string
	Clients map[net.Conn]bool
	Mu      sync.Mutex
	Rooms   map[string][]net.Conn
	RoomMu  sync.Mutex
}

// NewServer creates a new server instance with the specified address.
func NewServer(address string) *Server {
	return &Server{
		Address: address,
		Clients: make(map[net.Conn]bool),
		Rooms:   make(map[string][]net.Conn),
	}
}

// Start starts the server to listen for incoming client connections.
func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", s.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection handles a client connection in a separate goroutine.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Add the connection to the clients map
	s.Mu.Lock()
	s.Clients[conn] = true
	s.Mu.Unlock()

	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	// Handle client communication here (e.g., reading from/writing to conn)
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			break
		}
		message := string(buf[:n])
		fmt.Printf("Received message: %s\n", message)

		// Echo back the message
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			break
		}
	}

	// Remove the connection from the clients map when done
	s.Mu.Lock()
	delete(s.Clients, conn)
	s.Mu.Unlock()

	fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr().String())
}
