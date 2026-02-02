package game

import "github.com/erwaen/type-glish/internal/llm"

type PlayerStats struct {
	HP         int
	XP         int
	Level      int
	Vocabulary []string // words "collected"
	Weaknesses []string // Grammar issues tracked
}

type Context struct {
	Stats          PlayerStats
	History        []llm.ChatMessage
	LastInput      string
	LastAssessment llm.Assessment // structure result from the llm
	LLMClient      *llm.Client
}

func NewContext() *Context {
	return &Context{
		Stats:     PlayerStats{HP: 100, Level: 1},
		LLMClient: llm.NewClient(),
	}
}
