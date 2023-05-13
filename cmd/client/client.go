package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"netBlast/pkg"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Structure that stores application state
type model struct {
	input     textinput.Model
	conn      *websocket.Conn
	cursor    int
	name      string
	userColor string
	err       string
	messages  []pkg.Message
	lock      sync.Mutex
	UI        strings.Builder
}

func main() {
	rand.Seed(time.Now().UnixNano())

	p := tea.NewProgram(createModel())

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
		case tea.KeyCtrlC, tea.KeyEsc:
			m.conn.Close(websocket.StatusNormalClosure, "Connection Closed")
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

// Registers and establishes websocket connection with the server
func (m *model) registerNewUser(value string) {
	data := struct {
		Name string `json:"name"`
	}{Name: value}

	jData, err := json.Marshal(data)
	handleError("Client/register: couldnt parse to json.", err)

	req, err := http.NewRequest(http.MethodPost, "http://"+pkg.ServerURL+"/register", bytes.NewBuffer(jData))
	handleError("Client/register: couldnt create request.", err)

	res, err := http.DefaultClient.Do(req)
	handleError("Client/register: couldnt receive an http request.", err)

	if res.StatusCode == http.StatusAccepted {
		m.name = value

		c, _, err := websocket.Dial(context.Background(), "ws://"+pkg.ServerURL+"/message", nil)
		handleError("Client/register: couldnt connect websocket.", err)
		m.conn = c

		go m.receiveNewMessages()
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	handleError("Client/register: couldnt read response body.", err)

	err = json.Unmarshal(resBody, &m.err)
	handleError("Client/register: couldnt parse json.", err)
}

// Stores messages received from the websocket connection
func (m *model) receiveNewMessages() {
	for {
		message := struct {
			Username    string
			Message     string
			MessageTime time.Time
			Color       string
		}{}

		err := wsjson.Read(context.Background(), m.conn, &message)
		handleError("Client/add: couldnt read body.", err)

		msg := pkg.Message{
			Username:    message.Username,
			Message:     message.Message,
			MessageTime: message.MessageTime,
			Color:       message.Color,
		}

		m.messages = append(m.messages, msg)
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

	err := wsjson.Write(context.Background(), m.conn, message)
	handleError("Client/message: couldnt send message.", err)
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
	m.displayUserMessages()

	m.UI.WriteString("<-------------------------------------> \n Message: \n")
}

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

// Picks one random color from the scrapped color list
func getColor() string {
	body, err := ioutil.ReadFile("./tools/scrapper/colors.txt")
	handleError("Client/autolycus: couldnt open file.", err)

	colors := strings.Split(string(body), ", ")

	return colors[rand.Intn(len(colors))]
}

// Handles incoming error
func handleError(errMsg string, incomingErr error) {
	if incomingErr == nil {
		return
	}

	file, err := os.OpenFile("./logs/client/logs.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Print(err)
	}
	defer file.Close()

	// Writes error to logs file
	if _, err := file.WriteString(time.Now().Format("2006-01-02 15:04") + " " + incomingErr.Error() + "\n\n"); err != nil {
		fmt.Println(err)
	}

	// Exits program and gives message where error occured
	log.Fatal(errMsg)
}
