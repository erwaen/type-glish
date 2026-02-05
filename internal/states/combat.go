package states

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/ui"
)

type CombatState struct {
	textInput textinput.Model
}

func NewCombatState() *CombatState {
	ti := textinput.New()
	ti.Placeholder = "Describe your attack..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	return &CombatState{
		textInput: ti,
	}
}

func (s *CombatState) Init(ctx *game.Context) tea.Cmd {
	return textinput.Blink
}

func (s *CombatState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if s.textInput.Value() != "" {
				ctx.LastInput = s.textInput.Value()
				return &CombatProcessingState{}, nil
			}
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyEsc:
			return NewMenuState(nil), nil
		}
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s *CombatState) View(ctx *game.Context) string {
	if ctx.CurrentEnemy == nil {
		return ui.CenteredView("ERROR", "No enemy found!", true, ctx.Width, ctx.Height)
	}

	enemy := ctx.CurrentEnemy

	// Build combat view
	var content string

	// Status bar at top
	content += ui.RenderStatusBar(ctx.Stats.HP, 100, ctx.Stats.Gold, ctx.Stats.XP) + "\n\n"

	// Header: Location and Enemy
	content += ui.RenderCombatHeader(enemy.Location, enemy.Name) + "\n\n"

	// Enemy HP Bar
	content += ui.RenderHPBar(enemy.HP, enemy.MaxHP, enemy.Name, 20) + "\n\n"

	// DM Description
	content += ui.StyleSubTitle.Render("DM: "+enemy.Description) + "\n\n"

	// Narrative context if any
	if ctx.CurrentNarrative != "" {
		content += ctx.CurrentNarrative + "\n\n"
	}

	// Divider
	content += "───────────────────────────────────────────\n\n"

	// Input
	content += "YOUR ACTION:\n"
	content += s.textInput.View() + "\n\n"

	content += ui.StyleHelp.Render("(Type your combat action and press Enter)")

	return ui.CenteredView("COMBAT", content, true, ctx.Width, ctx.Height)
}
