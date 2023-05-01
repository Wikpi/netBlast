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
	Color       string    `json:"color"`
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

	return &model{
		input:     ti,
		userColor: colors[rand.Intn(len(colors))],
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

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

// Append messages received from the websocket connection
func (m *model) addNewMessages() {
	for {
		nMessage := struct {
			Username    string
			Message     string
			MessageTime time.Time
			Color       string
		}{}

		err := wsjson.Read(context.Background(), m.conn, &nMessage)
		handleError("Client/add: couldnt read body.", err)

		msg := message{
			Username:    nMessage.Username,
			Message:     nMessage.Message,
			MessageTime: nMessage.MessageTime,
			Color:       nMessage.Color,
		}

		m.messages = append(m.messages, msg)
	}
}

// Handles user input
func (m *model) handleMessage() {
	value := m.input.Value()
	if value == "" {
		return
	}

	if m.name == "" {
		m.registerUser(value)
	} else {
		m.writeMessage(value)
	}
}

// Registers and establishes websocket connection with the server
func (m *model) registerUser(value string) {
	data := struct {
		Name string `json:"name"`
	}{Name: value}

	jData, err := json.Marshal(data)
	handleError("Client/register: couldnt parse to json.", err)

	req, err := http.NewRequest(http.MethodPost, "http://"+reqURL+"/register", bytes.NewBuffer(jData))
	handleError("Client/register: couldnt create request.", err)

	res, err := http.DefaultClient.Do(req)
	handleError("Client/register: couldnt send an http request.", err)

	if res.StatusCode == http.StatusAccepted {
		m.name = value

		c, _, err := websocket.Dial(context.Background(), "ws://"+reqURL+"/message", nil)
		handleError("Client/register: couldnt connect websocket.", err)
		m.conn = c

		go m.addNewMessages()

	} else {
		resBody, err := ioutil.ReadAll(res.Body)
		handleError("Client/register: couldnt read response body.", err)

		err = json.Unmarshal(resBody, &m.err)
		handleError("Client/register: couldnt parse json.", err)
	}
}

// Writes user message to websocket connection
func (m model) writeMessage(value string) {
	message := message{
		Username:    m.name,
		Message:     value,
		MessageTime: time.Now(),
		Color:       m.userColor,
	}

	err := wsjson.Write(context.Background(), m.conn, message)
	handleError("Client/message: couldnt send message.", err)
}

func (m model) View() string {
	var s strings.Builder

	m.displayUserMessages(&s)

	// Listens for input
	s.WriteString(m.input.View())
	s.WriteString("\nPress Esc to quit.\n")

	return s.String()
}

// Displays messages
func (m model) displayUserMessages(s *strings.Builder) {
	if m.name == "" {
		// Doesnt let thorugh, if name is invalid
		if m.err == "" {
			s.WriteString("Name: \n")
		} else {
			s.WriteString(m.err + "Try again: \n")
		}
	} else {
		// Displays previous messages
		for _, msg := range m.messages {
			// Username color styler
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(msg.Color))

			// Displays message time (dd/mm/yyyy) if it was sent on a different day
			if msg.MessageTime.Day() < time.Now().Day() {
				s.WriteString(msg.MessageTime.Format("01-02-2006") + "\n\n")
			}
			msgMin := ""

			if msg.MessageTime.Minute() < 10 {
				msgMin = "0"
			}

			s.WriteString(strconv.Itoa(msg.MessageTime.Hour()) + ":" + msgMin + strconv.Itoa(msg.MessageTime.Minute()) + " " + style.Render(msg.Username) + ": " + msg.Message + "\n")
		}

		s.WriteString("<-------------------------------------> \n Message: \n")
	}
}

// Handles incoming error
func handleError(errMsg string, pErr error) {
	if pErr != nil {
		file, err := os.OpenFile("client/logs.txt", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Print(err)
		}
		defer file.Close()

		// Writes error to logs file
		if _, err := file.WriteString(pErr.Error()); err != nil {
			fmt.Println(err)
		}

		// Exits program and gives message where error occured
		log.Fatal(errMsg)
	}

}
