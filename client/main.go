package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Colors for usernames
var colors = []string{"#000000", "#FFFFFF", "#800080", "#00FF00", "#FFA500", "#FF0000", "#FF00FF", "#00FFFF", "#000080"}

// Server url and port
const reqURL = "localhost:8080"

type model struct {
	input     textinput.Model
	conn      *websocket.Conn
	ctx       context.Context
	cursor    int
	name      string
	userColor string
	err       string
	messages  []message
}

// Structure of individual user message
type message struct {
	Username    string    `json:"username"`
	Message     string    `json:"message"`
	MessageTime time.Time `json:"messageTime"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(createModel())

	if _, err := p.Run(); err != nil {
		log.Fatal("Err: ", err)
	}
}

// Intial model
func createModel() *model {
	ti := textinput.New()
	ti.Placeholder = "Your text"
	ti.CharLimit = 256
	ti.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
	ti.Width = 30

	ti.Focus()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return &model{
		input:     ti,
		userColor: colors[rand.Intn(len(colors))],
		ctx:       ctx,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Checks what input was entered
	switch key := msg.(type) {
	case tea.KeyMsg:
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.conn.Close(websocket.StatusNormalClosure, "")
			return m, tea.Quit
		case tea.KeyEnter:
			m.handleMessage()

			return m, nil
		}
	}
	// if m.name != "" {
	// 	go m.addNewMessages()
	// }

	// Updates input
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

// Append messages, so they could be shown to the user
func (m *model) addNewMessages() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var nMessage struct {
			Name        string
			Message     string
			MessageTime time.Time
		}

		err := wsjson.Read(ctx, m.conn, &nMessage)
		checkError("Client/add: couldnt read body: ", err)

		msg := message{
			Username:    nMessage.Name,
			Message:     nMessage.Message,
			MessageTime: nMessage.MessageTime,
		}

		m.messages = append(m.messages, msg)
	}
}

// Checks if there is an error and exist
func checkError(errMsg string, err error) {
	if err != nil {
		log.Fatal(errMsg, err)
	}
}

// Sets entered message into respective field
func (m *model) handleMessage() {
	value := m.input.Value()
	if value == "" {
		return
	}

	if m.name == "" {
		m.registerUser(value)
	} else {
		message := message{
			Username:    m.name,
			Message:     value,
			MessageTime: time.Now(),
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := wsjson.Write(ctx, m.conn, message)
		checkError("Client/message: couldnt send message: ", err)
	}

	// Resets input field
	m.input.SetValue("")
}

// Registers user if it hasnt already
func (m *model) registerUser(value string) {
	data := struct {
		Name string `json:"name"`
	}{Name: value}

	jData, err := json.Marshal(data)
	checkError("Client/register: couldnt parse to json: ", err)

	req, err := http.NewRequest(http.MethodPost, "http://"+reqURL+"/register", bytes.NewBuffer(jData))
	checkError("Client/register: couldnt create request: ", err)

	res, err := http.DefaultClient.Do(req)
	checkError("Client/register: couldnt send an http request: ", err)

	if res.StatusCode == http.StatusAccepted {
		m.name = value

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		c, _, err := websocket.Dial(ctx, "ws://"+reqURL+"/message", nil)
		checkError("Client/register: couldnt connect websocket: ", err)
		m.conn = c

	} else {
		resBody, err := ioutil.ReadAll(res.Body)
		checkError("Client/register: couldnt read response body: ", err)

		err = json.Unmarshal(resBody, &m.err)
		checkError("Client/register: couldnt parse json: ", err)
	}
}

func (m model) View() string {
	var s strings.Builder

	displayUserMessages(m, &s)

	// Listens for input
	s.WriteString(m.input.View())
	s.WriteString("\nPress Esc to quit.\n")

	return s.String()
}

// Displays messages
func displayUserMessages(m model, s *strings.Builder) {
	// Username color styler
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(m.userColor))

	if m.name == "" {
		if m.err == "" {
			s.WriteString("Name: \n")
		} else {
			s.WriteString(m.err)
			s.WriteString("Try again: \n")
		}
	} else {
		// Displays previous messages
		for _, m := range m.messages {
			// Displays message time (dd/mm/yyyy) if it was sent on a different day
			if m.MessageTime.Day() < time.Now().Day() {
				s.WriteString(m.MessageTime.Format("01-02-2006") + "\n\n")
			}

			s.WriteString(strconv.Itoa(m.MessageTime.Hour()) + ":" + strconv.Itoa(m.MessageTime.Minute()) + " ")
			s.WriteString(style.Render(m.Username) + ": " + m.Message + "\n")
		}

		s.WriteString("Message: \n")
	}
}
