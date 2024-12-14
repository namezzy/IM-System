package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	// Create a client object
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// connect server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dail error:", err)
	}
	client.conn = conn

	// return object
	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>> Failed to connect to the server.")
		return
	}

	fmt.Println(">>>>> Successfully to connected to the server.")

	// start client's business
	select {}
}
