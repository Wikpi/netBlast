package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"unicode/utf8"

	"netBlast/pkg"

	"nhooyr.io/websocket"
)

// Server struct that holds all the essential info
type serverInfo struct {
	messages    []pkg.Message
	errMsg      string
	statusCode  int
	users       map[string]string
	connections map[*websocket.Conn]string
	lock        sync.RWMutex
	mux         *http.ServeMux
}

func newServer() *serverInfo {
	server := &serverInfo{}

	server.users = make(map[string]string)
	server.connections = make(map[*websocket.Conn]string)
	server.mux = http.NewServeMux()

	return server
}

// Stores all server handlers
func (server *serverInfo) handleServer() {
	// Registers user and establishes a connection
	server.mux.HandleFunc("/register", server.registerUser)
	// Receives and handles user messages
	server.mux.HandleFunc("/message", server.handleSession)

	err := http.ListenAndServe(pkg.ServerURL, server.mux)
	pkg.HandleError(pkg.Sv, err, 0)
}

func main() {
	newServer().handleServer()
}

/* ----------------Main Handler Functions---------------- */

// Registers new users
func (s *serverInfo) registerUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	pkg.HandleError(pkg.SvRegister+pkg.BadRead, err, 1)

	name := struct {
		Name string
	}{}

	err = json.Unmarshal(body, &name)
	pkg.HandleError(pkg.SvRegister+pkg.BadParse, err, 1)

	s.lock.Lock()
	errMsg, status := checkName(name.Name, s)
	s.lock.Unlock()

	if status == http.StatusAccepted {
		s.lock.Lock()
		s.users[name.Name] = ""
		s.lock.Unlock()
	}

	w.WriteHeader(status)
	if errMsg != "" {
		jErr, err := json.Marshal(errMsg)
		pkg.HandleError(pkg.SvRegister+pkg.BadParse, err, 1)
		w.Write(jErr)
	}
}

// Handle websocket connection
func (s *serverInfo) handleSession(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	pkg.HandleError(pkg.SvMessage+pkg.BadConn, err, 1)

	defer c.Close(websocket.StatusInternalError, "")

	s.lock.Lock()
	s.connections[c] = ""
	s.lock.Unlock()

	s.readMessage(c)
}

/* ----------------Additional Functions---------------- */

// Reads received messages
func (s *serverInfo) readMessage(c *websocket.Conn) {
	for {
		message := pkg.WsRead(c, pkg.SvMessage+pkg.BadRead)
		if (message == pkg.Message{}) {
			delete(s.connections, c)
			return
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
	for ic := range s.connections {
		pkg.WsWrite(ic, message, pkg.SvMessage+pkg.BadWrite)
	}
	s.lock.RUnlock()
}

/* ----------------Standalone Functions---------------- */

// Validates username
func checkName(name string, s *serverInfo) (string, int) {
	errMsg := ""
	statusCode := 0

	if utf8.RuneCountInString(name) < 3 {
		errMsg = "Name too short. "
		statusCode = http.StatusNotAcceptable
	} else if utf8.RuneCountInString(name) > 10 {
		errMsg = "Name too long. "
		statusCode = http.StatusNotAcceptable
	}

	if _, ok := s.users[name]; !ok {
		errMsg = ""
		statusCode = http.StatusAccepted
	} else {
		errMsg = "Name already exists. "
		statusCode = http.StatusNotAcceptable
	}

	return errMsg, statusCode
}
