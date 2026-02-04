package states

import (
	"fmt"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/ui"
)

type VictoryState struct {
	defeatedEnemy string
}

func (s *VictoryState) Init(ctx *game.Context) tea.Cmd {
	if ctx.CurrentEnemy != nil {
		s.defeatedEnemy = ctx.CurrentEnemy.Name
	}
	return nil
}

func (s *VictoryState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// Award XP
			ctx.Stats.XP += 10

			// 50% chance: new combat or path choice
			if rand.Float32() < 0.5 {
				// New combat with random enemy
				ctx.CurrentEnemy = game.RandomEnemy()
				ctx.CurrentNarrative = fmt.Sprintf("A %s appears! %s", ctx.CurrentEnemy.Name, ctx.CurrentEnemy.Description)
				return NewCombatState(), nil
			} else {
				// Path choice for healing opportunity
				return NewPathChoiceState(), nil
			}
		}
		if msg.Type == tea.KeyCtrlC {
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s *VictoryState) View(ctx *game.Context) string {
	content := fmt.Sprintf(`
    ╔═══════════════════════════════════════╗
    ║                                       ║
    ║          V I C T O R Y !              ║
    ║                                       ║
    ║   You have defeated the %s!
    ║                                       ║
    ║   Your mastery of grammar prevails.   ║
    ║                                       ║
    ║   +10 XP                              ║
    ║                                       ║
    ╚═══════════════════════════════════════╝
`, s.defeatedEnemy)

	content += "\n\n"
	content += ui.RenderHPBar(ctx.Stats.HP, 100, "Your HP", 20) + "\n"
	content += fmt.Sprintf("XP: %d\n\n", ctx.Stats.XP)
	content += ui.StyleHelp.Render("Press [Enter] to continue your journey...")

	return ui.CenteredView("🏆 VICTORY 🏆", content, true, ctx.Width, ctx.Height)
}

