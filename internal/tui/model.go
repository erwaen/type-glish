package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/llm"
	"math/rand/v2"
	"github.com/erwaen/type-glish/internal/targets"
	"strings"
)

type WordMsg string
type ErrMsg error

type Model struct {
	Cursor  int
	Target  string
	Typed   string
	Loading bool
}

func InitialModel() Model {

	randomIndex:= rand.IntN(len(targets.StartupExamples)	)
	randomTarget:= targets.StartupExamples[randomIndex]
	return Model{
		Target: randomTarget,
		Typed:  "",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case WordMsg:
        m.Target = string(msg)
        m.Typed = ""
        m.Cursor = 0
        m.Loading = false
        return m, nil

    case ErrMsg:
        m.Target = "Error: " + msg.Error()
        m.Loading = false
        return m, nil

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        }
				

				// for block inputs when loading
        if m.Loading {
            return m, nil
        }

        // FINISHED STATE (Waiting for user to restart)
        // If cursor is at the end, ANY key triggers the new word
        if m.Cursor >= len(m.Target) {
            m.Loading = true
            return m, fetchNewWordCmd
        }

        // TYPING STATE
        switch msg.String() {
        case "backspace":
            if m.Cursor > 0 {
                m.Cursor--
                m.Typed = m.Typed[:len(m.Typed)-1]
            }

        default:
            expected := string(m.Target[m.Cursor])
            if msg.String() == expected {
                m.Typed += msg.String()
                m.Cursor++
            }
        }
    }

    return m, nil
}

func (m Model) View() string {
    if m.Loading {
        return "\n  Thinking...\n"
    }

    // Show the "Press any key" prompt when finished
    if m.Cursor >= len(m.Target) {
        return fmt.Sprintf("\n  %s\n\n  %s\n", 
            TypedStyle.Render(m.Target), // Show full word in green
            CursorStyle.Render("Press any key for the next word..."),
        )
    }

    s := "Type to learn!\n\n"
    
    // Slice safety check
    rest := ""
    if m.Cursor+1 < len(m.Target) {
        rest = m.Target[m.Cursor+1:]
    }

    completed := TypedStyle.Render(m.Target[:m.Cursor])
    cursorChar := CursorStyle.Render(string(m.Target[m.Cursor]))
    remaining := GreyStyle.Render(rest)

    s += fmt.Sprintf("%s%s%s\n", completed, cursorChar, remaining)
		s += fmt.Sprintf("%s\n",strings.Repeat(" ",m.Cursor)+"^")
    s += "\nPress ctrl+c to quit.\n"

    return s
}
func fetchNewWordCmd() tea.Msg {
	word, err := llm.NewRandmoWord()
	if err != nil {
		return ErrMsg(err)
	}
	return WordMsg(word)
}
