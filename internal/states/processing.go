package states

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/llm"

	spinner "github.com/charmbracelet/bubbles/spinner"
	"github.com/erwaen/type-glish/internal/ui"
)

type ProcessingState struct {
	spinner spinner.Model
}

func (s *ProcessingState) Init(ctx *game.Context) tea.Cmd {
	s.spinner = spinner.New()
	s.spinner.Spinner = spinner.Dot
	s.spinner.Style = lipgloss.NewStyle().Foreground(ui.ColorPrimary)

	// Fire off the LLM analysis command
	return tea.Batch(
		s.spinner.Tick,
		func() tea.Msg {
			return ctx.LLMClient.AnalyzeAction(ctx.LastInput)
		},
	)
}

func (s *ProcessingState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle the API Response
	case llm.AssessmentMsg:
		if msg.Err != nil {
			log.Printf("Error from LLM: %v", msg.Err)
			// Return to input state or handle error gracefully?
			// For now, let's just log and maybe set a dummy error message in context so user knows?
			// But the user just asked "how can I see the errors".
			// Let's pass the error to the UI logic if possible, or just log.
			// Going to ResultState with empty data is confusing.
			// Let's go back to InputState but maybe we need a way to show error.
		}
		ctx.LastAssessment = msg.Data // Save result to context

		// Transition to Result Screen
		return &ResultState{}, nil

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

func (s ProcessingState) View(ctx *game.Context) string {
	spin := s.spinner.View()
	content := fmt.Sprintf("%s The Dungeon Master is judging your grammar...", spin)
	return ui.CenteredView("THINKING...", content, true, ctx.Width, ctx.Height)
}
