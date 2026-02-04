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

// CombatAssessment is the LLM response for combat actions
type CombatAssessment struct {
	CorrectedSentence string `json:"corrected"`
	GrammarScore      int    `json:"score"`
	DamageDealt       int    `json:"damage_dealt"`
	DamageReceived    int    `json:"damage_received"`
	DMComment         string `json:"dm_comment"`
	Outcome           string `json:"outcome"`
	IsRelevant        bool   `json:"is_relevant"`
}

type CombatAssessmentMsg struct {
	Data CombatAssessment
	Err  error
}

// PathAssessment is the LLM response for path choices
type PathAssessment struct {
	CorrectedSentence string `json:"corrected"`
	GrammarScore      int    `json:"score"`
	Healing           int    `json:"healing"`
	DMComment         string `json:"dm_comment"`
	Outcome           string `json:"outcome"`
	IsRelevant        bool   `json:"is_relevant"`
}

type PathAssessmentMsg struct {
	Data PathAssessment
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
		return AssessmentMsg{Err: fmt.Errorf("failed to parse JSON: %w", err)}
	}

	return AssessmentMsg{Data: assessment}
}

// AnalyzeCombatAction analyzes a combat action with enemy context
func (c *Client) AnalyzeCombatAction(userAction, enemyName, location string) tea.Msg {
	prompt := fmt.Sprintf(CombatPromptTemplate, enemyName, location)

	messages := []ChatMessage{
		{Role: "system", Content: prompt},
		{Role: "user", Content: userAction},
	}

	resp, err := c.provider.Call(messages)
	if err != nil {
		return CombatAssessmentMsg{Err: err}
	}

	var assessment CombatAssessment
	err = json.Unmarshal([]byte(resp), &assessment)
	if err != nil {
		return CombatAssessmentMsg{Err: fmt.Errorf("failed to parse combat JSON: %w", err)}
	}

	return CombatAssessmentMsg{Data: assessment}
}

// AnalyzePathChoice analyzes a path choice for healing
func (c *Client) AnalyzePathChoice(userChoice, pathOptions string) tea.Msg {
	prompt := fmt.Sprintf(PathChoicePromptTemplate, pathOptions)

	messages := []ChatMessage{
		{Role: "system", Content: prompt},
		{Role: "user", Content: userChoice},
	}

	resp, err := c.provider.Call(messages)
	if err != nil {
		return PathAssessmentMsg{Err: err}
	}

	var assessment PathAssessment
	err = json.Unmarshal([]byte(resp), &assessment)
	if err != nil {
		return PathAssessmentMsg{Err: fmt.Errorf("failed to parse path JSON: %w", err)}
	}

	return PathAssessmentMsg{Data: assessment}
}
