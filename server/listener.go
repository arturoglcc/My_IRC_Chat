package main

import (
	"encoding/json"
	"log"
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

	for {
		n, err := client.Conn.Read(buf)
		if err != nil {
			log.Printf("Error leyendo mensaje del cliente %s: %v", client.ID, err)
			break
		}
		message := string(buf[:n])

		// Verify message is a valid json
		var msg map[string]interface{}
		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			// The message is not a valid json, disconnects the client and send a message
			sendInvalidMessageResponse(client)
			break
		}

		processClientMessage(server, client, msg)
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
