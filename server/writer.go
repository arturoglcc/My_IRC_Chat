package main

import (
	"encoding/json"
	"log"
)

func notifyNewUser(server *Server, newClient *Client) {
	// Crear el mensaje JSON
	notification := map[string]string{
		"type":     "NEW_USER",
		"username": newClient.ID,
	}
	sendMessageToAll(server, newClient, notification, true)
}

// notifyClientsOfStatusChange envía un mensaje a todos los demás clientes informando del nuevo estado.
func notifyClientsOfStatusChange(server *Server, client *Client) {
	// Crear el mensaje JSON para notificar a los otros clientes
	notification := map[string]string{
		"type":     "NEW_STATUS",
		"username": client.ID,
		"status":   client.Status,
	}

	sendMessageToAll(server, client, notification, true)
}

// sendUserList envía la lista de usuarios conectados y sus estados al cliente solicitante.
func sendUserList(server *Server, requestingClient *Client) {
	// Crear el mapa de usuarios con sus estados
	users := make(map[string]string)

	// Bloqueo para leer el mapa de clientes de forma segura
	server.Mu.Lock()
	for username, client := range server.Clients {
		users[username] = client.Status
	}
	server.Mu.Unlock()

	// Crear el mensaje JSON
	userListMessage := map[string]interface{}{
		"type":  "USER_LIST",
		"users": users,
	}

	// Serializar el mensaje a formato JSON
	jsonUserList, err := json.Marshal(userListMessage)
	if err != nil {
		log.Printf("Error al serializar la lista de usuarios: %v", err)
		return
	}

	// Enviar el mensaje al cliente que lo solicitó
	_, err = requestingClient.Conn.Write(jsonUserList)
	if err != nil {
		log.Printf("Error al enviar la lista de usuarios al cliente %s: %v", requestingClient.ID, err)
	}
}

func handleTextMessage(server *Server, sender *Client, msg map[string]interface{}) {
	// Verificar que el campo "username" y "text" estén presentes
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

	// Bloqueo para acceder al mapa de clientes
	server.Mu.Lock()
	recipient, exists := server.Clients[recipientUsername]
	server.Mu.Unlock()

	if exists {
		// El destinatario existe, enviarle el mensaje
		sendTextMessageToRecipient(server, sender, recipient, text)
	} else {
		// El destinatario no existe, responder con un mensaje de error al remitente
		sendNoSuchUserResponse(sender, recipientUsername, "TEXT")
	}
}

// sendPublicTextToAll envía un mensaje público a todos los usuarios excepto al remitente.
func sendPublicTextToAll(server *Server, sender *Client, text string) {
	// Crear el mensaje JSON
	message := map[string]string{
		"type":     "PUBLIC_TEXT_FROM",
		"username": sender.ID,
		"text":     text,
	}

	// Usar la función auxiliar para enviar el mensaje a todos
	sendMessageToAll(server, sender, message, true)
}

// sendRoomAlreadyExistsResponse envía una respuesta indicando que la sala ya existe o es inválida.
func sendRoomAlreadyExistsResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "NEW_ROOM",
		"result":    "ROOM_ALREADY_EXISTS",
		"extra":     roomName,
	}

	// Serializar el mensaje a formato JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta ROOM_ALREADY_EXISTS: %v", err)
		return
	}

	// Enviar la respuesta al cliente
	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta ROOM_ALREADY_EXISTS a %s: %v", client.ID, err)
	}
}

// createRoom crea una nueva sala y agrega al cliente como el primer miembro.
func createRoom(server *Server, client *Client, roomName string) {
	// Crear la nueva sala
	newRoom := &Room{
		Name:    roomName,
		Members: make(map[string]*Client),
		Invited: make(map[string]bool),
	}

	// Agregar el cliente como miembro de la sala
	newRoom.Members[client.ID] = client

	// Actualizar el mapa de salas en el servidor
	server.RoomMu.Lock()
	server.Rooms[roomName] = newRoom
	server.RoomMu.Unlock()

	// Agregar la sala a la lista de salas unidas del cliente
	client.JoinedRooms[roomName] = newRoom

	// Enviar la respuesta de éxito al cliente
	sendRoomCreationSuccessResponse(client, roomName)
}

// sendRoomCreationSuccessResponse envía una respuesta indicando que la sala se creó exitosamente.
func sendRoomCreationSuccessResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "NEW_ROOM",
		"result":    "SUCCESS",
		"extra":     roomName,
	}

	// Serializar el mensaje a formato JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta SUCCESS: %v", err)
		return
	}

	// Enviar la respuesta al cliente
	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta SUCCESS a %s: %v", client.ID, err)
	}
}

