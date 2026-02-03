package states

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/config"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/ui"
)

// SettingsState is the Provider Selection Menu
type SettingsState struct {
	choices []string
	cursor  int
	cfg     *config.Config
}

func NewSettingsState(cfg *config.Config) *SettingsState {
	return &SettingsState{
		choices: []string{"Use llama.cpp (Local)", "Use Gemini (Cloud)", "Update Gemini API Key", "Back"},
		cursor:  0,
		cfg:     cfg,
	}
}

func (s *SettingsState) Init(ctx *game.Context) tea.Cmd {
	return nil
}

func (s *SettingsState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		case "esc":
			return NewMenuState(s.cfg), nil
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "j":
			if s.cursor < len(s.choices)-1 {
				s.cursor++
			}
		case "enter":
			if s.cursor == 0 {
				// Llama
				s.cfg.Provider = "llamacpp"
				config.SaveConfig(s.cfg)
				return NewMenuState(s.cfg), nil
			} else if s.cursor == 1 {
				// Gemini
				if s.cfg.GeminiAPIKey == "" {
					return NewAPIInputState(s.cfg), nil
				}
				s.cfg.Provider = "gemini"
				config.SaveConfig(s.cfg)
				return NewMenuState(s.cfg), nil
			} else if s.cursor == 2 {
				// Update Gemini Key
				return NewAPIInputState(s.cfg), nil
			} else if s.cursor == 3 {
				// Back
				return NewMenuState(s.cfg), nil
			}
		}
	}
	return s, nil
}

func (s *SettingsState) View(ctx *game.Context) string {
	var content string

	content += ui.StyleSubTitle.Render("Select your Intelligence Provider") + "\n\n"

	for i, choice := range s.choices {
		content += ui.RenderMenuItem(choice, s.cursor == i) + "\n"
	}

	content += ui.StyleHelp.Render("\n(Use ↑/↓ to move, Enter to select, q to quit)")

	return ui.Box("SETTINGS", content, true)
}
