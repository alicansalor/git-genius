package llms

import (
	"context"
	"fmt"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Gemini struct {
	Client *genai.Client
}

func NewGemini(ctx context.Context, apiKey string) (*Gemini, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey((apiKey)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	return &Gemini{Client: client}, nil
}

func (g *Gemini) GenerateResponse(ctx context.Context, prompt string, maxTokens int) (string, error) {
	model := g.Client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// select the first candidate
	if len(resp.Candidates) == 0 {
		return "", nil
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil {
		return "", nil
	}

	var content string
	for _, part := range candidate.Content.Parts {
		content += fmt.Sprint(part)
	}

	return content, nil
}
