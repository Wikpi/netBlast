package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
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

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Server/register: Could not read body: ", err)
		}
		name := string(body)
		length := utf8.RuneCountInString(name) - 2
		errMsg := ""
		statusCode := http.StatusNoContent

		if length < 3 {
			errMsg = "Name too short. "
			statusCode = http.StatusNotAcceptable
		} else if length > 10 {
			errMsg = "Name too long. "
			statusCode = http.StatusNotAcceptable
		} else if _, ok := users[name]; !ok {
			users[name] = struct{}{}
			statusCode = http.StatusAccepted
		} else {
			errMsg = "Name already exists. "
			statusCode = http.StatusNotAcceptable
		}

		w.WriteHeader(statusCode)
		if errMsg != "" {
			jErr, err := json.Marshal(errMsg)
			if err != nil {
				log.Fatal("Server/register: coudlnt parse name: ", err)
			}
			w.Write(jErr)
		}
	})
	mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Server/message: Could not read body: ", err)
		}
		message := string(body)

		io.WriteString(w, message)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Server: error: ", err)
	}
}
