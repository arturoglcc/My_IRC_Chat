package main

import (
	"encoding/json"
	"io"
	"log"
	"strings"
)

// sendInvalidMessageResponse Send error message to the client and then disconnects them
func sendInvalidMessageResponse(client *Client) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "INVALID",
		"result":    "INVALID",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta JSON: %v", err)
		return
	}

	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta al cliente %s: %v", client.ID, err)
	}

	client.Conn.Close()
	log.Printf("Conexión con el cliente %s cerrada por mensaje inválido.", client.ID)
}
func listenToClientMessages(server *Server, client *Client) {
	buf := make([]byte, 1024)
	var incompleteMessage string

	for {
		n, err := client.Conn.Read(buf)
		if err != nil {
			// Verificar si el error es por conexión cerrada
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("Cliente %s desconectado.", client.ID)
				handleDisconnectMessage(server, client)
				break
			}
			log.Printf("Error leyendo mensaje del cliente %s: %v", client.ID, err)
			break
		}

		// Añadir al buffer de mensajes incompletos
		incompleteMessage += string(buf[:n])

		// Verificar si el mensaje está completo (por ejemplo, si es un JSON válido)
		var msg map[string]interface{}
		err = json.Unmarshal([]byte(incompleteMessage), &msg)
		if err != nil {
			// Si no es JSON válido, espera por más datos y no desconectes al cliente
			continue
		}

		// Procesar el mensaje completo y reiniciar el buffer
		processClientMessage(server, client, msg)
		incompleteMessage = ""
	}
}

func processClientMessage(server *Server, client *Client, msg map[string]interface{}) {
	// verify the message has the "type" field
	msgType, ok := msg["type"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Manage each type of message
	switch msgType {
	case "STATUS":
		handleStatusMessage(server, client, msg)
	case "USERS":
		sendUserList(server, client)
	case "TEXT":
		handleTextMessage(server, client, msg)
	case "PUBLIC_TEXT":
		handlePublicTextMessage(server, client, msg)
	case "NEW_ROOM":
		handleNewRoomMessage(server, client, msg)
	case "INVITE":
		handleInviteMessage(server, client, msg)
	case "JOIN_ROOM":
		handleJoinRoomMessage(server, client, msg)
	case "ROOM_USERS":
		handleRoomUsersMessage(server, client, msg)
	case "ROOM_TEXT":
		handleRoomTextMessage(server, client, msg)
	case "LEAVE_ROOM":
		handleLeaveRoomMessage(server, client, msg)
	case "DISCONNECT":
		handleDisconnectMessage(server, client)
	default:
		// if the message is not a valid message, disconnect the client
		sendInvalidMessageResponse(client)
	}
}