// inviteUsersToRoom envía una invitación a todos los usuarios en la lista de invitados.
func inviteUsersToRoom(server *Server, client *Client, roomName string, invitedUsernames []interface{}) {
	server.RoomMu.Lock()
	room := server.Rooms[roomName]
	server.RoomMu.Unlock()

	for _, invitedUsername := range invitedUsernames {
		username, _ := invitedUsername.(string)

		// Verificar si el usuario ya está en la sala o si ya ha sido invitado
		if room.Members[username] != nil || room.Invited[username] {
			continue // No enviar la invitación si el usuario ya está en la sala o ya fue invitado
		}

		server.Mu.Lock()
		invitedClient, exists := server.Clients[username]
		server.Mu.Unlock()

		if exists {
			sendInvitationToUser(invitedClient, client.ID, roomName)
			// Marcar al usuario como invitado
			server.RoomMu.Lock()
			room.Invited[username] = true
			server.RoomMu.Unlock()
		}
	}
}

// sendInvitationToUser envía una invitación a un usuario.
func sendInvitationToUser(invitedClient *Client, inviterUsername string, roomName string) {
	invitation := map[string]string{
		"type":     "INVITATION",
		"username": inviterUsername,
		"roomname": roomName,
	}

	// Serializar el mensaje a formato JSON
	jsonInvitation, err := json.Marshal(invitation)
	if err != nil {
		log.Printf("Error al serializar la invitación: %v", err)
		return
	}

	// Enviar la invitación al cliente invitado
	_, err = invitedClient.Conn.Write(jsonInvitation)
	if err != nil {
		log.Printf("Error al enviar la invitación a %s: %v", invitedClient.ID, err)
	}
}

// sendNotInvitedResponse envía una respuesta indicando que el usuario no fue invitado a la sala.
func sendNotInvitedResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "JOIN_ROOM",
		"result":    "NOT_INVITED",
		"extra":     roomName,
	}

	// Serializar el mensaje a formato JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta NOT_INVITED: %v", err)
		return
	}

	// Enviar la respuesta al cliente
	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta NOT_INVITED a %s: %v", client.ID, err)
	}
}

// sendJoinRoomSuccessResponse envía una respuesta indicando que el usuario se ha unido exitosamente a la sala.
func sendJoinRoomSuccessResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "JOIN_ROOM",
		"result":    "SUCCESS",
		"extra":     roomName,
	}

	// Serializar el mensaje a formato JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta SUCCESS: %v", err)
		return
	}

	// Enviar la respuesta al cliente
	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta SUCCESS a %s: %v", client.ID, err)
	}
}

// notifyRoomMembersUserJoined notifica a los demás miembros de la sala que un nuevo usuario se ha unido.
func notifyRoomMembersUserJoined(server *Server, room *Room, client *Client) {
	message := map[string]string{
		"type":     "JOINED_ROOM",
		"roomname": room.Name,
		"username": client.ID,
	}

	sendMessageToAll(server, client, message, true)
}

// sendRoomUserList envía la lista de usuarios en la sala al cliente.
func sendRoomUserList(server *Server, client *Client, room *Room) {
	users := make(map[string]string)

	// Recolectar la lista de usuarios y sus estados
	for username, member := range room.Members {
		users[username] = member.Status
	}

	// Crear el mensaje de lista de usuarios
	response := map[string]interface{}{
		"type":     "ROOM_USER_LIST",
		"roomname": room.Name,
		"users":    users,
	}

	// Serializar el mensaje a formato JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la lista de usuarios: %v", err)
		return
	}

	// Enviar el mensaje al cliente
	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la lista de usuarios a %s: %v", client.ID, err)
	}
}

// broadcastRoomTextMessage envía un mensaje de texto a todos los usuarios de la sala excepto al remitente.
func broadcastRoomTextMessage(server *Server, sender *Client, room *Room, text string) {
	message := map[string]string{
		"type":     "ROOM_TEXT_FROM",
		"roomname": room.Name,
		"username": sender.ID,
		"text":     text,
	}

	// Serializar el mensaje a formato JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar el mensaje ROOM_TEXT_FROM: %v", err)
		return
	}

	// Enviar el mensaje a todos los miembros de la sala
	for _, member := range room.Members {
		if member.ID != sender.ID { // No enviar el mensaje al remitente
			_, err := member.Conn.Write(jsonMessage)
			if err != nil {
				log.Printf("Error al enviar el mensaje ROOM_TEXT_FROM a %s: %v", member.ID, err)
			}
		}
	}
}

// notifyRoomMembersUserLeft notifica a los demás miembros de la sala que un usuario ha salido.
func notifyRoomMembersUserLeft(server *Server, client *Client, room *Room) {
	message := map[string]string{
		"type":     "LEFT_ROOM",
		"roomname": room.Name,
		"username": client.ID,
	}

	// Serializar el mensaje a formato JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar el mensaje LEFT_ROOM: %v", err)
		return
	}

	// Enviar el mensaje a todos los miembros de la sala
	for _, member := range room.Members {
		if member.ID != client.ID {
			_, err := member.Conn.Write(jsonMessage)
			if err != nil {
				log.Printf("Error al enviar el mensaje LEFT_ROOM a %s: %v", member.ID, err)
			}
		}
	}
}

// notifyAllUsersDisconnected notifica a todos los usuarios que un usuario se ha desconectado.
func notifyAllUsersDisconnected(server *Server, client *Client) {
	message := map[string]string{
		"type":     "DISCONNECTED",
		"username": client.ID,
	}

	// Serializar el mensaje a formato JSON
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
