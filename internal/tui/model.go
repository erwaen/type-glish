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
	var initialState states.GameState

	// If provider is set, verify, otherwise show menu
	if cfg.Provider == "" {
		initialState = states.NewMenuState(cfg)
	} else {
		// If Gemini but no key, show input
		if cfg.Provider == "gemini" && cfg.GeminiAPIKey == "" {
			initialState = states.NewAPIInputState(cfg)
		} else {
			// Assume provider is set correctly, but Context needs to be synced if changed?
			// Actually, main.go loads config, creates Context *then* creates Model.
			// If config is valid, we start with Narrative (or whatever initial state).
			// Let's assume passed initialState is correct usually, BUT:
			// Main.go might not know if key is missing.

			// Simplification: Always start with Menu if Provider is empty,
			// else Narrative (if main passed it).
			// But wait, user requested "menu state where we can select".
			// If persisted, maybe we skip?
			// "restore the latest user option" -> Yes.

			// We'll trust `main.go` to provide the initial state based on config check,
			// OR we handle it here.
			// Since `main.go` calls `NewModel`, let's update `NewModel` to decide.
			initialState = &states.NarrativeState{}
			// But we need to use the one passed or inferred?
			// Let's change signature to NOT take initialState, but decide it?
			// Or keep it flexible.
		}
	}

	return MainModel{
		ctx:   ctx,
		state: initialState, // Will be overridden in Init logic maybe?
		cfg:   cfg,
	}
}

func (m MainModel) Init() tea.Cmd {
	// If state is nil (first run decision needed)
	if m.state == nil {
		if m.cfg.Provider == "" {
			m.state = states.NewMenuState(m.cfg)
		} else if m.cfg.Provider == "gemini" && m.cfg.GeminiAPIKey == "" {
			m.state = states.NewAPIInputState(m.cfg)
		} else {
			// Default start state
			m.state = &states.NarrativeState{}
		}
	}
	return m.state.Init(m.ctx)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.state = newState

		// Critical: If we just finished Menu or API Input, we likely updated Config but NOT the Context LLMClient.
		// We need to reload the Context's LLMClient based on the new Config.
		// We can do this by checking if the state was Menu or APIInput.
		_, wasMenu := m.state.(*states.MenuState)
		_, wasInput := m.state.(*states.APIInputState)

		// Actually, `m.state` here is the OLD state (before assignment below).
		if wasMenu || wasInput {
			// Reload Context
			m.ctx.ReloadLLM(m.cfg)
		}

		// Assign new state
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
