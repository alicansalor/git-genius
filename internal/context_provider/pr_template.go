package context_provider

import (
	"fmt"
	"os"
)

type PRTemplateContextProvider struct {
	FilePath string
}

type PRTemplateContext struct {
	Template string
}

func (pt *PRTemplateContext) IsEmpty() bool {
	return pt.Template == ""
}

func (pt *PRTemplateContextProvider) FetchContext() (ProvidedContext, error) {
	if pt.FilePath == "" {
		return nil, fmt.Errorf("PR template path is not configured")
	}

	content, err := os.ReadFile(pt.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PR template: %w", err)
	}

	return &PRTemplateContext{Template: string(content)}, nil
}
