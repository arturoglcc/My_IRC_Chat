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
	Status      string
}

type Room struct {
	Name    string
	Members map[string]*Client
	Invited map[string]bool
}

type Server struct {
	Address string
	Clients map[string]*Client
	Rooms   map[string]*Room
	Mu      sync.Mutex
	RoomMu  sync.Mutex
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
		Name:    "general",
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
		// Crear la respuesta en caso de que el tipo no sea "IDENTIFY"
		response := map[string]string{
			"type":      "RESPONSE",
			"operation": "INVALID",
			"result":    "NOT_IDENTIFIED",
		}

		// Serializar la respuesta JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error al serializar la respuesta JSON: %v", err)
			return
		}

		// Enviar la respuesta al cliente
		conn.Write(jsonResponse)

		// Desconectar al cliente
		log.Printf("Desconectando al cliente por mensaje no válido.")
		return
	}

	// Verificar la longitud del ID (debe ser a lo más 8 caracteres)
	if len(identifyMsg.Username) > 8 {
		response := map[string]string{
			"type":      "RESPONSE",
			"operation": "IDENTIFY",
			"result":    "INVALID_USERNAME_LENGTH",
			"extra":     identifyMsg.Username,
		}

		// Serializar la respuesta JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error al serializar la respuesta JSON: %v", err)
			return
		}

		// Enviar la respuesta al cliente
		conn.Write(jsonResponse)
		return
	}

	// Comprobar si el nombre de usuario es único
	s.Mu.Lock()
	_, exists := s.Clients[identifyMsg.Username]
	s.Mu.Unlock()

	if exists {
		// El nombre de usuario ya está en uso, enviar respuesta en formato JSON
		response := map[string]string{
			"type":      "RESPONSE",
			"operation": "IDENTIFY",
			"result":    "USER_ALREADY_EXISTS",
			"extra":     identifyMsg.Username,
		}

		// Serializar la respuesta JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error al serializar la respuesta JSON: %v", err)
			return
		}

		// Enviar la respuesta al cliente
		conn.Write(jsonResponse)
		return
	}

	// Crear el cliente
	client := &Client{
		ID:          identifyMsg.Username,
		Conn:        conn,
		JoinedRooms: make(map[string]*Room),
		Status:      "ACTIVE",
	}

	// Unir al cliente al cuarto general
	s.Mu.Lock()
	generalRoom := s.Rooms["general"]
	generalRoom.Members[identifyMsg.Username] = client
	client.JoinedRooms["general"] = generalRoom
	s.Clients[identifyMsg.Username] = client
	s.Mu.Unlock()

	conn.Write([]byte("Te has unido al cuarto general.\n"))

	// Notificar a los demás clientes que un nuevo usuario se ha conectado
	notifyNewUser(s, client)

	go listenToClientMessages(s, client)
}
