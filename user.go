package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// User's constructor
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}
	// Start a goroutine that listens for messages from the current user channel.
	go user.ListenMessage()

	return user

}

// User's online business
func (this *User) Online() {

	// User online, add user to OnlineMap.
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// Broadcasting current user online message
	this.server.BroadCast(this, "Online")

}

// User's offline business
func (this *User) Offline() {

	// User offline delete user from OnlineMap
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// Broadcasting current user offline message.
	this.server.BroadCast(this, "Offline")
}

// User's business of processing messages
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
