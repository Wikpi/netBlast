package server

import (
	"io/ioutil"
	"net/http"
	"sync"
	"unicode/utf8"

	"netBlast/pkg"

	"nhooyr.io/websocket"
)

// Server struct that holds all the essential info
type serverInfo struct {
	messages []pkg.Message
	users    []user
	lock     sync.RWMutex
	mux      *http.ServeMux
}

type user struct {
	name string
	conn *websocket.Conn
}

func newServer() *serverInfo {
	server := &serverInfo{}
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

func Server() {
	newServer().handleServer()
}

/* ----------------Main Handler Functions---------------- */

// Registers new users
func (s *serverInfo) registerUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	pkg.HandleError(pkg.SvRegister+pkg.BadRead, err, 1)

	name := pkg.Name{}

	pkg.ParseFromJson(body, &name, pkg.SvRegister+pkg.BadParse)

	s.lock.Lock()
	errMsg, status := checkName(name.Name, s)
	s.lock.Unlock()

	if status == http.StatusAccepted {
		client := user{
			name: name.Name,
		}

		s.lock.Lock()
		s.users = append(s.users, client)
		s.lock.Unlock()
	}

	w.WriteHeader(status)
	if errMsg != "" {
		data := pkg.ParseToJson(errMsg, pkg.SvRegister+pkg.BadParse)
		w.Write(data)
	}
}

// Handle websocket connection
func (s *serverInfo) handleSession(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	pkg.HandleError(pkg.SvMessage+pkg.BadConn, err, 1)

	defer c.Close(websocket.StatusInternalError, "")

	s.lock.Lock()
	s.users[len(s.users)-1].conn = c
	s.lock.Unlock()

	s.readMessage(c)
}

/* ----------------Additional Functions---------------- */

// Reads received messages
func (s *serverInfo) readMessage(c *websocket.Conn) {
	for {
		message := pkg.WsRead(c, pkg.SvMessage+pkg.BadRead)
		if (message == pkg.Message{}) {
			if userIdx := findUser(c, s); userIdx != -1 {
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
		pkg.WsWrite(ic.conn, message, pkg.SvMessage+pkg.BadWrite)
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
	} else {
		if user := findUser(name, s); user == -1 {
			errMsg = ""
			statusCode = http.StatusAccepted
		} else {
			errMsg = "Name already exists. "
			statusCode = http.StatusNotAcceptable
		}
	}

	return errMsg, statusCode
}

// Finds name in user slice
func findUser(key interface{}, s *serverInfo) int {
	for idx, user := range s.users {
		if user.name == key || user.conn == key {
			return idx
		}
	}
	return -1
}
