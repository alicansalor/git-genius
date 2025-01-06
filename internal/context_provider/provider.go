package context_provider

type ContextProviderType string

const (
	LinearContextProviderType     ContextProviderType = "linear"
	GitContextProviderType        ContextProviderType = "git"
	PRTemplateContextProviderType ContextProviderType = "pr_template"
)

type ContextProvider interface {
	FetchContext() (ProvidedContext, error)
}

type ProvidedContext interface {
	IsEmpty() bool
}
