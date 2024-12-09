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

// Send a message to the corresponding client of the current User.
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// User's business of processing messages
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// Query the current list of users
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + " Online...\n"
			this.SendMsg(onlineMsg)

		}
		this.server.mapLock.Unlock()

	} else {
		this.server.BroadCast(this, msg)
	}
	this.server.BroadCast(this, msg)
}

/*
Method to listen to the current User channel and
send messages directly to the remote client when there is a message
*/
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
