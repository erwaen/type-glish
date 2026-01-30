package tui

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor int
	target string
	typed  string
}

func InitialModel() model {
	return model{
		target: "Hey, do you want to eat pizza today?",
		typed:  "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "backspace":
			if m.cursor > 0 {
				m.cursor--
				m.typed = m.typed[:len(m.typed)-1]
			}

		default:
			if m.cursor >= len(m.target) {
				return m, nil
			}
			expected := string(m.target[m.cursor])
			if msg.String() == expected {
				m.typed += msg.String()
				m.cursor++
			}

		}
	}

	return m, nil
}

func (m model) View() string {
	if m.cursor >= len(m.target) {
		return "You finished! Press q to quit.\n"
	}
	s := "Type to learn!\n\n"

	completed := TypedStyle.Render(m.target[:m.cursor])
	cursorChar := CursorStyle.Render(string(m.target[m.cursor]))
	remaining := GreyStyle.Render(m.target[m.cursor+1:])

	s += fmt.Sprintf("%s[%s]%s\n", completed, cursorChar, remaining)

	s += "\nPress q to quit.\n"

	return s
}
