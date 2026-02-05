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
	goldEarned    int
}

func (s *VictoryState) Init(ctx *game.Context) tea.Cmd {
	if ctx.CurrentEnemy != nil {
		s.defeatedEnemy = ctx.CurrentEnemy.Name
		// Calculate gold reward: base + random bonus based on tier
		tier := ctx.CurrentEnemy.Tier
		if tier < 1 {
			tier = 1
		}
		baseGold := tier * 5
		bonusGold := rand.Intn(tier*3 + 1)
		s.goldEarned = baseGold + bonusGold
	}
	return nil
}

func (s *VictoryState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// Award XP and Gold
			ctx.Stats.XP += 10
			ctx.Stats.Gold += s.goldEarned

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
	// Gold display styling
	goldStyle := ui.StyleDamageDealt

	content := fmt.Sprintf(`
    ╔═══════════════════════════════════════╗
    ║                                       ║
    ║          V I C T O R Y !              ║
    ║                                       ║
    ║   You have defeated the %s!
    ║                                       ║
    ║   Your mastery of grammar prevails.   ║
    ║                                       ║
    ╚═══════════════════════════════════════╝
`, s.defeatedEnemy)

	content += "\n"
	content += fmt.Sprintf("    +10 XP    %s\n\n", goldStyle.Render(fmt.Sprintf("+%d Gold", s.goldEarned)))
	content += "───────────────────────────────────────────\n\n"
	content += ui.RenderStatusBar(ctx.Stats.HP, 100, ctx.Stats.Gold+s.goldEarned, ctx.Stats.XP+10) + "\n\n"
	content += ui.StyleHelp.Render("Press [Enter] to continue your journey...")

	return ui.CenteredView("VICTORY", content, true, ctx.Width, ctx.Height)
}

