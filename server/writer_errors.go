package main

func sendRoomAlreadyExistsResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "NEW_ROOM",
		"result":    "ROOM_ALREADY_EXISTS",
		"extra":     roomName,
	}

	sendJSONMessage(client, response, "Error al serializar la respuesta ROOM_ALREADY_EXISTS",
		"Error al enviar la respuesta ROOM_ALREADY_EXISTS")

}

func sendNotInvitedResponse(client *Client, roomName string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "JOIN_ROOM",
		"result":    "NOT_INVITED",
		"extra":     roomName,
	}

	sendJSONMessage(client, response, "Error al serializar la respuesta NOT_INVITED",
		"Error al enviar la respuesta NOT_INVITED")
}

func sendNoSuchUserResponse(client *Client, operation string, username string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": operation,
		"result":    "NO_SUCH_USER",
		"extra":     username,
	}

	sendJSONMessage(client, response, "Error al serializar el mensaje NO_SUCH_USER",
		"Error al enviar el mensaje NO_SUCH_USER")
}

func sendNoSuchRoomResponse(client *Client, roomName string, operation string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": operation,
		"result":    "NO_SUCH_ROOM",
		"extra":     roomName,
	}

	sendJSONMessage(client, response, "Error al serializar el mensaje NO_SUCH_ROOM",
		"Error al enviar el mensaje NO_SUCH_ROOM")
}

func sendNotJoinedRoomResponse(client *Client, roomName string, operation string) {
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": operation,
		"result":    "NOT_JOINED",
		"extra":     roomName,
	}

	sendJSONMessage(client, response, "Error al serializar el mensaje NOT_JOINED", "Error al enviar el mensaje NOT_JOINED")
}
