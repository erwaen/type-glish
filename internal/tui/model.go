package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/config"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/states"
)

type MainModel struct {
	ctx   *game.Context    // The Data
	state states.GameState // The Current State
	cfg   *config.Config
}

func NewModel(ctx *game.Context, cfg *config.Config) MainModel {
	return MainModel{
		ctx:   ctx,
		state: states.NewMenuState(cfg),
		cfg:   cfg,
	}
}

func (m MainModel) Init() tea.Cmd {
	if m.state == nil {
		m.state = states.NewMenuState(m.cfg)
	}
	return m.state.Init(m.ctx)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window resize
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.ctx.Width = sizeMsg.Width
		m.ctx.Height = sizeMsg.Height
	}

	// Global Key Bindings
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "ctrl+s" {
			// Switch to Menu State
			newState := states.NewMenuState(m.cfg)
			m.state = newState
			return m, newState.Init(m.ctx)
		}
	}

	// Delegate Update to the current state
	newState, cmd := m.state.Update(msg, m.ctx)

	// Transition: If the state returned is different, switch
	if newState != nil && newState != m.state {
		oldState := m.state
		m.state = newState

		// Critical: If we just finished Settings or API Input, we likely updated Config but NOT the Context LLMClient.
		// We need to reload the Context's LLMClient based on the new Config.
		_, wasSettings := oldState.(*states.SettingsState)
		_, wasInput := oldState.(*states.APIInputState)

		if wasSettings || wasInput {
			// Reload Context
			m.ctx.ReloadLLM(m.cfg)
		}

		// Run the Init() of the NEW state immediately
		return m, tea.Batch(cmd, m.state.Init(m.ctx))
	}

	return m, cmd
}

func (m MainModel) View() string {
	if m.state == nil {
		return "Loading..."
	}
	return m.state.View(m.ctx)
}
