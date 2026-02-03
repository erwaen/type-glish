package states

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/ui"
)

type ResultState struct{}

func (s *ResultState) Init(ctx *game.Context) tea.Cmd {
	return nil
}

func (s *ResultState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// Update Context Stats or History here if needed
			// Update current narrative for next turn
			ctx.CurrentNarrative = ctx.LastAssessment.OutcomeDescription

			// For now, loop back to Narrative with the outcome
			nextNarrative := &NarrativeState{
				Content: ctx.CurrentNarrative,
			}
			return nextNarrative, nil
		}
	}
	return s, nil
}

func (s *ResultState) View(ctx *game.Context) string {
	a := ctx.LastAssessment

	scoreColor := ui.ColorError
	if a.GrammarScore > 7 {
		scoreColor = ui.ColorSuccess
	}

	scoreStyle := lipgloss.NewStyle().Foreground(scoreColor).Bold(true)
	scoreText := scoreStyle.Render(fmt.Sprintf("%d/10", a.GrammarScore))

	content := a.OutcomeDescription +
		"\n\n" +
		fmt.Sprintf("Score: %s  |  Damage: %d", scoreText, a.DamageDealt) +
		"\n" +
		a.DMComment +
		"\n\nPress [Enter] to continue..."

	return ui.CenteredView("ASSESSMENT", content, true, ctx.Width, ctx.Height)
}
