package client

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"netBlast/pkg"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"nhooyr.io/websocket"
)

// Stores the application state
type model struct {
	user userInfo

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

// Stores user info
type userInfo struct {
	user     pkg.User
	messages []pkg.Message
}

// Creates the initial model that holds default values
func newClient() *model {
	model := &model{
		input:  textinput.New(),
		screen: "register",
		user:   userInfo{},
	}

	model.input.Placeholder = "Your text"
	model.input.CharLimit = 256
	model.input.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
	model.input.Width = 30
	model.style = lipgloss.NewStyle().Bold(true)

	model.input.Focus()

	model.user.user.UserColor = getColor()

	model.help.options = []option{
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

	model.settings.options = []option{
		{
			name:  "Color",
			value: "-color, to update your user color.",
		},
		{
			name:  "Name",
			value: "change username. (Not finished)",
		},
	}

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
			if m.screen != "quit" {
				m.prevScreen = m.screen
				m.screen = "quit"
			}

		case tea.KeyCtrlH:
			changeScreen(m, "help")
		case tea.KeyCtrlX:
			changeScreen(m, "settings")
		case tea.KeyCtrlC:
			getUserList(m)
			changeScreen(m, "users")
		case tea.KeyEnter:
			m.routeMessage()

			// Resets input field
			m.input.SetValue("")

			if m.user.user.Status == "offline" {
				if m.user.user.Conn != nil {
					m.user.user.Conn.Close(websocket.StatusNormalClosure, "Connection Closed")
				}

				return m, tea.Quit
			}

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

	return m.ui.String()
}
