package sdk

import (
	"context"
	"fmt"

	versioncontrol "git-genius/internal/version_control"
	llm "git-genius/internal/llm"
)

// GitGeniusSDK defines the public interface for the SDK
type GitGeniusSDK interface {
	GenerateCommitMessage(ctx context.Context, issueID string) (string, error)
	CreatePullRequest(ctx context.Context, issueID string) (string, error)
}

type gitGeniusSDK struct{
	llm llm.LLM,
	prCreator versioncontrol.PRCreator,

}

// NewGeniusSDK creates a new GeniusSDK instance
func NewGeniusSDK(llm llm.LLM, prCreator versioncontrol.PRCreator) GitGeniusSDK {
	return &gitGeniusSDK{
		llm,
		prCreator,
	}
}

func (g *gitGeniusSDK) GenerateCommitMessage(ctx context.Context, issueID string) (string, error) {
	// Simulated logic for generating a commit message
	if issueID == "" {
		return "", fmt.Errorf("issueID cannot be empty")
	}

	return fmt.Sprintf("Generated commit message for issue %s", issueID), nil
}

func (g *gitGenius) CreatePullRequest(ctx context.Context, issueID string) (string, error) {
	// Simulated logic for creating a pull request
	if issueID == "" {
		return "", fmt.Errorf("issueID cannot be empty")
	}

	return fmt.Sprintf("Pull Request URL for issue %s", issueID), nil
}
