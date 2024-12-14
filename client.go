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
	// The current user's mode
	flag int
}

func NewClient(serverIp string, serverPort int) *Client {
	// Create a client object
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.Public chat mode")
	fmt.Println("2.Private chat mode")
	fmt.Println("3.Rename username")
	fmt.Println("0.Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>Please enter a number within the valid range.<<<<")
		return false
	}

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		// Handle different business operations according to different models.
		switch client.flag {
		case 1:
			// public chat mode
			fmt.Println("Choose public chat mode.....")
			break
		case 2:
			// private chat mode
			fmt.Println("Choose private chat mode")
			break
		case 3:
			// rename username
			fmt.Println("choose rename username")
			break

		}
	}
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
	client.Run()
}
