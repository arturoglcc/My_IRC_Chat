package main

import (
	"encoding/json"
	"log"
)

// sendMessageToAll envía un mensaje a todos los clientes conectados, excepto al remitente si es necesario.
func sendMessageToAll(server *Server, sender *Client, message interface{}, excludeSender bool) {
	// Serializar el mensaje a formato JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar el mensaje: %v", err)
		return
	}

	// Enviar el mensaje a todos los demás clientes conectados
	server.Mu.Lock()
	defer server.Mu.Unlock()

	for _, otherClient := range server.Clients {
		if excludeSender && otherClient.ID == sender.ID {
			continue // No enviar al remitente si excludeSender es true
		}
		_, err := otherClient.Conn.Write(jsonMessage)
		if err != nil {
			log.Printf("Error al enviar mensaje a %s: %v", otherClient.ID, err)
		}
	}
}

// sendTextMessageToRecipient envía un mensaje de texto de un usuario a otro.
func sendTextMessageToRecipient(server *Server, sender *Client, recipient *Client, text string) {
	// Crear el mensaje JSON
	message := map[string]string{
		"type":     "TEXT_FROM",
		"username": sender.ID,
		"text":     text,
	}

	// Serializar el mensaje a formato JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar el mensaje de texto: %v", err)
		return
	}

	// Enviar el mensaje al destinatario
	_, err = recipient.Conn.Write(jsonMessage)
	if err != nil {
		log.Printf("Error al enviar el mensaje a %s: %v", recipient.ID, err)
	}
}

// sendNoSuchUserResponse envía una respuesta al remitente cuando el destinatario no existe.
func sendNoSuchUserResponse(client *Client, operation string, username string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": operation,
		"result":    "NO_SUCH_USER",
		"extra":     username,
	}

	// Serializar el mensaje a formato JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta NO_SUCH_USER: %v", err)
		return
	}

	// Enviar la respuesta al remitente
	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta NO_SUCH_USER a %s: %v", client.ID, err)
	}
}

// sendNoSuchRoomResponse envía una respuesta indicando que la sala no existe, con una operación específica.
func sendNoSuchRoomResponse(client *Client, roomName string, operation string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": operation,
		"result":    "NO_SUCH_ROOM",
		"extra":     roomName,
	}

	// Serializar el mensaje a formato JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta NO_SUCH_ROOM: %v", err)
		return
	}

	// Enviar la respuesta al cliente
	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta NO_SUCH_ROOM a %s: %v", client.ID, err)
	}
}

// sendNotJoinedRoomResponse envía una respuesta indicando que el usuario no está en la sala.
func sendNotJoinedRoomResponse(client *Client, roomName string, operation string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": operation,
		"result":    "NOT_JOINED",
		"extra":     roomName,
	}

	// Serializar el mensaje a formato JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta NOT_JOINED: %v", err)
		return
	}

	// Enviar la respuesta al cliente
	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta NOT_JOINED a %s: %v", client.ID, err)
	}
}
