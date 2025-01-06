package context_provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LinearContextProvider struct {
	APIKey  string
	IssueID string
}

type LinearContext struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (lc *LinearContext) IsEmpty() bool {
	return lc.Title == "" && lc.Description == ""
}

func (lp *LinearContextProvider) FetchContext() (ProvidedContext, error) {
	if lp.IssueID == "" {
		return nil, fmt.Errorf("issueID cannot be empty")
	}

	query := `
	query GetIssue($id: String!) {
		issue(id: $id) {
			title
			description
		}
	}`
	reqBody, _ := json.Marshal(map[string]interface{}{
		"query": query,
		"variables": map[string]string{
			"id": lp.IssueID,
		},
	})
	req, _ := http.NewRequest("POST", "https://api.linear.app/graphql", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", lp.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to contact Linear API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("linear api responded with status: %d, message: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Data struct {
			Issue LinearContext `json:"issue"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Linear response: %w", err)
	}

	if result.Data.Issue.Title == "" {
		return nil, fmt.Errorf("issue not found")
	}

	return &result.Data.Issue, nil
}
