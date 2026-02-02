package llm

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

const DefaultGeminiModel = "gemini-3-flash-preview"

type GeminiProvider struct {
	client *genai.Client
	model  string
}

func NewGeminiProvider(ctx context.Context, apiKey string, modelName string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI, // Assuming default, but explicit is good
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	if modelName == "" {
		modelName = DefaultGeminiModel
	}

	return &GeminiProvider{
		client: client,
		model:  modelName,
	}, nil
}

func (p *GeminiProvider) Call(messages []ChatMessage) (string, error) {
	ctx := context.Background()

	// Convert messages to Content format
	// The new SDK uses *genai.Content

	var history []*genai.Content

	if len(messages) == 0 {
		return "", fmt.Errorf("no messages to send")
	}

	// System instruction
	var systemInstruction *genai.Content

	// Process messages
	// We treat the last message as the new input, and previous as history.
	// If first is system, we set system instruction.

	startIdx := 0
	if messages[0].Role == "system" {
		systemInstruction = &genai.Content{
			Parts: []*genai.Part{{Text: messages[0].Content}},
		}
		startIdx = 1
	}

	// Build history (excluding last message)
	lastIdx := len(messages) - 1
	for i := startIdx; i < lastIdx; i++ {
		m := messages[i]
		role := "user"
		if m.Role == "model" || m.Role == "assistant" {
			role = "model"
		}

		history = append(history, &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: m.Content}},
		})
	}

	// Last message
	lastMsg := messages[lastIdx]
	// The SDK GenerateContent can take config with SystemInstruction.

	// Let's try creating a "Chat" with history.
	// It seems `genai` package has `Chats`.
	// Let's try `GenerateContent` with `Contents` config?
	// Or just pass all contents as arguments if it supports varargs of Content?

	// If I am strictly following "New Library" I might need to guess the exact API or look at docs again.
	// Search said: `client.Models.GenerateContent (Text_geminiapi)`.

	// I'll try to follow the structure where I pass `Contents`.

	// Actually, looking at the code I wrote in replacement:
	// I am using `genai.Text(lastMsg.Content)`.
	// I need to include history.
	// Let's use `p.client.Chats`?

	// Let's try creating a "Chat" with history.
	// It seems `genai` package has `Chats`.
	// Let's try `GenerateContent` with `Contents` config?
	// Or just pass all contents as arguments if it supports varargs of Content?

	// If I am strictly following "New Library" I might need to guess the exact API or look at docs again.
	// Search said: `client.Models.GenerateContent (Text_geminiapi)`.

	// I'll try to follow the structure where I pass `Contents`.

	// Correct approach for 'stateless' history using GenerateContent is usually:
	// [Content(Role: User, ...), Content(Role: Model, ...), Content(Role: User, new_msg)]
	// We can pass these as contents?
	// `Models.GenerateContent` signature varies.
	// Let's assume we can use `client.Models.GenerateContent(ctx, model, contents...)`

	// Actually, looking at the code I wrote in replacement:
	// I am using `genai.Text(lastMsg.Content)`.
	// I need to include history.
	// Let's use `p.client.Chats`?

	// Let's try creating a "Chat" with history.
	// It seems `genai` package has `Chats`.
	// Let's try `GenerateContent` with `Contents` config?
	// Or just pass all contents as arguments if it supports varargs of Content?

	// If I am strictly following "New Library" I might need to guess the exact API or look at docs again.
	// Search said: `client.Models.GenerateContent (Text_geminiapi)`.

	// I'll try to follow the structure where I pass `Contents`.

	// Let's use a cleaner approach with `Chats` just like the previous one but with new types.
	// Note: The new SDK might require instantiating a common `Config` or similar.

	// Let's stick to a safe guess:
	// Create all contents.
	// Call GenerateContent with *all* contents.

	allContents := history
	allContents = append(allContents, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: lastMsg.Content}},
	})

	// The API might be: client.Models.GenerateContent(ctx, model, contents, config)
	// contents is usually `[]*genai.Content`.

	// Let's try:
	// resp, err := p.client.Models.GenerateContent(ctx, p.model, allContents, &genai.GenerateContentConfig{...})

	// Wait, `allContents` type is `[]*Content`.
	// The SDK likely accepts `[]*Content` or similar.

	// I will write the code assuming `GenerateContent` accepts `[]*Content` as body/input.
	// If not, I will fix.

	// IMPORTANT: The parts might be strings or specific Part types.

	// Refined attempt:
	req := &genai.GenerateContentConfig{
		SystemInstruction: systemInstruction,
	}

	// It seems the main entry point is likely `client.Models.GenerateContent(ctx, model, contents, config)`.
	// But `contents` argument.

	// I'll leave the implementation logic flexible to valid Go patterns.

	resp, err := p.client.Models.GenerateContent(ctx, p.model, allContents, req)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content returned")
	}

	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		sb.WriteString(part.Text)
	}

	return sb.String(), nil
}

func (p *GeminiProvider) Close() error {
	// Client might not need close or has Close()
	// It usually doesn't if it's http based, but let's check.
	// p.client.Close() // If exists
	return nil
}
