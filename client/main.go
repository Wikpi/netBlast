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

var colors = []string{"#000000", "#FFFFFF", "#800080", "#00FF00", "#FFA500", "#FF0000", "#FF00FF", "#00FFFF", "#000080"}

type model struct {
	input     textinput.Model
	cursor    int
	name      string
	userColor string
	messages  []message
}

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
			return m, tea.Quit
		case tea.KeyEnter:
			value := m.input.Value()

			if value == "" {
				return m, nil
			}

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

			m.input.SetValue("")

			return m, nil
		}
	}
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m model) View() string {
	var s strings.Builder
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(m.userColor))
	if m.name == "" {
		s.WriteString("Name: \n")
	} else {
		for _, m := range m.messages {
			if m.messageTime.Day() < time.Now().Day() {
				s.WriteString(m.messageTime.Format("01-02-2006") + "\n\n")
			}

			s.WriteString(strconv.Itoa(m.messageTime.Hour()) + ":" + strconv.Itoa(m.messageTime.Minute()) + " ")
			s.WriteString(style.Render(m.username) + ": " + m.message + "\n")
		}

		s.WriteString("Message: \n")
	}

	s.WriteString(m.input.View())
	s.WriteString("\nPress Esc to quit.\n")

	return s.String()
}
