package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	input    textinput.Model
	cursor   int
	name     string
	messages []message
}

type message struct {
	username string
	message  string
}

func main() {
	p := tea.NewProgram(createModel())

	if _, err := p.Run(); err != nil {
		log.Fatal("Err: ", err)
	}
}

func createModel() *model {
	ti := textinput.New()
	ti.Placeholder = "Your text"
	ti.CharLimit = 256
	// inp.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
	ti.Width = 30

	ti.Focus()

	return &model{
		input: ti,
		//randomColor: colors[rand.Intn(len(colors))],
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
					username: m.name,
					message:  value,
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

	if m.name == "" {
		s.WriteString("Name: \n")
	} else {
		for _, m := range m.messages {
			s.WriteString(m.username + ": " + m.message + "\n")
		}

		s.WriteString("Message: \n")
	}

	s.WriteString(m.input.View())
	s.WriteString("\nPress Esc to quit.\n")

	return s.String()
}
