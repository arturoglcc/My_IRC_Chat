// server.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type Client struct {
	ID          string
	Conn        net.Conn
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

// Estructura del mensaje JSON para identificación
type IdentifyMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
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

// Manejar conexiones
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	// Leer el mensaje inicial del cliente (se espera que sea el JSON de identificación)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error leyendo del cliente: %v", err)
		return
	}

	// Deserializar el JSON
	var identifyMsg IdentifyMessage
	err = json.Unmarshal(buf[:n], &identifyMsg)
	if err != nil {
		conn.Write([]byte("Error: El mensaje no es un JSON válido.\n"))
		return
	}

	// Verificar que el tipo sea "IDENTIFY"
	if identifyMsg.Type != "IDENTIFY" {
		conn.Write([]byte("Error: Se esperaba un mensaje de tipo IDENTIFY.\n"))
		return
	}

	// Comprobar si el nombre de usuario es único
	s.Mu.Lock()
	_, exists := s.Clients[identifyMsg.Username]
	s.Mu.Unlock()

	if exists {
		conn.Write([]byte("Error: El nombre de usuario ya está en uso.\n"))
		return
	}
	// Crear el cliente
	client := &Client{
		ID:          identifyMsg.Username,
		Conn:        conn,
		JoinedRooms: make(map[string]*Room),
	}

	// Unir al cliente al cuarto general
	s.Mu.Lock()
	generalRoom := s.Rooms["general"]
	generalRoom.Members[identifyMsg.Username] = client
	client.JoinedRooms["general"] = generalRoom
	s.Clients[identifyMsg.Username] = client
	s.Mu.Unlock()

	conn.Write([]byte("Te has unido al cuarto general.\n"))

	go listenToClientMessages(s, client)
	go writeToClient(s, client)
}
