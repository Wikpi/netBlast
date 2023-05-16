package pkg

import (
	"context"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WsRead(conn *websocket.Conn, errMsg string) Message {
	message := struct {
		Username    string
		Message     string
		MessageTime time.Time
		Color       string
	}{}

	err := wsjson.Read(context.Background(), conn, &message)
	HandleError(errMsg, err, 0)

	msg := Message{
		Username:    message.Username,
		Message:     message.Message,
		MessageTime: message.MessageTime,
		Color:       message.Color,
	}

	return msg
}
