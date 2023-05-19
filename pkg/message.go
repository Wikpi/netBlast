package pkg

import (
	"time"

	"nhooyr.io/websocket"
)

// Structure of individual user message
type Message struct {
	Username    string    `json:"username"`
	Message     string    `json:"message"`
	MessageTime time.Time `json:"messageTime"`
	Color       string    `json:"color"`
}

type Name struct {
	Name string `json:"name"`
}

type User struct {
	Name string
	Conn *websocket.Conn
}
