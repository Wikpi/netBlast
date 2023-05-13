package pkg

import "time"

// Structure of individual user message
type Message struct {
	Username    string    `json:"username"`
	Message     string    `json:"message"`
	MessageTime time.Time `json:"messageTime"`
	Color       string    `json:"color"`
}
