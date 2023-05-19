package client

import (
	"context"
	"io/ioutil"
	"net/http"
	"netBlast/pkg"
	"strings"
	"time"

	"nhooyr.io/websocket"
)

// Routes message execution depending on the screen
func (m *model) routeMessage() {
	value := m.input.Value()
	if value == "" {
		return
	}

	switch m.screen {
	case "register":
		m.registerNewUser(value)
	case "chat":
		m.writeNewMessage(value)
	case "settings":
		m.updateSettings(value)
	case "users":
		m.listUsers(value)
	}
}

// Registers and establishes a websocket connection with the server
func (m *model) registerNewUser(value string) {
	name := pkg.Name{Name: value}

	data := pkg.ParseToJson(name, pkg.ClRegister+pkg.BadParse)

	res := handlePostRequest(data, "http://"+pkg.ServerURL+"/register", pkg.ClRegister)

	if res.StatusCode == http.StatusAccepted {
		m.user.name = value

		c, _, err := websocket.Dial(context.Background(), "ws://"+pkg.ServerURL+"/message", nil)
		pkg.HandleError(pkg.ClRegister+pkg.BadConn, err, 0)
		m.user.conn = c

		go m.receiveNewMessages()
		return
	}

	// Gives an error if registration failed
	resBody, err := ioutil.ReadAll(res.Body)
	pkg.HandleError(pkg.ClRegister+pkg.BadRead, err, 0)

	pkg.ParseFromJson(resBody, &m.err, pkg.ClRegister+pkg.BadParse)

	res.Body.Close()
}

// Stores messages received from the websocket connection
func (m *model) receiveNewMessages() {
	for {
		msg := pkg.WsRead(m.user.conn, pkg.ClMessage+pkg.BadRead)

		m.lock.Lock()
		m.user.messages = append(m.user.messages, msg)
		m.lock.Unlock()
	}
}

// Writes user message to websocket connection
func (m *model) writeNewMessage(value string) {
	message := pkg.Message{
		Username:    m.user.name,
		Message:     value,
		MessageTime: time.Now(),
		Color:       m.user.userColor,
	}

	pkg.WsWrite(m.user.conn, message, pkg.ClMessage+pkg.BadWrite)
}

// Updates user settings
func (m *model) updateSettings(value string) {
	if strings.ToLower(value) == "color" {
		m.user.userColor = getColor()
	}
}

func (m *model) listUsers(value string) {
	return
}
