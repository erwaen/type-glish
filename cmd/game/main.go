package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erwaen/type-glish/internal/config"
	"github.com/erwaen/type-glish/internal/game"
	"github.com/erwaen/type-glish/internal/tui"
)

func main() {
	// Because tea doesn't now allow me to see prints, I can store in a file to check it
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	
	// creates the game data and setup the llm provider
	ctx := game.NewContext(cfg)
	
	// set a new model and the current state of the game, so menu
	m := tui.NewModel(ctx, cfg)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
