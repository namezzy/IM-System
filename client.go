package main

import (
	"flag"
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

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Set the server IP address(defalut is 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "Set the server port (defalut is 8888)")
}

func main() {
	//Command line parsing.
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> Failed to connect to the server.")
		return
	}

	fmt.Println(">>>>> Successfully to connected to the server.")

	// start client's business
	select {}
}
