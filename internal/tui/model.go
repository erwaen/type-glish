package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/states"
)

type MainModel struct {
	ctx   *game.Context    // The Data
	state states.GameState // The Current State
}

func NewModel(ctx *game.Context, initialState states.GameState) MainModel {
	return MainModel{
		ctx:   ctx,
		state: initialState,
	}
}

func (m MainModel) Init() tea.Cmd {
	if m.state != nil {
		return m.state.Init(m.ctx)
	}
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Delegate Update to the current state
	newState, cmd := m.state.Update(msg, m.ctx)

	// Transition: If the state returned is different, switch
	if newState != nil && newState != m.state {
		m.state = newState
		// Run the Init() of the NEW state immediately
		return m, tea.Batch(cmd, m.state.Init(m.ctx))
	}

	return m, cmd
}

func (m MainModel) View() string {
	return m.state.View(m.ctx)
}
