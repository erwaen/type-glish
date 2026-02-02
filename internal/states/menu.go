package states

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/config"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/ui"
)

type MenuState struct {
	choices []string
	cursor  int
	cfg     *config.Config
}

func NewMenuState(cfg *config.Config) *MenuState {
	return &MenuState{
		choices: []string{"Use llama.cpp (Local)", "Use Gemini (Cloud)", "Update Gemini API Key"},
		cursor:  0,
		cfg:     cfg,
	}
}

func (s *MenuState) Init(ctx *game.Context) tea.Cmd {
	return nil
}

func (s *MenuState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
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
				return &NarrativeState{}, nil
			} else if s.cursor == 1 {
				// Gemini
				if s.cfg.GeminiAPIKey == "" {
					return NewAPIInputState(s.cfg), nil
				}
				s.cfg.Provider = "gemini"
				config.SaveConfig(s.cfg)
				return &NarrativeState{}, nil
			} else if s.cursor == 2 {
				// Update Gemini Key
				return NewAPIInputState(s.cfg), nil
			}
		}
	}
	return s, nil
}

func (s *MenuState) View(ctx *game.Context) string {
	var content string

	// Add description
	content += ui.StyleSubTitle.Render("Select your Intelligence Provider") + "\n\n"

	for i, choice := range s.choices {
		content += ui.RenderMenuItem(choice, s.cursor == i) + "\n"
	}

	content += ui.StyleHelp.Render("\n(Use ↑/↓ to move, Enter to select, q to quit)")

	// Centering or just box
	// Using Box with active state
	return ui.Box("CONFIGURATION", content, true)
}
