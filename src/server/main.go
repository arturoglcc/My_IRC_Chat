package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
)

func main() {
	portPtr := flag.Int("port", 8080, "Port number for the server to listen on")
	flag.Parse()

	if *portPtr < 1 || *portPtr > 65535 {
		log.Fatalf("Invalid port number: %d. Must be between 1 and 65535.", *portPtr)
	}

	port := strconv.Itoa(*portPtr)

	fmt.Printf("Starting server on port %s...\n", port)

	// Directly call the StartServer function
	err := StartServer(port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	fmt.Println("Server stopped.")
}
