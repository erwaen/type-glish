package llm

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ChatMessage represents a single message in the conversation history
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Assessment struct {
	CorrectedSentence  string `json:"corrected"`
	GrammarScore       int    `json:"score"`
	DamageDealt        int    `json:"damage"`
	DMComment          string `json:"dm_comment"`
	OutcomeDescription string `json:"outcome"`
}

type AssessmentMsg struct {
	Data Assessment
	Err  error
}

type Client struct {
	provider Provider
}

func NewClient(p Provider) *Client {
	return &Client{
		provider: p,
	}
}

func (c *Client) AnalyzeAction(userAction string) tea.Msg {
	// Construct the prompt
	messages := []ChatMessage{
		{Role: "system", Content: CriticPrompt},
		{Role: "user", Content: userAction},
	}

	resp, err := c.provider.Call(messages)
	if err != nil {
		return AssessmentMsg{Err: err}
	}

	var assessment Assessment
	err = json.Unmarshal([]byte(resp), &assessment)
	if err != nil {
		// Fallback if JSON is bad, or maybe the LLM refused JSON.
		// For now, return error.
		return AssessmentMsg{Err: fmt.Errorf("failed to parse JSON: %w", err)}
	}

	return AssessmentMsg{Data: assessment}
}
