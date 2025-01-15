package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Handle errors and exit if necessary
func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Continuously read messages from the server
func readMessages(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Server disconnected.")
			os.Exit(0)
		}
		fmt.Print(message)
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	handleError(err)
	defer conn.Close()

	fmt.Println("Connected to server.")

	// Start a goroutine to read messages from the server
	go readMessages(conn)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "exit" {
			fmt.Println("Closing connection.")
			break
		}

		_, err := conn.Write([]byte(message + "\n"))
		handleError(err)
	}
}
