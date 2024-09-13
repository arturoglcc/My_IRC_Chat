// server.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

// IdentifyMessage Json structure for identification
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

// Handle signals to shut down the server gracefully
func (s *Server) handleSignals() {
	// Create a channel to listen for interrupt signals (Ctrl+C, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	sig := <-sigChan
	fmt.Printf("\nReceived signal: %v, shutting down server...\n", sig)

	// Disconnect all clients
	s.Mu.Lock()
	for _, client := range s.Clients { // Accede a los valores, no a las claves
		client.Conn.Close() // Close the connection properly
	}
	s.Mu.Unlock()

	// Optionally, wait for goroutines to finish here (if needed)

	os.Exit(0) // Exit the program
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
			time.Sleep(1 * time.Second)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	//defer conn.Close()

	buf := make([]byte, 1024)

	// Reads the first message of the client (it's supposed to be the identify Json)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error leyendo del cliente: %v", err)
		return
	}

	//Deserialize JSON
	var identifyMsg IdentifyMessage
	err = json.Unmarshal(buf[:n], &identifyMsg)
	if err != nil {
		errorResponse := map[string]string{
			"Type":      "RESPONSE",
			"Operation": "INVALID",
			"Result":    "INVALID",
		}

		jsonResponse, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println("Error al generar el JSON de respuesta:", err)
			return
		}

		conn.Write(jsonResponse)
		return
	}

	// Verify type "IDENTIFY"
	if identifyMsg.Type != "IDENTIFY" {
		// response in case type is not "IDENTIFY"
		response := map[string]string{
			"type":      "RESPONSE",
			"operation": "INVALID",
			"result":    "NOT_IDENTIFIED",
		}

		// Serialice Json response
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error al serializar la respuesta JSON: %v", err)
			return
		}

		// Send the response to the client
		conn.Write(jsonResponse)

		// Disconnect client
		log.Printf("Desconectando al cliente por mensaje no válido.")
		return
	}

	// Verify length of username (It must be up to 8 characters)
	if len(identifyMsg.Username) > 8 {
		response := map[string]string{
			"type":      "RESPONSE",
			"operation": "IDENTIFY",
			"result":    "INVALID_USERNAME_LENGTH",
			"extra":     identifyMsg.Username,
		}

		// Serialice Json response
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error al serializar la respuesta JSON: %v", err)
			return
		}

		// Send the response to the client
		conn.Write(jsonResponse)
		return
	}

	// Check if username is unique
	s.Mu.Lock()
	_, exists := s.Clients[identifyMsg.Username]
	s.Mu.Unlock()

	if exists {
		// If username is not unique, send response in Json format
		response := map[string]string{
			"type":      "RESPONSE",
			"operation": "IDENTIFY",
			"result":    "USER_ALREADY_EXISTS",
			"extra":     identifyMsg.Username,
		}

		// Serialice Json response
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error al serializar la respuesta JSON: %v", err)
			return
		}

		// Send the response to the client
		conn.Write(jsonResponse)
		return
	}

	// create client
	client := &Client{
		ID:          identifyMsg.Username,
		Conn:        conn,
		JoinedRooms: make(map[string]*Room),
		Status:      "ACTIVE",
	}

	// Join client to general room
	s.Mu.Lock()
	generalRoom := s.Rooms["general"]
	generalRoom.Members[identifyMsg.Username] = client
	client.JoinedRooms["general"] = generalRoom
	s.Clients[identifyMsg.Username] = client
	s.Mu.Unlock()

	// Crear el mensaje de respuesta
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "IDENTIFY",
		"result":    "SUCCESS",
		"extra":     client.ID,
	}

	// Serializar el mensaje a formato JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar el mensaje de IDENTIFY: %v", err)
		return
	}

	// Enviar el mensaje al cliente
	_, err = conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar el mensaje de IDENTIFY a %s: %v", client.ID, err)
		return
	}

	fmt.Printf("Usuario %s identificado exitosamente.\n", client.ID)

	notifyNewUser(s, client)

	go listenToClientMessages(s, client)
}

func terminateServer() {
	log.Println("El cuarto general está vacío. Terminando el servidor...")
	os.Exit(0) // Cerrar el servidor
}
