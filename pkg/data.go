package pkg

import (
	"time"

	"nhooyr.io/websocket"
)

// Server url and port
const ServerURL = "localhost:8080"

// Tool list directories
const Scrapper = "./tools/scrapper"

// Client log file
//const ClLogs = "./logs/client/logs.txt"

// Server log file
//const SvLogs = "./logs/server/logs.txt"

// Logs
const Logs = "./logs/logs.txt"

// Structure of individual user message
type Message struct {
	Username    string    `json:"username"`
	Message     string    `json:"message"`
	MessageTime time.Time `json:"messageTime"`
	Color       string    `json:"color"`
}

// Structure of name message
type Name struct {
	Name string `json:"name"`
}

// Structure of userList message
type User struct {
	Name      string          `json:"Name"`
	Conn      *websocket.Conn `json:"Conn"`
	Status    string          `json:"Status"`
	UserColor string          `json:"UserColor"`
}
