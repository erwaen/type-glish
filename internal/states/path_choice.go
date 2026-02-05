package states

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/llm"
	"github.com/erwaen/type-glish/internal/ui"
)

// Path options
var PathOptions = []struct {
	Name        string
	Description string
}{
	{"The Misty Forest", "A winding path through ancient trees shrouded in fog."},
	{"The Crystal Cave", "A glittering cavern with echoing whispers."},
	{"The Old Bridge", "A creaky wooden bridge over a rushing river."},
	{"The Abandoned Tower", "A crumbling tower that once housed great scholars."},
	{"The Meadow of Echoes", "A peaceful meadow where your words linger."},
}

type PathChoiceState struct {
	textInput textinput.Model
	paths     []struct {
		Name        string
		Description string
	}
}

func NewPathChoiceState() *PathChoiceState {
	ti := textinput.New()
	ti.Placeholder = "Describe which path you take..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	// Pick 3 random paths
	shuffled := make([]struct {
		Name        string
		Description string
	}, len(PathOptions))
	copy(shuffled, PathOptions)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return &PathChoiceState{
		textInput: ti,
		paths:     shuffled[:3],
	}
}

func (s *PathChoiceState) Init(ctx *game.Context) tea.Cmd {
	return textinput.Blink
}

func (s *PathChoiceState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if s.textInput.Value() != "" {
				ctx.LastInput = s.textInput.Value()
				// Build path options string for LLM
				pathStr := ""
				for i, p := range s.paths {
					pathStr += fmt.Sprintf("%d. %s - %s\n", i+1, p.Name, p.Description)
				}
				return &PathProcessingState{pathOptions: pathStr}, nil
			}
		case tea.KeyCtrlC:
			return s, tea.Quit
		}
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s *PathChoiceState) View(ctx *game.Context) string {
	var content string

	// Status bar at top
	content += ui.RenderStatusBar(ctx.Stats.HP, 100, ctx.Stats.Gold, ctx.Stats.XP) + "\n\n"

	content += ui.StyleSubTitle.Render("You come to a crossroads...") + "\n\n"

	content += "Choose your path:\n\n"
	for i, p := range s.paths {
		content += fmt.Sprintf("  %d. %s\n", i+1, ui.StyleEnemyName.Render(p.Name))
		content += fmt.Sprintf("     %s\n\n", ui.StyleSubTitle.Render(p.Description))
	}

	content += "───────────────────────────────────────────\n\n"

	content += "Describe your choice in a complete sentence:\n"
	content += "> " + s.textInput.View() + "\n\n"

	content += ui.StyleHelp.Render("(Better grammar = more healing!)")

	return ui.CenteredView("CROSSROADS", content, true, ctx.Width, ctx.Height)
}

// PathProcessingState processes the path choice
type PathProcessingState struct {
	spinner     spinner.Model
	pathOptions string
}

func (s *PathProcessingState) Init(ctx *game.Context) tea.Cmd {
	s.spinner = spinner.New()
	s.spinner.Spinner = spinner.Dot
	s.spinner.Style = lipgloss.NewStyle().Foreground(ui.ColorPrimary)

	pathOpts := s.pathOptions

	return tea.Batch(
		s.spinner.Tick,
		func() tea.Msg {
			return ctx.LLMClient.AnalyzePathChoice(ctx.LastInput, pathOpts)
		},
	)
}

func (s *PathProcessingState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case llm.PathAssessmentMsg:
		if msg.Err != nil {
			log.Printf("Error from LLM: %v", msg.Err)
			// Default healing on error
			healing := 10
			ctx.Stats.HP += healing
			if ctx.Stats.HP > 100 {
				ctx.Stats.HP = 100
			}
			return &PathResultState{
				healing:   healing,
				outcome:   "You find a peaceful spot to rest...",
				dmComment: "The narrator lost their notes.",
				corrected: ctx.LastInput,
				score:     5,
			}, nil
		}

		// Apply healing
		healing := msg.Data.Healing
		if healing > 20 {
			healing = 20
		}
		ctx.Stats.HP += healing
		if ctx.Stats.HP > 100 {
			ctx.Stats.HP = 100
		}

		return &PathResultState{
			healing:   healing,
			outcome:   msg.Data.Outcome,
			dmComment: msg.Data.DMComment,
			corrected: msg.Data.CorrectedSentence,
			score:     msg.Data.GrammarScore,
		}, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return nil, tea.Quit
		}
	}

	return s, nil
}

func (s *PathProcessingState) View(ctx *game.Context) string {
	spin := s.spinner.View()
	content := fmt.Sprintf("%s The Dungeon Master considers your path...", spin)
	return ui.CenteredView("🛤 CROSSROADS 🛤", content, true, ctx.Width, ctx.Height)
}

// PathResultState shows the result of path choice
type PathResultState struct {
	healing   int
	outcome   string
	dmComment string
	corrected string
	score     int
}

func (s *PathResultState) Init(ctx *game.Context) tea.Cmd {
	return nil
}

func (s *PathResultState) Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// Spawn new enemy and go to combat
			ctx.CurrentEnemy = game.RandomEnemy()
			ctx.CurrentNarrative = fmt.Sprintf("As you travel, a %s blocks your path! %s",
				ctx.CurrentEnemy.Name, ctx.CurrentEnemy.Description)
			return NewCombatState(), nil
		}
		if msg.Type == tea.KeyCtrlC {
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s *PathResultState) View(ctx *game.Context) string {
	var content string

	content += ui.StyleSubTitle.Render("YOUR CHOICE:") + "\n"
	content += "> " + s.corrected + "\n\n"

	content += "───────────────────────────────────────────\n\n"

	content += s.outcome + "\n\n"

	// Score icons based on grammar score
	var scoreIcons string
	if s.score >= 9 {
		scoreIcons = "★★★"
	} else if s.score >= 7 {
		scoreIcons = "★★☆"
	} else if s.score >= 5 {
		scoreIcons = "★☆☆"
	} else {
		scoreIcons = "☆☆☆"
	}

	scoreColor := ui.ColorError
	if s.score >= 7 {
		scoreColor = ui.ColorSuccess
	} else if s.score >= 5 {
		scoreColor = ui.ColorWarning
	}

	scoreStyle := lipgloss.NewStyle().Foreground(scoreColor).Bold(true)
	healStyle := lipgloss.NewStyle().Foreground(ui.ColorSuccess).Bold(true)

	content += fmt.Sprintf("Score: %s %s  |  Health Restored: %s\n\n",
		scoreStyle.Render(fmt.Sprintf("%d/10", s.score)),
		scoreStyle.Render(scoreIcons),
		healStyle.Render(fmt.Sprintf("+%d", s.healing)))

	content += ui.StyleSubTitle.Render("DM:") + " " + s.dmComment + "\n\n"

	content += "───────────────────────────────────────────\n\n"
	content += ui.RenderStatusBar(ctx.Stats.HP, 100, ctx.Stats.Gold, ctx.Stats.XP) + "\n\n"

	content += ui.StyleHelp.Render("Press [Enter] to continue...")

	return ui.CenteredView("PATH RESULT", content, true, ctx.Width, ctx.Height)
}
