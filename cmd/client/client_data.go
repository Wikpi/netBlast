package client

import (
	"netBlast/pkg"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	welcomeScreen  = "welcome"
	registerScreen = "register"
	loginScreen    = "login"
	chatScreen     = "chat"
	dmScreen       = "dms"
	settingsScreen = "settings"
	helpScreen     = "help"
	usersScreen    = "users"
	quitScreen     = "quit"
)

// Stores the application state
type model struct {
	user pkg.User

	chat     chat
	help     help
	settings settings
	userList userList

	cursor     int
	screen     string
	prevScreen string
	err        string
	style      lipgloss.Style
	input      textinput.Model
	lock       sync.RWMutex
	ui         strings.Builder
}

// Option chassis
type option struct {
	name  string
	value string
}

// Additional model for chat screen
type chat struct {
	messages []pkg.Message
	dms      []pkg.Message
}

// Additional model for userlist screen
type userList struct {
	users []pkg.User
	err   string
}

// Additional model for settings screen
type settings struct {
	options []option
	err     string
}

// Additional model for help screen
type help struct {
	options []option
}
