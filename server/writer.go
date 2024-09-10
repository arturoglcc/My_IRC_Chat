package main

import (
	"encoding/json"
	"log"
)

func notifyNewUser(server *Server, newClient *Client) {
	notification := map[string]string{
		"type":     "NEW_USER",
		"username": newClient.ID,
	}
	sendMessageToAll(server, newClient, notification, true)
}

func notifyClientsOfStatusChange(server *Server, client *Client) {
	notification := map[string]string{
		"type":     "NEW_STATUS",
		"username": client.ID,
		"status":   client.Status,
	}

	sendMessageToAll(server, client, notification, true)
}

func sendUserList(server *Server, requestingClient *Client) {
	users := make(map[string]string)

	// Block to read client map in a safe way
	server.Mu.Lock()
	for username, client := range server.Clients {
		users[username] = client.Status
	}
	server.Mu.Unlock()

	userListMessage := map[string]interface{}{
		"type":  "USER_LIST",
		"users": users,
	}

	jsonUserList, err := json.Marshal(userListMessage)
	if err != nil {
		log.Printf("Error al serializar la lista de usuarios: %v", err)
		return
	}

	_, err = requestingClient.Conn.Write(jsonUserList)
	if err != nil {
		log.Printf("Error al enviar la lista de usuarios al cliente %s: %v", requestingClient.ID, err)
	}
}

func handleTextMessage(server *Server, sender *Client, msg map[string]interface{}) {
	recipientUsername, ok := msg["username"].(string)
	if !ok {
		sendInvalidMessageResponse(sender)
		return
	}

	text, ok := msg["text"].(string)
	if !ok {
		sendInvalidMessageResponse(sender)
		return
	}

	// Block to get clients map in a safe way
	server.Mu.Lock()
	recipient, exists := server.Clients[recipientUsername]
	server.Mu.Unlock()

	if exists {
		sendTextMessageToRecipient(server, sender, recipient, text)
	} else {
		sendNoSuchUserResponse(sender, recipientUsername, "TEXT")
	}
}

func sendPublicTextToAll(server *Server, sender *Client, text string) {
	message := map[string]string{
		"type":     "PUBLIC_TEXT_FROM",
		"username": sender.ID,
		"text":     text,
	}

	sendMessageToAll(server, sender, message, true)
}

func sendRoomAlreadyExistsResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "NEW_ROOM",
		"result":    "ROOM_ALREADY_EXISTS",
		"extra":     roomName,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta ROOM_ALREADY_EXISTS: %v", err)
		return
	}

	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta ROOM_ALREADY_EXISTS a %s: %v", client.ID, err)
	}
}

func createRoom(server *Server, client *Client, roomName string) {
	newRoom := &Room{
		Name:    roomName,
		Members: make(map[string]*Client),
		Invited: make(map[string]bool),
	}

	newRoom.Members[client.ID] = client

	server.RoomMu.Lock()
	server.Rooms[roomName] = newRoom
	server.RoomMu.Unlock()

	client.JoinedRooms[roomName] = newRoom

	sendRoomCreationSuccessResponse(client, roomName)
}

func sendRoomCreationSuccessResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "NEW_ROOM",
		"result":    "SUCCESS",
		"extra":     roomName,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta SUCCESS: %v", err)
		return
	}

	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta SUCCESS a %s: %v", client.ID, err)
	}
}

func inviteUsersToRoom(server *Server, client *Client, roomName string, invitedUsernames []interface{}) {
	server.RoomMu.Lock()
	room := server.Rooms[roomName]
	server.RoomMu.Unlock()

	for _, invitedUsername := range invitedUsernames {
		username, _ := invitedUsername.(string)

		if room.Members[username] != nil || room.Invited[username] {
			continue
		}

		server.Mu.Lock()
		invitedClient, exists := server.Clients[username]
		server.Mu.Unlock()

		if exists {
			sendInvitationToUser(invitedClient, client.ID, roomName)
			server.RoomMu.Lock()
			room.Invited[username] = true
			server.RoomMu.Unlock()
		}
	}
}

func sendInvitationToUser(invitedClient *Client, inviterUsername string, roomName string) {
	invitation := map[string]string{
		"type":     "INVITATION",
		"username": inviterUsername,
		"roomname": roomName,
	}

	jsonInvitation, err := json.Marshal(invitation)
	if err != nil {
		log.Printf("Error al serializar la invitación: %v", err)
		return
	}

	_, err = invitedClient.Conn.Write(jsonInvitation)
	if err != nil {
		log.Printf("Error al enviar la invitación a %s: %v", invitedClient.ID, err)
	}
}

func sendNotInvitedResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "JOIN_ROOM",
		"result":    "NOT_INVITED",
		"extra":     roomName,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta NOT_INVITED: %v", err)
		return
	}

	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta NOT_INVITED a %s: %v", client.ID, err)
	}
}

func sendJoinRoomSuccessResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "JOIN_ROOM",
		"result":    "SUCCESS",
		"extra":     roomName,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta SUCCESS: %v", err)
		return
	}

	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta SUCCESS a %s: %v", client.ID, err)
	}
}

func notifyRoomMembersUserJoined(server *Server, room *Room, client *Client) {
	message := map[string]string{
		"type":     "JOINED_ROOM",
		"roomname": room.Name,
		"username": client.ID,
	}
	sendMessageToAll(server, client, message, true)
}

func sendRoomUserList(server *Server, client *Client, room *Room) {
	users := make(map[string]string)

	for username, member := range room.Members {
		users[username] = member.Status
	}

	response := map[string]interface{}{
		"type":     "ROOM_USER_LIST",
		"roomname": room.Name,
		"users":    users,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la lista de usuarios: %v", err)
		return
	}

	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la lista de usuarios a %s: %v", client.ID, err)
	}
}

func broadcastRoomTextMessage(server *Server, sender *Client, room *Room, text string) {
	message := map[string]string{
		"type":     "ROOM_TEXT_FROM",
		"roomname": room.Name,
		"username": sender.ID,
		"text":     text,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar el mensaje ROOM_TEXT_FROM: %v", err)
		return
	}

	for _, member := range room.Members {
		if member.ID != sender.ID { // No enviar el mensaje al remitente
			_, err := member.Conn.Write(jsonMessage)
			if err != nil {
				log.Printf("Error al enviar el mensaje ROOM_TEXT_FROM a %s: %v", member.ID, err)
			}
		}
	}
}

func notifyRoomMembersUserLeft(server *Server, client *Client, room *Room) {
	message := map[string]string{
		"type":     "LEFT_ROOM",
		"roomname": room.Name,
		"username": client.ID,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar el mensaje LEFT_ROOM: %v", err)
		return
	}

	for _, member := range room.Members {
		if member.ID != client.ID {
			_, err := member.Conn.Write(jsonMessage)
			if err != nil {
				log.Printf("Error al enviar el mensaje LEFT_ROOM a %s: %v", member.ID, err)
			}
		}
	}
}

func notifyAllUsersDisconnected(server *Server, client *Client) {
	message := map[string]string{
		"type":     "DISCONNECTED",
		"username": client.ID,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar el mensaje DISCONNECTED: %v", err)
		return
	}

	// Enviar el mensaje a todos los usuarios conectados
	server.Mu.Lock()
	for _, otherClient := range server.Clients {
		if otherClient.ID != client.ID {
			_, err := otherClient.Conn.Write(jsonMessage)
			if err != nil {
				log.Printf("Error al enviar el mensaje DISCONNECTED a %s: %v", otherClient.ID, err)
			}
		}
	}
	server.Mu.Unlock()
}
