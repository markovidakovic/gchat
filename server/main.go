package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

// Client represents a connected client
type Client struct {
	conn net.Conn
	name string
}

var clients []*Client
var clientsMtx sync.Mutex

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started. Listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		client := &Client{conn: conn}
		clients = append(clients, client)
		go handleClient(client)
	}
}

func handleClient(client *Client) {
	defer client.conn.Close()

	fmt.Printf("Received connection from %v\n", client.conn.RemoteAddr())

	// Without the '\n' character the initial message is not send due to the tcp protocol
	client.conn.Write([]byte("Enter your name: \n"))

	scanner := bufio.NewScanner(client.conn)
	scanner.Scan()
	client.name = scanner.Text()

	broadcast(client.name + " has joined the chat\n")
	fmt.Println(client.name + " has joined the chat")

	for {
		scanner := bufio.NewScanner(client.conn)
		for scanner.Scan() {
			msg := scanner.Text()
			broadcast(client.name + ": " + msg + "\n")
		}
	}
}

func broadcast(message string) {
	clientsMtx.Lock()
	defer clientsMtx.Unlock()

	for _, client := range clients {
		_, err := client.conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error broadcasting message to", client.name, ":", err)
		}
	}
}
