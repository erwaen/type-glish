package states

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/config"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/ui"
)

// MenuState is the Main Menu (Start Game, Settings)
type MenuState struct {
	choices []string
	cursor  int
	cfg     *config.Config
}

func NewMenuState(cfg *config.Config) *MenuState {
	return &MenuState{
		choices: []string{"Start Game", "Settings"},
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
				// Start Game
				if s.cfg.Provider == "" {
					return NewSettingsState(s.cfg), nil
				}
				if s.cfg.Provider == "gemini" && s.cfg.GeminiAPIKey == "" {
					return NewAPIInputState(s.cfg), nil
				}
				return &NarrativeState{}, nil
			} else if s.cursor == 1 {
				// Settings
				return NewSettingsState(s.cfg), nil
			}
		}
	}
	return s, nil
}

func (s *MenuState) View(ctx *game.Context) string {
	var content string

	content += ui.StyleSubTitle.Render("Welcome to Type-Glish") + "\n\n"

	for i, choice := range s.choices {
		content += ui.RenderMenuItem(choice, s.cursor == i) + "\n"
	}

	content += ui.StyleHelp.Render("\n(Use ↑/↓ to move, Enter to select, q to quit)")

	return ui.Box("MAIN MENU", content, true)
}
