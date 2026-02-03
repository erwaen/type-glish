package states

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/config"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/ui"
)

type APIInputState struct {
	textInput textinput.Model
	cfg       *config.Config
}

func NewAPIInputState(cfg *config.Config) *APIInputState {
	ti := textinput.New()
	ti.Placeholder = "Enter your Gemini API Key"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 60
	ti.EchoMode = textinput.EchoPassword // Hide key

	return &APIInputState{
		textInput: ti,
		cfg:       cfg,
	}
}

func (s *APIInputState) Init(ctx *game.Context) tea.Cmd {
	return textinput.Blink
}

func (s *APIInputState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			apiKey := s.textInput.Value()
			if apiKey != "" {
				s.cfg.GeminiAPIKey = apiKey
				s.cfg.Provider = "gemini"
				// Save config
				if err := config.SaveConfig(s.cfg); err != nil {
					// Handle error? For now print to stdout is bad in TUI.
					// Just proceed.
				}
				// Re-initialize Context with new config?
				// The MainModel needs to handle this transition or Context update.
				// Since Context is passed by pointer, we can try to update it here.
				// But Context has LLMClient which needs recreation.
				// For now, let's assume we return a State that triggers a Context Reload or we update Context here.
				// BUT `Init` of the next state (Narrative?) assumes Context is ready.

				// Let's defer context update to the transition logic in MainModel if possible,
				// or just do it here:
				// ctx.ReloadLLM(s.cfg) -> We need to implement this method on Context.

				// Returning to Settings
				return NewSettingsState(s.cfg), nil
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return s, tea.Quit
		}
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s *APIInputState) View(ctx *game.Context) string {
	content := "To use Gemini, we need an API Key.\nIt will be saved locally.\n\n" + s.textInput.View()
	return ui.CenteredView("GEMINI SETUP", content, true, ctx.Width, ctx.Height)
}
