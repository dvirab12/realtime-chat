package main

import (
	"bufio"
	"fmt"
	"net"
)

var (
	clients = make(map[net.Conn]string)
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
	
	clients[conn] = name
	

	message := []byte("Hello from client")
	n, err := conn.Write(message)
	if err != nil {
		fmt.Println("Couldn't write mesasge to %s: ", conn.RemoteAddr(), err)
		return
	}
	 if n < len(message) {
		fmt.Println("Message was not fully written %d/%d", n, len(message))
	 }
}
