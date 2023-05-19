package client

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var options []option

type option struct {
	name  string
	value string
}

func init() {
	options = []option{
		{
			name:  "Color",
			value: "update settings related to color option.",
		},
	}
}

// Routes to different UI screens
func (m *model) routeScreen() {
	switch m.screen {
	case "register":
		m.displayRegister()
	case "chat":
		m.displayChat()
	case "settings":
		m.displaySettings()
	}
}

// Displays the register screen
func (m *model) displayRegister() {
	m.ui = strings.Builder{}

	if m.user.name == "" {
		// Doesnt let thorugh, if name is invalid
		if m.err == "" {
			m.ui.WriteString("Name: \n")
			return
		}
		m.ui.WriteString(m.err + "Try again: \n")
		return
	}
	m.screen = "chat"
}

// Displays the chatroom screen
func (m *model) displayChat() {
	m.ui = strings.Builder{}

	m.lock.RLock()

	defer m.lock.RUnlock()

	m.logUserMessages()

	m.ui.WriteString("<-------------------------------------> \nMessage: \n")
}

// Displays the settings screen
func (m *model) displaySettings() {
	m.ui = strings.Builder{}

	m.ui.WriteString("Options: \n")

	for _, option := range options {
		m.ui.WriteString(option.name + ": " + option.value + "\n")
	}
}

// Display all user messages
func (m *model) logUserMessages() {
	var currentTime time.Time

	// Displays previous messages
	for _, msg := range m.user.messages {
		// Username color styler
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(msg.Color))

		// Dsiplays new time
		if currentTime.After(msg.MessageTime) {
			currentTime = msg.MessageTime

			m.ui.WriteString(currentTime.Format("15:04") + "\n")
		}

		m.ui.WriteString(style.Render(msg.Username) + ": " + msg.Message + "\n")
	}
}
