package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"unicode/utf8"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type serverModel struct {
	//users    []string // For database
	messages []message
}

type message struct {
	Username    string    `json:"username"`
	Message     string    `json:"message"`
	MessageTime time.Time `json:"messageTime"`
	Color       string    `json:"color"`
}

func main() {
	mux := http.NewServeMux()

	// Stores usernames and connections // Potentially use Database in the future?
	users := make(map[string]string)
	connections := make(map[*websocket.Conn]string)

	var server serverModel

	// Registers user and establishes a connection
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		handleError("Server/register: Could not read body: ", err, 1)

		name := struct {
			Name string
		}{}

		err = json.Unmarshal(body, &name)
		handleError("Server/register: could not parse body: ", err, 1)

		errMsg, statusCode := checkName(name.Name, users)

		if statusCode == http.StatusAccepted {
			users[name.Name] = ""
		}

		w.WriteHeader(statusCode)
		if errMsg != "" {
			jErr, err := json.Marshal(errMsg)
			handleError("Server/register: coudlnt parse name: ", err, 1)
			w.Write(jErr)
		}
	})
	// Receives and handles user messages
	mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		handleError("Server/message: couldnt upgrade connection: ", err, 1)

		defer c.Close(websocket.StatusInternalError, "")

		connections[c] = ""

		server.readMessage(&connections, c)
	})

	err := http.ListenAndServe("localhost:8080", mux)
	handleError("Server: error: ", err, 0)
}

// Validates user name
func checkName(name string, users map[string]string) (string, int) {
	statusCode := http.StatusNoContent
	errMsg := ""

	if utf8.RuneCountInString(name) < 3 {
		errMsg = "Name too short. "
		statusCode = http.StatusNotAcceptable
	} else if utf8.RuneCountInString(name) > 10 {
		errMsg = "Name too long. "
		statusCode = http.StatusNotAcceptable
	} else if _, ok := users[name]; !ok {
		statusCode = http.StatusAccepted
	} else {
		errMsg = "Name already exists. "
		statusCode = http.StatusNotAcceptable
	}

	return errMsg, statusCode
}

// Read user sent message
func (sm *serverModel) readMessage(connections *map[*websocket.Conn]string, c *websocket.Conn) {
	for {
		msg := struct {
			Username    string
			Message     string
			MessageTime time.Time
			Color       string
		}{}

		err := wsjson.Read(context.Background(), c, &msg)
		if err != nil {
			delete(*connections, c)
			return
		}

		message := message{
			Username:    msg.Username,
			Message:     msg.Message,
			MessageTime: msg.MessageTime,
			Color:       msg.Color,
		}
		sm.messages = append(sm.messages, message)

		writeToAll(*connections, message)
	}
}

// Writes user message to all other connections
func writeToAll(connections map[*websocket.Conn]string, message message) {
	for ic := range connections {
		err := wsjson.Write(context.Background(), ic, message)
		handleError("Server/message: couldnt write: ", err, 1)
	}
}

// Handles incoming error
func handleError(errMsg string, pErr error, exc ...int) {
	if pErr != nil {
		file, err := os.OpenFile("./logs/server/logs.txt", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Print(err)
		}
		defer file.Close()

		// Writes error to logs file
		if _, err := file.WriteString(pErr.Error()); err != nil {
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
}
