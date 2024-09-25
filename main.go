package main

import (
	"bufio"
	"fmt"
	"net"
)

var (
	clients = make(map[string]net.Conn)
	broadcast chan string
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error creating listener:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	fmt.Println("Hello, Please enter your username:")
	name, _ := reader.ReadString('\n') // listen to input until \n
	
	clients[name] = conn
	
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			leaveMsg, _ := fmt.Printf("%s has left the chat", name)
			broadcast <- leaveMsg
			delete(clients, name)
		}
		// Broad cast message from user
		msg, _ :=  fmt.Printf("%s: %s", name, msg)
		broadcast <- msg
	}

}
