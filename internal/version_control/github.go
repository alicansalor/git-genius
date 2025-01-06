package versioncontrol

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

type GitHubManager struct {
	client *github.Client
	owner  string
	repo   string
}

// NewGitHubPRCreator initializes a new GitHub PR creator
func NewGitHubManager(ctx context.Context, token string) (*GitHubManager, error) {
	owner, repo, err := getOwnerAndRepo()
	if err != nil {
		return nil, fmt.Errorf("failed to determine owner and repo: %w", err)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubManager{
		client: github.NewClient(tc),
		owner:  owner,
		repo:   repo,
	}, nil
}

// CreatePullRequest creates a pull request on GitHub
func (g *GitHubManager) CreatePR(ctx context.Context, title, body, head, base string) (string, error) {
	newPR := &github.NewPullRequest{
		Title: github.Ptr(title),
		Head:  github.Ptr(head),
		Base:  github.Ptr(base),
		Body:  github.Ptr(body),
	}

	pr, _, err := g.client.PullRequests.Create(ctx, g.owner, g.repo, newPR)
	if err != nil {
		return "", fmt.Errorf("failed to create pull request: %w", err)
	}

	return pr.GetHTMLURL(), nil
}

func getOwnerAndRepo() (string, string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get remote URL: %w", err)
	}

	url := strings.TrimSpace(string(out))

	// Handle SSH and HTTPS URLs
	if strings.HasPrefix(url, "git@") {
		// SSH URL format: git@github.com:owner/repo.git
		parts := strings.SplitN(strings.TrimPrefix(url, "git@github.com:"), "/", 2)
		if len(parts) < 2 {
			return "", "", fmt.Errorf("invalid SSH URL format: %s", url)
		}
		return parts[0], strings.TrimSuffix(parts[1], ".git"), nil
	} else if strings.HasPrefix(url, "https://") {
		// HTTPS URL format: https://github.com/owner/repo.git
		parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
		if len(parts) < 2 {
			return "", "", fmt.Errorf("invalid HTTPS URL format: %s", url)
		}
		return parts[0], strings.TrimSuffix(parts[1], ".git"), nil
	}

	return "", "", fmt.Errorf("unsupported remote URL format: %s", url)
}
