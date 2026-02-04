package states

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/ui"
)

type CombatResultState struct{}

func (s *CombatResultState) Init(ctx *game.Context) tea.Cmd {
	return nil
}

func (s *CombatResultState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// Check if player is dead
			if ctx.Stats.HP <= 0 {
				return &GameOverState{}, nil
			}

			// Check if enemy is dead
			if ctx.CurrentEnemy != nil && ctx.CurrentEnemy.HP <= 0 {
				return &VictoryState{}, nil
			}

			// Continue combat
			ctx.CurrentNarrative = ctx.CombatAssessment.Outcome
			return NewCombatState(), nil
		}
		if msg.Type == tea.KeyCtrlC {
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s *CombatResultState) View(ctx *game.Context) string {
	a := ctx.CombatAssessment

	var content string

	// Show what the player typed
	content += ui.StyleSubTitle.Render("YOU SAID:") + "\n"
	content += fmt.Sprintf("> %s\n\n", ctx.LastInput)

	// Show corrected version if different
	if a.CorrectedSentence != ctx.LastInput {
		content += ui.StyleSubTitle.Render("CORRECTED:") + "\n"
		content += fmt.Sprintf("> %s\n\n", a.CorrectedSentence)
	}

	// Outcome
	content += "───────────────────────────────────────────\n\n"
	content += ui.StyleSubTitle.Render("RESULT:") + "\n"
	content += a.Outcome + "\n\n"

	// Damage dealt/received
	scoreColor := ui.ColorError
	if a.GrammarScore >= 7 {
		scoreColor = ui.ColorSuccess
	} else if a.GrammarScore >= 5 {
		scoreColor = ui.ColorWarning
	}

	scoreStyle := lipgloss.NewStyle().Foreground(scoreColor).Bold(true)
	damageDealtStyle := lipgloss.NewStyle().Foreground(ui.ColorSuccess).Bold(true)
	damageReceivedStyle := lipgloss.NewStyle().Foreground(ui.ColorError).Bold(true)

	content += fmt.Sprintf("Grammar Score: %s\n", scoreStyle.Render(fmt.Sprintf("%d/10", a.GrammarScore)))
	content += fmt.Sprintf("Damage Dealt: %s    Damage Received: %s\n\n",
		damageDealtStyle.Render(fmt.Sprintf("-%d", a.DamageDealt)),
		damageReceivedStyle.Render(fmt.Sprintf("-%d", a.DamageReceived)))

	// DM Comment
	content += ui.StyleSubTitle.Render("DM:") + " " + a.DMComment + "\n\n"

	// Current HP status
	if ctx.CurrentEnemy != nil {
		content += ui.RenderHPBar(ctx.CurrentEnemy.HP, ctx.CurrentEnemy.MaxHP, ctx.CurrentEnemy.Name, 15) + "\n"
	}
	content += ui.RenderHPBar(ctx.Stats.HP, 100, "You", 15) + "\n\n"

	// Show error if any (muted grey)
	if ctx.LastError != "" {
		errorStyle := lipgloss.NewStyle().Foreground(ui.ColorSubtext).Italic(true)
		content += errorStyle.Render("(LLM error: "+ctx.LastError+")") + "\n\n"
	}

	content += ui.StyleHelp.Render("Press [Enter] to continue...")

	return ui.CenteredView("⚔ COMBAT RESULT ⚔", content, true, ctx.Width, ctx.Height)
}

// GameOverState handles player death
type GameOverState struct{}

func (s *GameOverState) Init(ctx *game.Context) tea.Cmd {
	return nil
}

func (s *GameOverState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" || msg.String() == "q" {
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s *GameOverState) View(ctx *game.Context) string {
	content := `
    ╔═══════════════════════════════════════╗
    ║                                       ║
    ║            G A M E   O V E R          ║
    ║                                       ║
    ║   Your grammar failed you...          ║
    ║   The Kingdom of Lexicon mourns.      ║
    ║                                       ║
    ╚═══════════════════════════════════════╝

    Press [Enter] or [Q] to exit.
`
	return ui.CenteredView("💀 DEFEAT 💀", content, true, ctx.Width, ctx.Height)
}
