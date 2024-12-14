package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

// Process the response message from the server and display it directly to standard output.
func (client *Client) DealResponse() {
	// As soon as there is data on client.conn, directly copy it to standard output and block listening indefinitely.
	io.Copy(os.Stdout, client.conn)
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

// Query online user
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err: ", err)
		return
	}
}

// private chat mode
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>> Please enter the chat recipient's [username], type \"exit\"to quit. ")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>Please input your messages, or type \"exit\" to quit.")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// No sent if the message is null
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err: ", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>Please input your messages, or type \"exit\" to quit.")
			fmt.Scanln(&chatMsg)

		}

		client.SelectUsers()
		fmt.Println(">>>> Please enter the chat recipient's [username], type \"exit\"to quit. ")
		fmt.Scanln(&remoteName)

	}
}

func (client *Client) PublicChat() {
	// Prompt the user enter messages
	var chatMsg string

	fmt.Println(">>>>Please input your message, or type \"exit\" to quit.")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// send server

		// If the message is not empty, send it.
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err: ", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>Please input your message, or type \"exit\" to quit.")
		fmt.Scanln(&chatMsg)
	}

}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>> Please input username: ")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		// Handle different business operations according to different models.
		switch client.flag {
		case 1:
			// public chat mode
			// fmt.Println("Choose public chat mode.....")
			client.PublicChat()
			break
		case 2:
			// private chat mode
			//fmt.Println("Choose private chat mode")
			client.PrivateChat()
			break
		case 3:
			// rename username
			client.UpdateName()
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

	// Start a separate goroutine to handle the server's response messages.
	go client.DealResponse()

	fmt.Println(">>>>> Successfully to connected to the server.")

	// start client's business
	client.Run()
}
