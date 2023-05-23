package client

import (
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Routes to different UI screens
func (m *model) routeScreen() {
	switch m.screen {
	case registerScreen:
		m.displayRegister()
	case chatScreen:
		m.displayChat()
	case settingsScreen:
		m.displaySettings()
	case usersScreen:
		m.displayUsers()
	case quitScreen:
		m.displayQuit()
	case helpScreen:
		m.displayHelp()
	}
}

// Displays the register screen
func (m *model) displayRegister() {
	m.ui = strings.Builder{}

	if m.user.Name == "" {
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

	m.listenInput()
	m.ui.WriteString("\n\nPress CtrlH to see helpful commands\n")
}

// Displays the settings screen
func (m *model) displaySettings() {
	m.ui = strings.Builder{}

	m.ui.WriteString("Options: \n")

	for _, option := range m.settings.options {
		m.ui.WriteString("\t" + option.name + ": " + option.value + "\n")
	}

	m.listenInput()
	m.ui.WriteString("\n\nPress CtrlH to see helpful commands\n")
}

// Displays the users screen
func (m *model) displayUsers() {
	m.ui = strings.Builder{}

	m.ui.WriteString("Current Users: \n")

	for idx, user := range m.userList.users {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(user.UserColor))

		m.ui.WriteString("\t" + strconv.Itoa(idx+1) + ". " + style.Render(user.Name) + ": " + user.Status + "\n")
	}

	m.listenInput()
	m.ui.WriteString("\n\nPress CtrlH to see helpful commands\n")
}

// Displays the quit screen
func (m *model) displayQuit() {
	m.ui = strings.Builder{}

	if m.user.Status == "offline" {
		m.ui.WriteString("Logging Out! Be sure to come back! \n\n")
		return
	}
	m.ui.WriteString("Are you sure you want to log out? (Y/N)\n")

	m.listenInput()
}

// Displayss the help screen
func (m *model) displayHelp() {
	m.ui = strings.Builder{}

	m.ui.WriteString("Available Commands: \n")

	for _, option := range m.help.options {
		m.ui.WriteString(option.name + ": " + option.value + "\n")
	}
	m.ui.WriteString("\n\nPress CtrlH to return to " + m.prevScreen + " screen.\n")
}

// Displays all user messages
func (m *model) logUserMessages() {
	var currentTime time.Time

	// Displays previous messages
	for _, msg := range m.chat.messages {
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
func (m *model) listenInput() {
	m.ui.WriteString("<-------------------------------------> \nMessage: \n")
	// Listens for input
	m.ui.WriteString(m.input.View())
}
