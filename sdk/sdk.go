package sdk

import (
	"context"
	"fmt"

	"git-genius/config"
	context_provider "git-genius/internal/context_provider"
	llm "git-genius/internal/llm"
	versioncontrol "git-genius/internal/version_control"
)

type GitGenius interface {
	GenerateCommitMessage(ctx context.Context) (string, error)
	GeneratePullRequestContent(ctx context.Context) (*PullRequestContent, error)
}

type GitGeniusSDK struct {
	llm            llm.LLM
	prCreator      versioncontrol.PRCreator
	contextManager *context_provider.ContextManager
}

// NewGitGeniusSDK creates a new GeniusSDK instance
func NewGitGeniusSDK(ctx context.Context, cfg *config.Config) (*GitGeniusSDK, error) {

	if cfg == nil {
		return nil, fmt.Errorf("config is not defined")
	}

	// crete a new LLM instance
	llm, err := cfg.NewLLM(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create llm: %v", err)
	}

	// create a new PRCreator instance
	prCreator, err := cfg.NewPRCreator(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create pr creator: %v", err)
	}
	// create context manager
	contextManager := context_provider.NewContextManager(cfg)

	return &GitGeniusSDK{
		llm,
		prCreator,
		contextManager,
	}, nil
}

func (g *GitGeniusSDK) GenerateCommitMessage(ctx context.Context) (string, error) {
	context, err := g.contextManager.CollectContext()
	if err != nil {
		return "", fmt.Errorf("failed to collect context: %v", err)
	}

	var gitPrompt string
	if context.Git != nil && !context.Git.IsEmpty() {
		gitPrompt = fmt.Sprintf(`
			Git Diff: %v,
			Git New Files: %v,
			Git Previous commit messages: %v
		`, context.Git.Diff, context.Git.NewFiles, context.Git.PreviousMessage)
	}

	// if there is no git context then we can't generate a commit message
	if gitPrompt == "" {
		return "", fmt.Errorf("no git context found")
	}

	var linearPrompt string
	if context.Linear != nil && !context.Linear.IsEmpty() {
		linearPrompt = fmt.Sprintf(`
			Linear ticket description: %v,
			Linear ticket title: %v
		`, context.Linear.Description, context.Linear.Title)
	}

	prompt := fmt.Sprintf("Generate a concise git commit message using the following: %v %v", linearPrompt, gitPrompt)

	// Generate commit message
	commitMessage, err := g.llm.GenerateResponse(ctx, prompt, 150)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %v", err)
	}

	return commitMessage, nil
}

func (g *GitGeniusSDK) GeneratePullRequestContent(ctx context.Context) (*PullRequestContent, error) {

	context, err := g.contextManager.CollectContext()
	if err != nil {
		return nil, fmt.Errorf("failed to collect context: %v", err)
	}

	var prTemplatePrompt string
	if context.PRTemplate != nil && !context.PRTemplate.IsEmpty() {
		prTemplatePrompt = fmt.Sprintf(`Write the PR according to the template: %v
			`, context.PRTemplate.Template)
	}

	// if there is no PR template context then we can't generate a PR description
	if prTemplatePrompt == "" {
		return nil, fmt.Errorf("no PR template context found")
	}

	var linearPrompt string
	if context.Linear != nil && !context.Linear.IsEmpty() {
		linearPrompt = fmt.Sprintf(`
				Linear ticket description: %v,
				Linear ticket title: %v
			`, context.Linear.Description, context.Linear.Title)
	}

	// all the commits made in the branch
	var gitPrompt string
	if context.Git != nil && !context.Git.IsEmpty() {
		gitPrompt = fmt.Sprintf(`
				Git commit messages: %v
			`, context.Git.PreviousMessage)
	}

	// if there is no git context then we can't generate a commit message
	if gitPrompt == "" {
		return nil, fmt.Errorf("no git context found")
	}

	titlePrompt := fmt.Sprintf(`Generate a one line PR title
		using the following: %v %v`, linearPrompt, gitPrompt)

	title, err := g.llm.GenerateResponse(ctx, titlePrompt, 150)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %v", err)
	}

	bodyPrompt := fmt.Sprintf(`Generate a PR description
		using the following: %v %v %v`, linearPrompt, gitPrompt, prTemplatePrompt)

	body, err := g.llm.GenerateResponse(ctx, bodyPrompt, 150)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %v", err)
	}

	return &PullRequestContent{
		Title: title,
		Body:  body,
		Tags:  []string{},
	}, nil
}
