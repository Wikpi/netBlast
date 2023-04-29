package main

import (
	"log"
	"math/rand"
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
	messages  []message
}

// Structure of individual user message
type message struct {
	username    string
	message     string
	messageTime time.Time
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

	checkMessageType(m, msg, &cmd)

	// Updates input
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

// Checks what input was entered
func checkMessageType(m *model, msg tea.Msg, cmd *tea.Cmd) {
	switch key := msg.(type) {
	case tea.KeyMsg:
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			*cmd = tea.Quit
		case tea.KeyEnter:
			handleMessage(m, cmd)

			*cmd = nil
		}
	}
}

// Sets entered message into respective field
func handleMessage(m *model, cmd *tea.Cmd) {
	value := m.input.Value()
	if value == "" {
		*cmd = nil
	}

	// Registers user if it hasnt already
	if m.name == "" {
		m.name = value
	} else {
		message := message{
			username:    m.name,
			message:     value,
			messageTime: time.Now(),
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
		s.WriteString("Name: \n")
	} else {
		// Displays previous messages
		for _, m := range m.messages {
			// Displays message time (dd/mm/yyyy) if it was sent on a different day
			if m.messageTime.Day() < time.Now().Day() {
				s.WriteString(m.messageTime.Format("01-02-2006") + "\n\n")
			}

			s.WriteString(strconv.Itoa(m.messageTime.Hour()) + ":" + strconv.Itoa(m.messageTime.Minute()) + " ")
			s.WriteString(style.Render(m.username) + ": " + m.message + "\n")
		}

		s.WriteString("Message: \n")
	}
}
