package context_provider

import (
	"errors"

	"git-genius/config"
)

type Context struct {
	Linear     *LinearContext
	Git        *GitContext
	PRTemplate *PRTemplateContext
}

type ContextManager struct {
	Config *config.Config
}

func NewContextManager(cfg *config.Config) *ContextManager {
	return &ContextManager{Config: cfg}
}

func (cm *ContextManager) CollectContext() (*Context, error) {
	context := &Context{}

	for _, contextProviderConfig := range cm.Config.ContextProviders {

		var contextProvider ContextProvider
		switch contextProviderConfig.Name {
		case string(LinearContextProviderType):
			contextProvider = &LinearContextProvider{
				APIKey:  contextProviderConfig.APIKey,
				IssueID: cm.Config.IssueID,
			}
		case string(GitContextProviderType):
			contextProvider = &GitContextProvider{}
		case string(PRTemplateContextProviderType):
			contextProvider = &PRTemplateContextProvider{
				FilePath: contextProviderConfig.Path,
			}
		default:
			return nil, errors.New("unknown context_provider: " + contextProviderConfig.Name)
		}

		data, err := contextProvider.FetchContext()
		if err != nil {
			return nil, err
		}

		// Populate the typed context
		switch v := data.(type) {
		case *LinearContext:
			context.Linear = v
		case *GitContext:
			context.Git = v
		case *PRTemplateContext:
			context.PRTemplate = v
		default:
			return nil, errors.New("unexpected context_provider response type")
		}
	}

	return context, nil
}
