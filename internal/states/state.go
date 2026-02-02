package states

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/game"
)

type GameState interface{
		Init(ctx *game.Context) tea.Cmd

		Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd)

		View(ctx *game.Context) string
}


