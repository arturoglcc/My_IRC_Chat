package main

import (
	"encoding/json"
	"log"
)

// sendJSONMessage serialice a message and sends it to a client
func sendJSONMessage(client *Client, message interface{}, serializationErrMsg, sendErrMsg string) {
	jsonResponse, err := json.Marshal(message)
	if err != nil {
		log.Printf("%s: %v", serializationErrMsg, err)
		return
	}

	_, err = client.Conn.Write(jsonResponse)
	if err != nil {
		log.Printf("%s a %s: %v", sendErrMsg, client.ID, err)
	}
}

func sendToRoomMembers(sender *Client, room *Room, message interface{}, serializationErrMsg, sendErrMsg string) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("%s: %v", serializationErrMsg, err)
		return
	}

	for _, member := range room.Members {
		if member.ID != sender.ID {
			_, err := member.Conn.Write(jsonMessage)
			if err != nil {
				log.Printf("%s a %s: %v", sendErrMsg, member.ID, err)
			}
		}
	}
}

func sendMessageToAll(server *Server, sender *Client, message interface{}, excludeSender bool) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar el mensaje: %v", err)
		return
	}

	server.Mu.Lock()
	defer server.Mu.Unlock()

	for _, otherClient := range server.Clients {
		if excludeSender && otherClient.ID == sender.ID {
			continue
		}
		_, err := otherClient.Conn.Write(jsonMessage)
		if err != nil {
			log.Printf("Error al enviar mensaje a %s: %v", otherClient.ID, err)
		}
	}
}
