package llm

// Provider defines the interface for LLM backends
type Provider interface {
	Call(messages []ChatMessage) (string, error)
}
