package main

import (
	"log"
	"time"
)

func handleStatusMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verify status field is present and in a string
	newStatus, ok := msg["status"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Verify valid status
	validStatuses := map[string]bool{
		"ACTIVE": true,
		"BUSY":   true,
		"AWAY":   true,
	}

	if _, valid := validStatuses[newStatus]; !valid {
		sendInvalidMessageResponse(client)
		return
	}

	// Change the status of the client and notify to other users
	client.Status = newStatus
	log.Printf("El cliente %s ha cambiado su estado a %s", client.ID, newStatus)
	notifyClientsOfStatusChange(server, client)
}

func handlePublicTextMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verify if the field "text" is present
	publicText, ok := msg["text"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	sendPublicTextToAll(server, client, publicText)
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

func handleNewRoomMessage(server *Server, client *Client, msg map[string]interface{}) {
	//Verify if the field "roomname" is present
	roomName, ok := msg["roomname"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	if len(roomName) > 16 {
		sendRoomAlreadyExistsResponse(client, roomName)
		return
	}

	// Verify if the room already exist
	server.RoomMu.Lock()
	_, exists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if exists {
		sendRoomAlreadyExistsResponse(client, roomName)
		return
	}

	// Create the room and join the client
	createRoom(server, client, roomName)
}

func handleInviteMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verify the fields "roomname" and "usernames" are present
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

	// Verify the room exists and the client is in it
	server.RoomMu.Lock()
	room, roomExists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if !roomExists || room.Members[client.ID] == nil {
		sendNoSuchRoomResponse(client, roomName, "INVITE")
		return
	}

	// Verify all the invited clients exists
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

	inviteUsersToRoom(server, client, roomName, invitedUsernames)
}

func handleJoinRoomMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verify the field "roomname exists"
	roomName, ok := msg["roomname"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Verify if the room exists
	server.RoomMu.Lock()
	room, roomExists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if !roomExists {
		sendNoSuchRoomResponse(client, roomName, "JOIN_ROOM")
		return
	}

	// Verify if the user has been invited
	if !room.Invited[client.ID] {
		sendNotInvitedResponse(client, roomName)
		return
	}

	// Join the client to the room and delate them from the invited list
	server.RoomMu.Lock()
	delete(room.Invited, client.ID)
	room.Members[client.ID] = client
	server.RoomMu.Unlock()

	sendJoinRoomSuccessResponse(client, roomName)
	notifyRoomMembersUserJoined(server, room, client)
}

func handleRoomUsersMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verify the field "roomname" is present
	roomName, ok := msg["roomname"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	// Verify if the rooms exists
	server.RoomMu.Lock()
	room, roomExists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if !roomExists {
		sendNoSuchRoomResponse(client, roomName, "ROOM_USERS")
		return
	}

	//Verify if the user is in the room
	if room.Members[client.ID] == nil {
		sendNotJoinedRoomResponse(client, roomName, "ROOM_USERS")
		return
	}

	sendRoomUserList(server, client, room)
}

func handleRoomTextMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verify if the field "roomname" and "text" are present
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

	// Verify if the room exists
	server.RoomMu.Lock()
	room, roomExists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if !roomExists {
		sendNoSuchRoomResponse(client, roomName, "ROOM_TEXT")
		return
	}

	// Verify if the user is the room
	if room.Members[client.ID] == nil {
		sendNotJoinedRoomResponse(client, roomName, "ROOM_TEXT")
		return
	}

	broadcastRoomTextMessage(server, client, room, text)
}

func handleLeaveRoomMessage(server *Server, client *Client, msg map[string]interface{}) {
	// Verify if the field "roomname" is present
	roomName, ok := msg["roomname"].(string)
	if !ok {
		sendInvalidMessageResponse(client)
		return
	}

	//Verify if the room exists
	server.RoomMu.Lock()
	room, roomExists := server.Rooms[roomName]
	server.RoomMu.Unlock()

	if !roomExists || roomName == "general" {
		sendNoSuchRoomResponse(client, roomName, "LEAVE_ROOM")
		return
	}

	// Verify if the user is in the room
	if room.Members[client.ID] == nil {
		sendNotJoinedRoomResponse(client, roomName, "LEAVE_ROOM")
		return
	}

	// Eliminate the user from the room members
	server.RoomMu.Lock()
	delete(room.Members, client.ID)
	server.RoomMu.Unlock()

	notifyRoomMembersUserLeft(server, client, room)

	if len(room.Members) == 0 {
		delete(server.Rooms, roomName)
	}
}

func handleDisconnectMessage(server *Server, client *Client) {
	notifyAllUsersDisconnected(server, client)

	// Eliminate the user from all the rooms and notify room members
	server.RoomMu.Lock()
	for _, room := range server.Rooms {
		if _, exists := room.Members[client.ID]; exists {
			time.Sleep(10 * time.Millisecond)
			notifyRoomMembersUserLeft(server, client, room)
			delete(room.Members, client.ID)
			if room.Name == "General" && len(room.Members) == 0 {
				server.RoomMu.Unlock()
				terminateServer()
				return
			}

			if len(room.Members) == 0 {
				delete(server.Rooms, room.Name)
			}
		}
	}
	server.RoomMu.Unlock()

	server.Mu.Lock()
	delete(server.Clients, client.ID)
	server.Mu.Unlock()
	client.Conn.Close()
}
