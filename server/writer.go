package main

// Esta función es un placeholder para manejar el envío de mensajes al cliente
func writeToClient(server *Server, client *Client) {
	for {
		// Aquí podrías implementar una cola de mensajes o lógica para enviar mensajes al cliente
		// En este ejemplo, simplemente enviamos un mensaje de prueba cada 10 segundos
		client.Conn.Write([]byte("Este es un mensaje del servidor.\n"))

		// Puedes implementar la lógica para enviar mensajes específicos en respuesta
		// a comandos o eventos.
	}
}
