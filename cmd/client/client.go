package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"netBlast/pkg"
	"netBlast/tools/scrapper"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"nhooyr.io/websocket"
)

// Structure that stores the application state
type model struct {
	input     textinput.Model
	conn      *websocket.Conn
	cursor    int
	name      string
	userColor string
	err       string
	messages  []pkg.Message
	lock      sync.RWMutex
	UI        strings.Builder
}

// Creates the initial model that holds default values
func createModel() *model {
	ti := textinput.New()
	ti.Placeholder = "Your text"
	ti.CharLimit = 256
	ti.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
	ti.Width = 30

	ti.Focus()

	return &model{
		input:     ti,
		userColor: getColor(),
	}
}

// Starts a new client UI
func newClient() {
	rand.Seed(time.Now().UnixNano())

	p := tea.NewProgram(createModel())

	if _, err := p.Run(); err != nil {
		log.Fatal("Err: ", err)
	}
}

func Client() {
	newClient()
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
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.conn != nil {
				m.conn.Close(websocket.StatusNormalClosure, "Connection Closed")
			}
			return m, tea.Quit
		case tea.KeyEnter:
			m.handleMessage()

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
	//m. UI = ""

	m.displayUI()

	// Listens for input
	m.UI.WriteString(m.input.View())
	m.UI.WriteString("\nPress Esc to quit.\n")

	return m.UI.String()
}

/* ----------------Additional Functions---------------- */

// Handles user input
func (m *model) handleMessage() {
	value := m.input.Value()
	if value == "" {
		return
	}

	if m.name == "" {
		m.registerNewUser(value)
		return
	}
	m.writeNewMessage(value)
}

// Registers and establishes a websocket connection with the server
func (m *model) registerNewUser(value string) {
	name := pkg.Name{Name: value}

	data := pkg.ParseToJson(name, pkg.ClRegister+pkg.BadParse)

	res := handleHTTPRequest(data, "http://"+pkg.ServerURL+"/register")

	if res.StatusCode == http.StatusAccepted {
		m.name = value

		c, _, err := websocket.Dial(context.Background(), "ws://"+pkg.ServerURL+"/message", nil)
		pkg.HandleError(pkg.ClRegister+pkg.BadConn, err, 0)
		m.conn = c

		go m.receiveNewMessages()
		return
	}

	// Gives an error if registration failed
	resBody, err := ioutil.ReadAll(res.Body)
	pkg.HandleError(pkg.ClRegister+pkg.BadRead, err, 0)

	pkg.ParseFromJson(resBody, &m.err, pkg.ClRegister+pkg.BadParse)
}

// Stores messages received from the websocket connection
func (m *model) receiveNewMessages() {
	for {
		msg := pkg.WsRead(m.conn, pkg.ClMessage+pkg.BadRead)

		m.lock.Lock()
		m.messages = append(m.messages, msg)
		m.lock.Unlock()
	}
}

// Writes user message to websocket connection
func (m *model) writeNewMessage(value string) {
	message := pkg.Message{
		Username:    m.name,
		Message:     value,
		MessageTime: time.Now(),
		Color:       m.userColor,
	}

	pkg.WsWrite(m.conn, message, pkg.ClMessage+pkg.BadWrite)
}

// Adds neccessary info to display UI
func (m *model) displayUI() {
	m.UI = strings.Builder{}

	if m.name == "" {
		// Doesnt let thorugh, if name is invalid
		if m.err == "" {
			m.UI.WriteString("Name: \n")
			return
		}
		m.UI.WriteString(m.err + "Try again: \n")
		return
	}
	m.lock.RLock()
	defer m.lock.RUnlock()
	m.displayUserMessages()

	m.UI.WriteString("<-------------------------------------> \n Message: \n")
}

// Display all user messages
func (m *model) displayUserMessages() {
	var currentTime time.Time

	// Displays previous messages
	for _, msg := range m.messages {
		// Username color styler
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(msg.Color))

		// Dsiplays new time
		if currentTime.After(msg.MessageTime) {
			currentTime = msg.MessageTime

			m.UI.WriteString(currentTime.Format("15:04") + "\n")
		}

		m.UI.WriteString(style.Render(msg.Username) + ": " + msg.Message + "\n")
	}
}

/* ----------------Standalone Functions---------------- */

// Handles the POST request to server
func handleHTTPRequest(data []byte, URL string) *http.Response {
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(data))
	pkg.HandleError(pkg.ClRegister+pkg.BadReq, err, 0)

	res, err := http.DefaultClient.Do(req)
	pkg.HandleError(pkg.ClRegister+pkg.BadRes, err, 0)

	return res
}

// Picks one random color from the scrapped color list
func getColor() string {
	body, err := ioutil.ReadFile(pkg.Scrapper + "/colors.txt")
	pkg.HandleError(pkg.Cl+pkg.BadOpen, err, 0)

	colors := strings.Split(string(body), ", ")

	return colors[rand.Intn(len(colors))]
}

// Scrapes colors using Autolycus module
func useAutolycus() {
	scrapper.Scrape()
}
