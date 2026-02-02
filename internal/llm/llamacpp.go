package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	serverURL = "http://127.0.0.1:8080/v1/chat/completions"
)

type LlamaCppProvider struct{}

func NewLlamaCppProvider() *LlamaCppProvider {
	return &LlamaCppProvider{}
}

type chatRequest struct {
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
	Stream      bool          `json:"stream"`
}

type chatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func parseLLMResponse(resp *http.Response, err error) (string, error) {
	if err != nil {
		return "", fmt.Errorf("failed to call llm server: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server error: %s", string(body))
	}

	var result chatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no content received from LLM")
	}

	// Clean up the output (remove extra spaces or newlines)
	cleanText := strings.TrimSpace(result.Choices[0].Message.Content)

	// Sometimes LLMs add quotes even when told not to, handle that safely
	cleanText = strings.Trim(cleanText, "\"")

	cleanText = strings.ReplaceAll(cleanText, "’", "'")

	cleanText = strings.ReplaceAll(cleanText, "‘", "'")

	return cleanText, nil
}

// Call sends a request to the local LLM
func (p *LlamaCppProvider) Call(messages []ChatMessage) (string, error) {
	reqBody := chatRequest{
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   300,
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	return parseLLMResponse(resp, err)
}
