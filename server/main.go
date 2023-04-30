package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleMessage)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("server error: ", err)
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, Bob here from bob vance refrigerators")
}
