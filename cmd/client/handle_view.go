package client

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var settingOptions []option
var helpOptions []option

type option struct {
	name  string
	value string
}

func init() {
	settingOptions = []option{
		{
			name:  "Color",
			value: "update settings related to color option.",
		},
	}

	helpOptions = []option{
		{
			name:  "Help",
			value: "press -CtrlH to see help.",
		},
		{
			name:  "Settings",
			value: "press -CtrlX to enter settings.",
		},
		{
			name:  "Users",
			value: "press -CtrlC to see users.",
		},
		{
			name:  "Quit",
			value: "press -Esc to quit.",
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
	case "users":
		m.displayUsers()
	case "quit":
		m.displayQuit()
	case "help":
		m.displayHelp()
	}
}

// Displays the register screen
func (m *model) displayRegister() {
	m.ui = strings.Builder{}

	if m.user.user.Name == "" {
		// Doesnt let thorugh, if name is invalid
		if m.err == "" {
			m.ui.WriteString("Name: \n")
		} else {
			m.ui.WriteString(m.err + "Try again: \n")
		}
	}
	// Listens for input
	m.ui.WriteString(m.input.View())

	m.ui.WriteString("\n\nPress Esc to quit.\n")
}

// Displays the chatroom screen
func (m *model) displayChat() {
	m.ui = strings.Builder{}

	m.lock.RLock()

	defer m.lock.RUnlock()

	m.logUserMessages()

	m.listenInput("\n\nPress CtrlZ to see helpful commands\n")
}

// Displays the settings screen
func (m *model) displaySettings() {
	m.ui = strings.Builder{}

	m.ui.WriteString("Options: \n")

	for _, option := range settingOptions {
		m.ui.WriteString(option.name + ": " + option.value + "\n")
	}

	m.listenInput("\n\nPress CtrlZ to see helpful commands\n")
}

// Displays the users screen
func (m *model) displayUsers() {
	m.ui = strings.Builder{}

	m.ui.WriteString("Current Users: \n")

	for _, user := range m.userList.users {
		m.ui.WriteString(user.Name + ": " + user.Status + "\n")
	}

	m.listenInput("\n\nPress CtrlZ to see helpful commands\n")
}

// Displays the quit screen
func (m *model) displayQuit() {
	m.ui = strings.Builder{}

	if m.user.user.Status == "offline" {
		m.ui.WriteString("Logging Out! Be sure to come back! \n\n")
		return
	}
	m.ui.WriteString("Are you sure you want to log out? (Y/N)\n")

	m.listenInput("")
}

// Displayss the help screen
func (m *model) displayHelp() {
	m.ui = strings.Builder{}

	m.ui.WriteString("Available Commands: \n")

	for _, option := range helpOptions {
		m.ui.WriteString(option.name + ": " + option.value + "\n")
	}
	m.ui.WriteString("\n\nPress CtrlH to return to the chatroom\n")
}

// Displays all user messages
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

// Listens for client input
func (m *model) listenInput(action ...string) {
	m.ui.WriteString("<-------------------------------------> \nMessage: \n")
	// Listens for input
	m.ui.WriteString(m.input.View())

	if action[0] != "" {
		m.ui.WriteString(action[0])
	}
}
