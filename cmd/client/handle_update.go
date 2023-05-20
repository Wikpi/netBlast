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
	case "quit":
		m.quitUser(value)
	case "help":
		m.handleHelp(value)
	}
}

// Registers and establishes a websocket connection with the server
func (m *model) registerNewUser(value string) {
	name := pkg.Name{Name: value}

	data := pkg.ParseToJson(name, pkg.ClRegister+pkg.BadParse)

	res := handlePostRequest(data, "http://"+pkg.ServerURL+"/register", pkg.ClRegister)

	if res.StatusCode == http.StatusAccepted {
		m.user.user.Name = value

		c, _, err := websocket.Dial(context.Background(), "ws://"+pkg.ServerURL+"/message", nil)
		pkg.HandleError(pkg.ClRegister+pkg.BadConn, err, 0)
		m.user.user.Conn = c

		m.user.user.Status = "online"
		m.screen = "chat"

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
		msg := pkg.WsRead(m.user.user.Conn, pkg.ClMessage+pkg.BadRead)

		m.lock.Lock()
		m.user.messages = append(m.user.messages, msg)
		m.lock.Unlock()
	}
}

// Writes user message to websocket connection
func (m *model) writeNewMessage(value string) {
	user := m.user.user

	message := pkg.Message{
		Username:    user.Name,
		Message:     value,
		MessageTime: time.Now(),
		Color:       user.UserColor,
	}

	pkg.WsWrite(user.Conn, message, pkg.ClMessage+pkg.BadWrite)
}

// Updates user settings
func (m *model) updateSettings(value string) {
	if strings.ToLower(value) == "color" {
		m.user.user.UserColor = getColor()
	}
}

func (m *model) listUsers(value string) {
	return
}

func (m *model) quitUser(value string) {
	if value == "Y" {
		m.user.user.Status = "offline"
		return
	} else if value == "N" {
		m.screen = "chat"
		return
	}
}

func (m *model) handleHelp(value string) {
	return
}
