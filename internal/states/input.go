package states

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/ui"
)

type InputState struct {
	textInput textinput.Model
}

func (s *InputState) Init(ctx *game.Context) tea.Cmd {
	ti := textinput.New()
	ti.Placeholder = "Write here what you are gonna do?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	s.textInput = ti
	return textinput.Blink
}

func (s *InputState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			ctx.LastInput = s.textInput.Value()
			return &ProcessingState{}, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return s, tea.Quit
		}
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s *InputState) View(ctx *game.Context) string {
	return ui.Box("YOUR ACTION",
		"Describe your action in English:\n\n"+
			s.textInput.View(),
	)
}
