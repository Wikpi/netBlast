package pkg

import (
	"time"

	"nhooyr.io/websocket"
)

// Server url and port
const (
	ServerURL = "localhost:8080"

	// Tool list directories
	Scrapper = "./tools/scrapper"

	// Logs
	Logs = "./logs/logs.txt"
)

// Client log file
//const ClLogs = "./logs/client/logs.txt"

// Server log file
//const SvLogs = "./logs/server/logs.txt"

// Structure of individual user message
type Message struct {
	Username    string    `json:"username"`
	Message     string    `json:"message"`
	MessageTime time.Time `json:"messageTime"`
	Color       string    `json:"color"`
	MessageType string    `json:"messageType"`

	Receiver      string `json:"receiver"`
	ReceiverColor string `json:"receiverColor"`
}

// Structure of name message
type Name struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Structure of userList message
type User struct {
	Id        int             `json:"id"`
	Name      string          `json:"name"`
	Conn      *websocket.Conn `json:"conn"`
	Status    string          `json:"status"`
	UserColor string          `json:"userColor"`
}
