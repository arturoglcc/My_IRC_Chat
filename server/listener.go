package main

import (
	"log"
)

// Función que escucha los mensajes del cliente
func listenToClientMessages(server *Server, client *Client) {
	buf := make([]byte, 1024)

	for {
		n, err := client.Conn.Read(buf)
		if err != nil {
			log.Printf("Error leyendo mensaje del cliente %s: %v", client.ID, err)
			break
		}
		message := string(buf[:n])

		// Aquí puedes agregar la lógica para procesar comandos como unirse a cuartos
		log.Printf("Mensaje recibido de %s: %s", client.ID, message)

		// Procesar el mensaje (puedes agregar más lógica aquí)
		processClientMessage(server, client, message)
	}

	// Cuando el cliente se desconecte, removerlo de los cuartos
	removeClientFromRooms(server, client)
}
