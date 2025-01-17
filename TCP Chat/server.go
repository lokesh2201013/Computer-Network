package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var clients []net.Conn
var mu sync.Mutex

// Handle individual client connections
func handleClient(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected: %s\n", clientAddr)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Client disconnected: %s\n", clientAddr)
			removeClient(conn)
			return
		}

		message = strings.TrimSpace(message)
		fmt.Printf("Message from %s: %s\n", clientAddr, message)

		if message == "exit" {
			fmt.Printf("Closing connection with: %s\n", clientAddr)
			removeClient(conn)
			return
		}

		broadcastMessage(fmt.Sprintf("%s says: %s\n", clientAddr, message), conn)
	}
}

// Broadcast a message to all connected clients except the sender
func broadcastMessage(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for _, client := range clients {
		if client != sender {
			client.Write([]byte(message))
		}
	}
}

// Remove a disconnected client from the client list
func removeClient(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for i, client := range clients {
		if client == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port 8080.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		mu.Lock()
		clients = append(clients, conn)
		mu.Unlock()

		go handleClient(conn)
	}
}
