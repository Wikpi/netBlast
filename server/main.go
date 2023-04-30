package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

		if _, ok := users[name]; ok {
			io.WriteString(w, "Name already exists")
		} else if len(name) < 2 {
			io.WriteString(w, "Name too short")
		} else if len(name) > 10 {
			io.WriteString(w, "name too long")
		} else {
			users[name] = struct{}{}
			io.WriteString(w, "Registered! Hello ")
			io.WriteString(w, name)
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
