package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"netBlast/pkg"

	"nhooyr.io/websocket"
)

// Server struct that holds all the essential info
type serverInfo struct {
	s        http.Server
	messages []pkg.Message
	users    []pkg.User
	lock     sync.RWMutex
	mux      *http.ServeMux
	shutdown chan os.Signal
}

// Creates server with default parameters
func newServer(sd chan os.Signal) *serverInfo {
	serverInfo := &serverInfo{}
	serverInfo.s = http.Server{Addr: pkg.ServerURL, Handler: nil}
	serverInfo.mux = http.NewServeMux()
	serverInfo.shutdown = sd

	return serverInfo
}

// Stores all server handlers
func (server *serverInfo) handleServer() {
	fmt.Println("Running server on: ", pkg.ServerURL)

	// Check if server is running
	server.mux.HandleFunc("/", ping)
	// Registers user and establishes a connection
	server.mux.HandleFunc("/register", server.registerUser)
	// Receives and handles user messages
	server.mux.HandleFunc("/message", server.handleSession)
	// Give connected user list
	server.mux.HandleFunc("/userList", server.sendUserList)

	//go server.serverShutdown()
	/* not working? */
	// err := server.s.ListenAndServe()

	err := http.ListenAndServe(pkg.ServerURL, server.mux)
	pkg.HandleError(pkg.Sv, err, 0)
}

func Server(shutdown chan os.Signal) {
	newServer(shutdown).handleServer()
}

// Shutdowns the server
func (server *serverInfo) serverShutdown() {
	signal.Notify(server.shutdown, os.Interrupt)
	<-server.shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.s.Shutdown(ctx)
	pkg.HandleError(pkg.Sv+pkg.BadClose, err, 1)
}

// Reads received messages
func (s *serverInfo) readMessage(c *websocket.Conn) {
	for {
		message := pkg.WsRead(c, pkg.SvMessage+pkg.BadRead)
		if (message == pkg.Message{}) {
			if userIdx := findUser(c, s); userIdx != -1 {
				fmt.Println("User left the server: ", s.users[userIdx].Name)
				s.users = append(s.users[:userIdx], s.users[userIdx+1:]...)
				return
			}
		}

		s.lock.Lock()
		s.messages = append(s.messages, message)
		s.lock.Unlock()

		s.writeToAll(message)
	}
}

// Writes user message to all other connections
func (s *serverInfo) writeToAll(message pkg.Message) {
	s.lock.RLock()
	for _, ic := range s.users {
		pkg.WsWrite(ic.Conn, message, pkg.SvMessage+pkg.BadWrite)
	}
	s.lock.RUnlock()
}
