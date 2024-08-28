// server_test.go
package main

import (
	"net"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	address := "localhost:8080"
	s := NewServer(address)

	if s == nil {
		t.Fatal("Expected server to be initialized, got nil")
	}

	if s.Address != address {
		t.Errorf("Expected server address to be %s, got %s", address, s.Address)
	}

	if len(s.Clients) != 0 {
		t.Errorf("Expected server clients map to be empty, got length %d", len(s.Clients))
	}

	if len(s.Rooms) != 0 {
		t.Errorf("Expected server rooms map to be empty, got length %d", len(s.Rooms))
	}
}

func TestServerStart(t *testing.T) {
	address := "localhost:8081"
	s := NewServer(address)

	// Run the server in a goroutine so that it doesn't block the test
	go s.Start()

	// Give the server a moment to start up
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Server failed to start: %v", err)
	}
	conn.Close()
}
