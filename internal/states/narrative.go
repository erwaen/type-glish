package states

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/ui"
)

type NarrativeState struct {
	Content string // The story text
}

func (s NarrativeState) Init(ctx *game.Context) tea.Cmd {
	return nil
}

func (s NarrativeState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// Transition to Input Mode
			return &InputState{}, nil
		}
	}
	return s, nil
}

func (s NarrativeState) View(ctx *game.Context) string {
	hint := ui.StyleHelp.Render("Press [Enter] to take action... (Ctrl+S for Settings)")

	return ui.CenteredView("DUNGEON MASTER", s.Content+"\n\n"+hint, true, ctx.Width, ctx.Height)
}
