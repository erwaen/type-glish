package game

import (
	"context"
	"fmt"

	"github.com/erwaen/type-glish/internal/config"
	"github.com/erwaen/type-glish/internal/llm"
)

type PlayerStats struct {
	HP         int
	XP         int
	Level      int
	Vocabulary []string // words "collected"
	Weaknesses []string // Grammar issues tracked
}

type Context struct {
	Stats            PlayerStats
	History          []llm.ChatMessage
	LastInput        string
	LastAssessment   llm.Assessment // structure result from the llm
	CurrentNarrative string         // The current story text displayed to the user
	LLMClient        *llm.Client

	// Terminal dimensions
	Width  int
	Height int
}

func NewContext(cfg *config.Config) *Context {
	ctx := &Context{
		Stats: PlayerStats{HP: 100, Level: 1},
	}
	ctx.ReloadLLM(cfg)
	return ctx
}

// ReloadLLM recreates the LLM client based on the provided config
func (c *Context) ReloadLLM(cfg *config.Config) {
	var provider llm.Provider
	var err error

	if cfg.Provider == "gemini" {
		if cfg.GeminiAPIKey == "" {
			fmt.Println("Warning: GEMINI_API_KEY not set, falling back to llamacpp")
			provider = llm.NewLlamaCppProvider()
		} else {
			provider, err = llm.NewGeminiProvider(context.Background(), cfg.GeminiAPIKey, cfg.GeminiModel)
			if err != nil {
				fmt.Printf("Error initializing Gemini: %v. Falling back to llamacpp\n", err)
				provider = llm.NewLlamaCppProvider()
			}
		}
	} else {
		provider = llm.NewLlamaCppProvider()
	}

	c.LLMClient = llm.NewClient(provider)
}
