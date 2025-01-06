package llms

import "context"

type LLM interface {
	GenerateResponse(ctx context.Context, prompt string, maxTokens int) (string, error)
}
