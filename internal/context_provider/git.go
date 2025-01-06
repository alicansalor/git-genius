package context_provider

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

type GitContextProvider struct{}

type GitContext struct {
	Diff            string   // The diff of changes
	NewFiles        []string // A list of newly added or untracked files
	PreviousMessage []string // A list of previous commit messages
}

func (gc *GitContext) IsEmpty() bool {
	return gc.Diff == "" && len(gc.NewFiles) == 0 && len(gc.PreviousMessage) == 0
}

func (gp *GitContextProvider) FetchContext() (ProvidedContext, error) {
	// Fetch the diff
	diff, err := gp.getDiff()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Git diff: %w", err)
	}

	// Fetch the new files
	newFiles, err := gp.getNewFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch new files: %w", err)
	}

	// Fetch the previous commit messages
	previousMessages, err := gp.getPreviousMessages()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch previous commit messages: %w", err)
	}

	return &GitContext{
		Diff:            diff,
		NewFiles:        newFiles,
		PreviousMessage: previousMessages,
	}, nil
}

func (gp *GitContextProvider) getDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to fetch diff: %w", err)
	}
	return string(out), nil
}

func (gp *GitContextProvider) getNewFiles() ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch new files: %w", err)
	}

	var newFiles []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		// Identify new files marked with "A " or "??"
		if len(line) > 3 && (line[:2] == "A " || line[:2] == "??") {
			filePath := line[3:] // Extract file path
			newFiles = append(newFiles, filePath)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read Git status output: %w", err)
	}

	return newFiles, nil
}

func (gp *GitContextProvider) getPreviousMessages() ([]string, error) {
	cmd := exec.Command("git", "log", "-n", "20", "--pretty=format:%s")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch previous commit messages: %w", err)
	}

	// Parse commit messages into a slice
	var messages []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		messages = append(messages, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read Git log output: %w", err)
	}

	return messages, nil
}
