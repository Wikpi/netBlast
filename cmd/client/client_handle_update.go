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
	case registerScreen:
		m.handleRegister(value)
	case chatScreen:
		m.handleWrite(value)
	case settingsScreen:
		m.handleSettings(value)
	case usersScreen:
		m.handleUserList(value)
	case quitScreen:
		m.handleQuit(value)
	case helpScreen:
		m.handleHelp(value)
	}
}

// Registers and establishes a websocket connection with the server
func (m *model) handleRegister(value string) {
	name := pkg.Name{Name: value, Color: m.user.UserColor}

	data := pkg.ParseToJson(name, pkg.ClRegister+pkg.BadParse)

	res := handlePostRequest(data, "http://"+pkg.ServerURL+"/register", pkg.ClRegister)

	if res.StatusCode == http.StatusAccepted {
		m.user.Name = value

		c, _, err := websocket.Dial(context.Background(), "ws://"+pkg.ServerURL+"/message", nil)
		pkg.HandleError(pkg.ClRegister+pkg.BadConn, err, 0)
		m.user.Conn = c

		m.user.Status = "online"
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
		msg := pkg.WsRead(m.user.Conn, pkg.ClMessage+pkg.BadRead)

		m.lock.Lock()
		m.chat.messages = append(m.chat.messages, msg)
		m.lock.Unlock()
	}
}

// Writes user message to websocket connection
func (m *model) handleWrite(value string) {
	user := m.user

	message := pkg.Message{
		Username:    user.Name,
		Message:     value,
		MessageTime: time.Now(),
		Color:       user.UserColor,
	}

	pkg.WsWrite(user.Conn, message, pkg.ClMessage+pkg.BadWrite)
}

// Updates user settings
func (m *model) handleSettings(value string) {
	if strings.ToLower(value) == "color" {
		m.user.UserColor = getColor()
	}
}

func (m *model) handleUserList(value string) {
	return
}

func (m *model) handleQuit(value string) {
	if value == "Y" {
		m.user.Status = "offline"
		return
	} else if value == "N" {
		m.screen = m.prevScreen
		return
	}
}

func (m *model) handleHelp(value string) {
	return
}
