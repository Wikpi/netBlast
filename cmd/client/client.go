package client

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"netBlast/pkg"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"nhooyr.io/websocket"
)

// Stores the application state
type model struct {
	cursor   int
	screen   string
	err      string
	user     userInfo
	settings settings
	userList userList
	input    textinput.Model
	lock     sync.RWMutex
	ui       strings.Builder
}

// Additional model for userlist screen
type userList struct {
	users []list.Item
	list  list.Model
}

// Additional model for settings screen
type settings struct {
}

// Stores user info
type userInfo struct {
	name      string
	userColor string
	messages  []pkg.Message
	conn      *websocket.Conn
}

// Creates the initial model that holds default values
func newClient() *model {
	ti := textinput.New()
	ti.Placeholder = "Your text"
	ti.CharLimit = 256
	ti.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
	ti.Width = 30

	ti.Focus()

	color := getColor()
	if color == "" {
		color = "#FFF"
	}

	items := []list.Item{item{title: "Bobby", desc: "broo"}}

	model := &model{
		input:    ti,
		screen:   "register",
		user:     userInfo{userColor: color},
		userList: userList{users: nil, list: list.New(items, list.NewDefaultDelegate(), 0, 0)},
	}
	model.userList.list.Title = "Current Users"

	return model
}

// Starts a new client
func Client() {
	if status := pingServer(); status != http.StatusOK {
		log.Fatal("Server not responding!")
	}

	rand.Seed(time.Now().UnixNano())

	p := tea.NewProgram(newClient())

	if _, err := p.Run(); err != nil {
		log.Fatal("Err: ", err)
	}
}

/* ----------------Main UI Functions---------------- */

// Returns an initial command for the application to run
func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

// Handles incoming events and updates the model accordingly
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch key := msg.(type) {
	case tea.KeyMsg:
		switch key.Type {
		case tea.KeyEsc:
			if m.user.conn != nil {
				m.user.conn.Close(websocket.StatusNormalClosure, "Connection Closed")
			}
			return m, tea.Quit
		case tea.KeyCtrlX:
			if m.screen == "register" {
				return m, nil
			}

			if m.screen == "chat" {
				m.screen = "settings"
				return m, nil
			}
			m.screen = "chat"
			return m, nil
		case tea.KeyCtrlC:
			if m.screen == "register" {
				return m, nil
			}

			if m.screen == "chat" {
				m.screen = "users"
				getUserList(m)
				cmd := m.userList.list.SetItems(m.userList.users)
				return m, cmd
			}
			m.screen = "chat"
			return m, nil
		case tea.KeyEnter:
			m.routeMessage()

			// Resets input field
			m.input.SetValue("")

			return m, nil
		}
	}

	// Updates input
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

// Renders the UI based on the data in the model
func (m *model) View() string {
	m.routeScreen()

	// Listens for input
	m.ui.WriteString(m.input.View())

	// if m.screen == "chat" {
	// 	m.ui.WriteString("\n\nPress CtrlX to enter settings.")
	// } else if m.screen == "settings" {
	// 	m.ui.WriteString("\n\n Press CtrlX to return to the chatroom.")
	// }

	m.ui.WriteString("\nPress Esc to quit.\n")

	return m.ui.String()
}
