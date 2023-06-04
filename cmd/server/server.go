package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"netBlast/pkg"
	"netBlast/tools/database"

	"nhooyr.io/websocket"
)

// Creates server with default parameters
func newServer(sd chan os.Signal) *serverInfo {
	serverInfo := &serverInfo{}
	serverInfo.s = http.Server{Addr: pkg.ServerURL, Handler: nil}

	serverInfo.mux = http.NewServeMux()

	serverInfo.shutdown = sd

	serverInfo.connections = make(map[*websocket.Conn]string)

	serverInfo.db = database.OpenDB()

	return serverInfo
}

// Stores all server handlers
func (server *serverInfo) handleServer() {
	fmt.Println("Running server on: ", pkg.ServerURL)

	// Check if server is running
	server.mux.HandleFunc(pingHandler, ping)
	// Registers user and establishes a connection
	server.mux.HandleFunc(registerHandler, server.registerUser)
	// Receives and handles user messages
	server.mux.HandleFunc(sessionHandler, server.handleSession)
	// Give connected user list
	server.mux.HandleFunc(userListHandler, server.sendUserList)
	// Send a dm from one user to another
	server.mux.HandleFunc(dmHandler, server.directMessage)

	//go server.serverShutdown()
	/* not working? */
	// err := server.s.ListenAndServe()

	err := http.ListenAndServe(pkg.ServerURL, server.mux)
	if err != nil {
		pkg.LogError(err)
		log.Fatal(pkg.Sv + pkg.BadOpen)
	}
}

func Server(shutdown chan os.Signal) {
	newServer(shutdown).handleServer()
}
