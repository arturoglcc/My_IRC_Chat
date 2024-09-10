package main

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
