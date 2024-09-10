package main

import (
	"encoding/json"
	"log"
)

// sendInvalidMessageResponse envía un mensaje de error al cliente y cierra la conexión.
func sendInvalidMessageResponse(client *Client) {
	// Crear la respuesta JSON
	response := map[string]string{
		"type":      "RESPONSE",
		"operation": "INVALID",
		"result":    "INVALID",
	}

	// Serializar la respuesta JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error al serializar la respuesta JSON: %v", err)
		return
	}

	// Enviar la respuesta al cliente
	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("Error al enviar la respuesta al cliente %s: %v", client.ID, err)
	}

	// Cerrar la conexión con el cliente
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

		// Validar que el mensaje sea un JSON válido
		var msg map[string]interface{}
		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			// El mensaje no es un JSON válido, enviar respuesta y desconectar
			sendInvalidMessageResponse(client)
			break // Salir del bucle para terminar la goroutine
		}

		// Procesar el mensaje recibido
		processClientMessage(server, client, msg)
	}
}

// Función que escucha los mensajes del cliente
func processClientMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verificar que el mensaje tenga el campo "type"
	msgType, ok := msg["type"].(string)
	if !ok {
		// Si el campo "type" no es una cadena, enviar respuesta de error y desconectar
		sendInvalidMessageResponse(client)
		return
	}

	// Manejar diferentes tipos de mensajes
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
		// Si el tipo de mensaje no es reconocido, enviar respuesta de error y desconectar
		sendInvalidMessageResponse(client)
	}
}

func handleStatusMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verificar que el campo "status" esté presente y sea una cadena
	newStatus, ok := msg["status"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Verificar si el estado es válido
	validStatuses := map[string]bool{
		"ACTIVE": true,
		"BUSY":   true,
		"AWAY":   true,
	}

	if _, valid := validStatuses[newStatus]; !valid {
		// Si el estado no es válido, enviar una respuesta de error y desconectar
		sendInvalidMessageResponse(client)
		return
	}

	// Cambiar el estado del cliente al nuevo estado
	client.Status = newStatus
	log.Printf("El cliente %s ha cambiado su estado a %s", client.ID, newStatus)

	// Notificar a todos los demás clientes sobre el nuevo estado del cliente
	notifyClientsOfStatusChange(server, client)
}

func handlePublicTextMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verificar que el campo "text" esté presente
	publicText, ok := msg["text"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Enviar el mensaje a todos los demás usuarios
	sendPublicTextToAll(server, client, publicText)
}

func handleNewRoomMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verificar que el campo "roomname" esté presente
	roomName, ok := msg["roomname"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Verificar la longitud del nombre de la sala
	if len(roomName) > 16 {
		sendRoomAlreadyExistsResponse(client, roomName)
		return
	}

	// Verificar si la sala ya existe
	server.RoomMu.Lock()
	_, exists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if exists {
		sendRoomAlreadyExistsResponse(client, roomName)
		return
	}

	// Crear la sala y agregar al cliente
	createRoom(server, client, roomName)
}

func handleInviteMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verificar que los campos "roomname" y "usernames" estén presentes
	roomName, ok := msg["roomname"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	invitedUsernames, ok := msg["usernames"].([]interface{})
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Verificar si la sala existe y si el cliente pertenece a la sala
	server.RoomMu.Lock()
	room, roomExists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if !roomExists || room.Members[client.ID] == nil {
		sendNoSuchRoomResponse(client, roomName, "INVITE")
		return
	}

	// Verificar que todos los usuarios invitados existan
	for _, invitedUsername := range invitedUsernames {
		username, valid := invitedUsername.(string)
		if !valid {
			sendInvalidMessageResponse(client)
			return
		}

		server.Mu.Lock()
		_, userExists := server.Clients[username]
		server.Mu.Unlock()

		if !userExists {
			sendNoSuchUserResponse(client, "INVITE", username)
			return
		}
	}

	// Enviar la invitación a los usuarios
	inviteUsersToRoom(server, client, roomName, invitedUsernames)
}

func handleJoinRoomMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verificar que el campo "roomname" esté presente
	roomName, ok := msg["roomname"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Verificar si la sala existe
	server.RoomMu.Lock()
	room, roomExists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if !roomExists {
		sendNoSuchRoomResponse(client, roomName, "JOIN_ROOM")
		return
	}

	// Verificar si el usuario fue invitado
	if !room.Invited[client.ID] {
		sendNotInvitedResponse(client, roomName)
		return
	}

	// Unir al cliente a la sala y eliminarlo de la lista de invitados
	server.RoomMu.Lock()
	delete(room.Invited, client.ID)
	room.Members[client.ID] = client
	server.RoomMu.Unlock()

	// Enviar la respuesta de éxito al cliente
	sendJoinRoomSuccessResponse(client, roomName)

	// Notificar a los demás miembros de la sala
	notifyRoomMembersUserJoined(server, room, client)
}

func handleRoomUsersMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verificar que el campo "roomname" esté presente
	roomName, ok := msg["roomname"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Verificar si la sala existe
	server.RoomMu.Lock()
	room, roomExists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if !roomExists {
		sendNoSuchRoomResponse(client, roomName, "ROOM_USERS")
		return
	}

	// Verificar si el usuario está en la sala
	if room.Members[client.ID] == nil {
		sendNotJoinedRoomResponse(client, roomName, "ROOM_USERS")
		return
	}

	// Enviar la lista de usuarios de la sala
	sendRoomUserList(server, client, room)
}

func handleRoomTextMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verificar que los campos "roomname" y "text" estén presentes
	roomName, ok := msg["roomname"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	text, ok := msg["text"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Verificar si la sala existe
	server.RoomMu.Lock()
	room, roomExists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if !roomExists {
		sendNoSuchRoomResponse(client, roomName, "ROOM_TEXT")
		return
	}

	// Verificar si el usuario está en la sala
	if room.Members[client.ID] == nil {
		sendNotJoinedRoomResponse(client, roomName, "ROOM_TEXT")
		return
	}

	// Enviar el mensaje a los demás usuarios en la sala
	broadcastRoomTextMessage(server, client, room, text)
}

func handleLeaveRoomMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verificar que el campo "roomname" esté presente
	roomName, ok := msg["roomname"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Verificar si la sala existe
	server.RoomMu.Lock()
	room, roomExists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if !roomExists {
		sendNoSuchRoomResponse(client, roomName, "LEAVE_ROOM")
		return
	}

	// Verificar si el usuario está en la sala
	if room.Members[client.ID] == nil {
		sendNotJoinedRoomResponse(client, roomName, "LEAVE_ROOM")
		return
	}

	// Eliminar al cliente de los miembros de la sala
	server.RoomMu.Lock()
	delete(room.Members, client.ID)
	server.RoomMu.Unlock()

	// Notificar a los demás miembros de la sala
	notifyRoomMembersUserLeft(server, client, room)

	// Eliminar el cuarto si está vacío
	if len(room.Members) == 0 {
		delete(server.Rooms, roomName)
	}
}

func handleDisconnectMessage(server *Server, client *Client) {
	// Notificar a todos los usuarios del servidor que este usuario se ha desconectado
	notifyAllUsersDisconnected(server, client)

	// Eliminar al usuario de todos los cuartos en los que estaba y notificar a los miembros
	server.RoomMu.Lock()
	for _, room := range server.Rooms {
		if _, exists := room.Members[client.ID]; exists {
			// Notificar a los demás miembros del cuarto
			notifyRoomMembersUserLeft(server, client, room)
			// Eliminar al cliente de los miembros del cuarto
			delete(room.Members, client.ID)

			// Verificar si el cuarto se queda vacío y no es el cuarto general
			if len(room.Members) == 0 && room.Name != "General" {
				delete(server.Rooms, room.Name)
			}
		}
	}
	server.RoomMu.Unlock()

	// Eliminar al usuario de la lista de usuarios del servidor
	server.Mu.Lock()
	delete(server.Clients, client.ID)
	server.Mu.Unlock()

	// Cerrar la conexión con el cliente y finalizar la goroutine
	client.Conn.Close()
}
