package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"unicode/utf8"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type info struct {
	users    []string
	messages []message
}

type message struct {
	username    string
	message     string
	messageTime time.Time
}

func main() {
	mux := http.NewServeMux()

	users := make(map[string]struct{})
	connections := make(map[*websocket.Conn]string)

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		checkError("Server/register: Could not read body: ", err)

		name := struct {
			Name string
		}{}

		err = json.Unmarshal(body, &name)
		checkError("Server/register: could not parse body: ", err)

		length := utf8.RuneCountInString(name.Name)
		errMsg := ""
		statusCode := http.StatusNoContent

		if length < 3 {
			errMsg = "Name too short. "
			statusCode = http.StatusNotAcceptable
		} else if length > 10 {
			errMsg = "Name too long. "
			statusCode = http.StatusNotAcceptable
		} else if _, ok := users[name.Name]; !ok {
			statusCode = http.StatusAccepted

			users[name.Name] = struct{}{}
		} else {
			errMsg = "Name already exists. "
			statusCode = http.StatusNotAcceptable
		}

		w.WriteHeader(statusCode)
		if errMsg != "" {
			jErr, err := json.Marshal(errMsg)
			checkError("Server/register: coudlnt parse name: ", err)
			w.Write(jErr)
		}
	})
	mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		checkError("Server/message: couldnt upgrade connection: ", err)

		defer c.Close(websocket.StatusInternalError, "")

		connections[c] = ""
		fmt.Println(connections, " ", users)

		// ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		// defer cancel()

		for {
			message := struct {
				Username    string    `json:"username"`
				Message     string    `json:"message"`
				MessageTime time.Time `json:"messageTime"`
				Color       string    `json:"color"`
			}{}

			err = wsjson.Read(context.Background(), c, &message)
			if err != nil {
				delete(connections, c)
				return
			}
			fmt.Println(message)
			for ic := range connections {

				err = wsjson.Write(context.Background(), ic, message)
				checkError("Server/message: couldnt write: ", err)
			}
		}
	})

	err := http.ListenAndServe(":8080", mux)
	checkError("Server: error: ", err)
}

// Checks if there is an error and exits program
func checkError(errMsg string, err error) {
	if err != nil {
		log.Fatal(errMsg, err)
	}
}
