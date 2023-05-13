package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"unicode/utf8"

	"netBlast/pkg"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Server struct that holds all the essential info
type serverInfo struct {
	messages    []pkg.Message
	m           sync.Mutex
	errMsg      string
	statusCode  int
	users       map[string]string
	connections map[*websocket.Conn]string
	lock        sync.RWMutex
	mux         *http.ServeMux
}

func main() {
	server := &serverInfo{}

	server.users = make(map[string]string)
	server.connections = make(map[*websocket.Conn]string)
	server.mux = http.NewServeMux()

	// Registers user and establishes a connection
	server.mux.HandleFunc("/register", server.registerUser)
	// Receives and handles user messages
	server.mux.HandleFunc("/message", server.handleRequest)

	err := http.ListenAndServe(pkg.ServerURL, server.mux)
	handleError(pkg.Sv, err, 0)
}

// Registers new users
func (s *serverInfo) registerUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	handleError(pkg.SvRegister+pkg.BadRead, err, 1)

	name := struct {
		Name string
	}{}

	err = json.Unmarshal(body, &name)
	handleError(pkg.SvRegister+pkg.BadParse, err, 1)

	s.lock.Lock()
	s.checkName(name.Name)
	s.lock.Unlock()

	if s.statusCode == http.StatusAccepted {
		s.lock.Lock()
		s.users[name.Name] = ""
		s.lock.Unlock()
	}

	w.WriteHeader(s.statusCode)
	if s.errMsg != "" {
		jErr, err := json.Marshal(s.errMsg)
		handleError(pkg.SvRegister+pkg.BadParse, err, 1)
		w.Write(jErr)
	}
}

// Handle websocket connection
func (s *serverInfo) handleRequest(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	handleError(pkg.SvMessage+pkg.BadConn, err, 1)

	defer c.Close(websocket.StatusInternalError, "")

	s.lock.Lock()
	s.connections[c] = ""
	s.lock.Unlock()

	s.readMessage(c)
}

// Validates username
func (s *serverInfo) checkName(name string) {
	if utf8.RuneCountInString(name) < 3 {
		s.errMsg = "Name too short. "
		s.statusCode = http.StatusNotAcceptable
	} else if utf8.RuneCountInString(name) > 10 {
		s.errMsg = "Name too long. "
		s.statusCode = http.StatusNotAcceptable
	} else if _, ok := s.users[name]; !ok {
		s.statusCode = http.StatusAccepted
	} else {
		s.errMsg = "Name already exists. "
		s.statusCode = http.StatusNotAcceptable
	}
}

// Reads received messages
func (s *serverInfo) readMessage(c *websocket.Conn) {
	for {
		msg := struct {
			Username    string
			Message     string
			MessageTime time.Time
			Color       string
		}{}

		err := wsjson.Read(context.Background(), c, &msg)
		if err != nil {
			delete(s.connections, c)
			return
		}

		message := pkg.Message{
			Username:    msg.Username,
			Message:     msg.Message,
			MessageTime: msg.MessageTime,
			Color:       msg.Color,
		}
		s.lock.Lock()
		s.messages = append(s.messages, message)
		s.lock.Unlock()

		s.lock.RLock()
		s.writeToAll(message)
		s.lock.RUnlock()
	}
}

// Writes user message to all other connections
func (s *serverInfo) writeToAll(message pkg.Message) {
	for ic := range s.connections {
		err := wsjson.Write(context.Background(), ic, message)
		handleError(pkg.SvMessage+pkg.BadWrite, err, 1)
	}
}

// Handles incoming error
func handleError(errMsg string, incomingErr error, exc ...int) {
	if incomingErr == nil {
		return
	}
	file, err := os.OpenFile(pkg.SvLogs, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Print(err)
	}
	defer file.Close()

	// Writes error to logs file
	if _, err := file.WriteString(time.Now().Format("2006-01-13 12:60") + " " + incomingErr.Error() + "\n\n"); err != nil {
		fmt.Println(err)
	}

	// Exits program and gives message where error occured
	switch exc[0] {
	case 0:
		log.Fatal(errMsg)
	case 1:
		fmt.Println(errMsg)
	}
}
