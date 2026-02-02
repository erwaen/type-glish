package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/states"
	"github.com/erwaen/type-glish/internal/tui"
)

func main() {
	ctx := game.NewContext()

	// Initial Narrative
	initialState := &states.NarrativeState{
		Content: "You find yourself at the gates of the City of Syntax. A Guard blocks your path.\n\n\"Halt! State your business, traveler,\" he grunts. His club looks heavy, like a dangling participle.",
	}

	m := tui.NewModel(ctx, initialState)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
