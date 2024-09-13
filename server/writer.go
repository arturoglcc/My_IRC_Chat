package main

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

func sendUserList(server *Server, client *Client) {
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

	sendJSONMessage(client, userListMessage, "Error al serializar la lista de usuarios",
		"Error al enviar la lista de usuarios al usuario")
}

func sendTextMessageToRecipient(server *Server, sender *Client, recipient *Client, text string) {
	// Crear el mensaje JSON
	message := map[string]string{
		"type":     "TEXT_FROM",
		"username": sender.ID,
		"text":     text,
	}

	sendJSONMessage(recipient, message, "Error al serializar el mensaje", "Error al enviar el mensaje")
}

func sendPublicTextToAll(server *Server, sender *Client, text string) {
	message := map[string]string{
		"type":     "PUBLIC_TEXT_FROM",
		"username": sender.ID,
		"text":     text,
	}

	sendMessageToAll(server, sender, message, false)
}

func sendRoomCreationSuccessResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "NEW_ROOM",
		"result":    "SUCCESS",
		"extra":     roomName,
	}

	sendJSONMessage(client, response, "Error al serializar la respuesta SUCCESS",
		"Error al enviar la respuesta SUCCESS")

}

func sendInvitationToUser(invitedClient *Client, inviterUsername string, roomName string) {
	invitation := map[string]string{
		"type":     "INVITATION",
		"username": inviterUsername,
		"roomname": roomName,
	}

	sendJSONMessage(invitedClient, invitation, "Error al serializar la respuesta INVITATION",
		"Error al enviar la respuesta INVITATION")
}

func sendJoinRoomSuccessResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "JOIN_ROOM",
		"result":    "SUCCESS",
		"extra":     roomName,
	}

	sendJSONMessage(client, response, "Error al serializar la respuesta SUCCESS",
		"Error al enviar la respuesta SUCCESS")

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

	sendJSONMessage(client, response, "Error al serializar la lista de usuarios",
		"Error al enviar la lista de usuarios")
}

func broadcastRoomTextMessage(server *Server, sender *Client, room *Room, text string) {
	message := map[string]string{
		"type":     "ROOM_TEXT_FROM",
		"roomname": room.Name,
		"username": sender.ID,
		"text":     text,
	}

	sendToRoomMembers(sender, room, message, "Error al serializar el mensaje ROOM_TEXT_FROM",
		"Error al enviar el mensaje ROOM_TEXT_FROM")

}

func notifyRoomMembersUserLeft(server *Server, client *Client, room *Room) {
	message := map[string]string{
		"type":     "LEFT_ROOM",
		"roomname": room.Name,
		"username": client.ID,
	}

	sendToRoomMembers(client, room, message, "Error al serializar el mensaje LEFT_ROOM",
		"Error al enviar el mensaje LEFT_ROOM")

}

func notifyAllUsersDisconnected(server *Server, client *Client) {
	message := map[string]string{
		"type":     "DISCONNECTED",
		"username": client.ID,
	}

	sendMessageToAll(server, client, message, true)
}
