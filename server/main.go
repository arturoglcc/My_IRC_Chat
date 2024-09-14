package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
)

func main() {
	portPtr := flag.Int("port", 1234, "Port number for the server to listen on")
	flag.Parse()

	// Check if the provided port is within the valid range
	if *portPtr < 1 || *portPtr > 65535 {
		log.Fatalf("Invalid port number: %d. Must be between 1 and 65535.", *portPtr)
	}

	// Convert the port number to a string
	port := strconv.Itoa(*portPtr)
	fmt.Printf("Starting server on port %s...\n", port)

	// Use the provided port in the server initialization
	server := NewServer("0.0.0.0:" + port) // Initialize the server with the specified port
	server.Start()                         // Start the server

	fmt.Println("Server stopped.")
}
