package pkg

import (
	"context"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WsWrite(conn *websocket.Conn, message Message, errMsg string) {
	err := wsjson.Write(context.Background(), conn, message)
	HandleError(errMsg, err, 1)
}
