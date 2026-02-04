package states

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/llm"
	"github.com/erwaen/type-glish/internal/ui"

	spinner "github.com/charmbracelet/bubbles/spinner"
)

type CombatProcessingState struct {
	spinner spinner.Model
}

func (s *CombatProcessingState) Init(ctx *game.Context) tea.Cmd {
	s.spinner = spinner.New()
	s.spinner.Spinner = spinner.Dot
	s.spinner.Style = lipgloss.NewStyle().Foreground(ui.ColorPrimary)

	enemyName := "Unknown"
	location := "Unknown"
	if ctx.CurrentEnemy != nil {
		enemyName = ctx.CurrentEnemy.Name
		location = ctx.CurrentEnemy.Location
	}

	return tea.Batch(
		s.spinner.Tick,
		func() tea.Msg {
			return ctx.LLMClient.AnalyzeCombatAction(ctx.LastInput, enemyName, location)
		},
	)
}

func (s *CombatProcessingState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case llm.CombatAssessmentMsg:
		if msg.Err != nil {
			log.Printf("Error from LLM: %v", msg.Err)
			// Store error for display
			ctx.LastError = msg.Err.Error()
			// On error, create a default assessment
			ctx.CombatAssessment = llm.CombatAssessment{
				CorrectedSentence: ctx.LastInput,
				GrammarScore:      5,
				DamageDealt:       7,
				DamageReceived:    5,
				DMComment:         "The Dungeon Master is momentarily distracted...",
				Outcome:           "Your attack connects, but so does the enemy's!",
				IsRelevant:        true,
			}
		} else {
			ctx.LastError = "" // Clear any previous error
			ctx.CombatAssessment = msg.Data
		}

		// Apply damage
		if ctx.CurrentEnemy != nil {
			ctx.CurrentEnemy.HP -= ctx.CombatAssessment.DamageDealt
			if ctx.CurrentEnemy.HP < 0 {
				ctx.CurrentEnemy.HP = 0
			}
		}
		ctx.Stats.HP -= ctx.CombatAssessment.DamageReceived
		if ctx.Stats.HP < 0 {
			ctx.Stats.HP = 0
		}

		return &CombatResultState{}, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return nil, tea.Quit
		}
	}

	return s, nil
}

func (s *CombatProcessingState) View(ctx *game.Context) string {
	spin := s.spinner.View()
	content := fmt.Sprintf("%s The Dungeon Master judges your attack...", spin)
	return ui.CenteredView("⚔ COMBAT ⚔", content, true, ctx.Width, ctx.Height)
}
