package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn net.Conn
	username string
	isAdmin bool
}

var (
	clients = []Client{}
	broadcast = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on :8080")

	go broadcastMsg()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error creating listener:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func broadcastMsg() {
	for {
		msg := <-broadcast

		for _, client := range clients {
			_, err := fmt.Fprintf(client.conn, "%s", msg)
			if err != nil {
				fmt.Printf("Error sending msg to %s: %v\n", client.username, err)
				client.conn.Close()
				removeClient(client.username)
			}
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	fmt.Fprintf(conn, "Enter your username: ")
	username, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("Error reading %s username: %v\n", conn.RemoteAddr(), err)
		return
	}

	username = strings.TrimSpace(username)

	clients = append(clients, Client{conn, username, false})	

	joinMsg :=fmt.Sprintf("%s has join the chat!\n", username)
	broadcast <- joinMsg

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			leaveMsg := fmt.Sprintf("%s has left the chat\n", username)
			broadcast <- leaveMsg
			removeClient(username)
			return
		}

		msg = strings.TrimSpace(msg)
		// Broad cast message from user
		broadcastMsg :=  fmt.Sprintf("%s: %s", username, msg)
		broadcast <- broadcastMsg
	}
}

// removeClient removes a client from the clients slice
func removeClient(username string) {
    for i, client := range clients {
        if client.username == username {
            // Remove the client from the slice
            clients = append(clients[:i], clients[i+1:]...) // Slice removal
            return
        }
    }
}
