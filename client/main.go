package main

import (
	"bytes"
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
)

// Colors for usernames
var colors = []string{"#000000", "#FFFFFF", "#800080", "#00FF00", "#FFA500", "#FF0000", "#FF00FF", "#00FFFF", "#000080"}

type model struct {
	input     textinput.Model
	cursor    int
	name      string
	userColor string
	err       string
	messages  []message
}

// Structure of individual user message
type message struct {
	Username    string
	Message     string
	MessageTime time.Time
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

	// Checks what input was entered
	switch key := msg.(type) {
	case tea.KeyMsg:
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			handleMessage(m)

			return m, nil
		}
	}

	// Updates input
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

// Sets entered message into respective field
func handleMessage(m *model) {
	value := m.input.Value()
	if value == "" {
		return
	}

	// Registers user if it hasnt already
	if m.name == "" {
		data := value
		jData, err := json.Marshal(data)
		if err != nil {
			log.Fatal("Client/register: couldnt parse to json: ", err)
		}
		reqURL := "http://localhost:8080/register"

		req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(jData))
		if err != nil {
			log.Fatal("Client/register: couldnt create request: ", err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal("Client/register: couldnt send an http request: ", err)
		}

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal("Client/register: couldnt read response body: ", err)
		}
		if res.StatusCode == http.StatusAccepted {
			m.name = value
		} else {
			err := json.Unmarshal(resBody, &m.err)
			if err != nil {
				log.Fatal("Client/register: couldnt parse json: ", err)
			}
		}
	} else {
		message := message{
			Username:    m.name,
			Message:     value,
			MessageTime: time.Now(),
		}

		m.messages = append(m.messages, message)
	}

	// Resets input field
	m.input.SetValue("")
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
