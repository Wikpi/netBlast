package pkg

import (
	"context"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Read message contents from connection
func WsRead(conn *websocket.Conn, errMsg string) Message {
	message := Message{}

	wsjson.Read(context.Background(), conn, &message)
	// No need to handle the error, since it would only occur if the connection was closed off cleanly
	//HandleError(errMsg, err, 2)

	return message
}

// Write message content to connection
func WsWrite(conn *websocket.Conn, message Message, errMsg string) {
	err := wsjson.Write(context.Background(), conn, message)
	HandleError(errMsg, err, 1)
}
