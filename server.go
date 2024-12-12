package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Create class of Server
type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

// Create a server interface
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg

}

func (this *Server) Handler(conn net.Conn) {
	// ...Currently connected businesses
	// fmt.Println("Connection established successfully. ")

	user := NewUser(conn, this)
	user.Online()

	// Use a dedicated timer.
	timer := time.NewTimer(10 * time.Second)

	// Timer reset channel.
	resetTime := make(chan bool, 1)

	// Listen for user activity
	//isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err: ", err)
				return
			}
			//Extracting user's message(removing '\n')
			msg := string(buf[:n-1])

			// The user is processing the message for 'msg'
			user.DoMessage(msg)

			// Any message from a user indicates that the user is currently active
			// isLive <- true

			// Send a signal to the timer reset channel.
			select {
			case resetTime <- true:
			default:
			}

		}
	}()

	//  The current handler is blocking.
	for {
		select {
		case <-resetTime:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(10 * time.Second)

		case <-timer.C:
			// It's already overtime
			//  Force logout the current user.
			user.SendMsg("You have been logged out duo to inactivity.")

			// Delete from the online user list
			this.mapLock.Lock()
			delete(this.OnlineMap, user.Name)
			this.mapLock.Unlock()

			// Close user channel
			close(user.C)

			// Close connect
			conn.Close()

			// Exit the current handler
			return // runtime.Goexit()

		}
	}

}

// Start server interface
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// close listen socket
	defer listener.Close()

	go this.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err: ", err)
			continue
		}

		// do handler
		go this.Handler(conn)

	}

}
