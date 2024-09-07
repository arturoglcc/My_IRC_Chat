// server.go
package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Client struct {
	ID          string
	Conn        net.Conn
	OwnedRooms  map[string]*Room
	JoinedRooms map[string]*Room
}

type Room struct {
	ID      string
	Owner   *Client
	Members map[string]*Client
}

// Server struct represents the server with necessary fields to manage connections and handle client communication.
type Server struct {
	Address string
	Clients map[string]*Client
	Rooms   map[string]*Room
	Mu      sync.Mutex
}

// NewServer creates a new server instance with the specified address.
func NewServer(address string) *Server {
	server := &Server{
		Address: address,
		Clients: make(map[string]*Client),
		Rooms:   make(map[string]*Room),
	}

	generalRoom := &Room{
		ID:      "general",
		Owner:   nil,
		Members: make(map[string]*Client),
	}

	server.Rooms["general"] = generalRoom

	return server
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

	buf := make([]byte, 1024)

	var clientID string
	for {
		conn.Write([]byte("Ingrese un ID único: "))
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error leyendo ID del cliente: %v", err)
			return
		}
		clientID = string(buf[:n])

		s.Mu.Lock()
		_, exists := s.Clients[clientID]
		s.Mu.Unlock()

		if exists {
			conn.Write([]byte("El ID ya está en uso. Intente con otro ID.\n"))
		} else {
			break
		}
	}

	// Crear el cliente
	client := &Client{
		ID:          clientID,
		Conn:        conn,
		OwnedRooms:  make(map[string]*Room),
		JoinedRooms: make(map[string]*Room),
	}

	// Unir al cliente al cuarto general
	s.Mu.Lock()
	generalRoom := s.Rooms["general"]
	generalRoom.Members[clientID] = client
	client.JoinedRooms["general"] = generalRoom
	s.Clients[clientID] = client
	s.Mu.Unlock()

	conn.Write([]byte("Te has unido al cuarto general.\n"))
	go listenToClientMessages(s, client)
	go writeToClient(s, client)
}
