package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	llm "git-genius/internal/llm"
	versioncontrol "git-genius/internal/version_control"
)

type Config struct {
	ContextProviders []ProviderConfig     `yaml:"context_providers"`
	IssueID          string               // dynamically set
	LLM              LLMConfig            `yaml:"llm"`
	VersionControl   VersionControlConfig `yaml:"version_control"`
}

type VersionControlConfig struct {
	Provider string `yaml:"provider"`
	Token    string `yaml:"token"`
}

type LLMConfig struct {
	Name   string `yaml:"name"`
	APIKey string `yaml:"api_key"`
}

type ProviderConfig struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	APIKey string `yaml:"api_key"`
}

// DefaultConfigPath returns the default configuration file path
func DefaultConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to determine user config directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(configDir, "git-genius", "config.yaml")
}

// LoadConfig loads the configuration from the specified path
func LoadConfig(path string) (*Config, error) {

	if path == "" {
		path = DefaultConfigPath()
	}

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at path: %s", path)
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Decode YAML
	var cfg Config
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

func (cfg Config) NewLLM(ctx context.Context) (llm.LLM, error) {
	switch cfg.LLM.Name {
	case "gemini":
		return llm.NewGemini(ctx, cfg.LLM.APIKey)
	default:
		return nil, fmt.Errorf("unsupported LLM: %s", cfg.LLM.Name)
	}
}

// NewPRCreator generates the appropriate PRCreator based on the configuration
func (cfg Config) NewPRCreator(ctx context.Context) (versioncontrol.PRCreator, error) {

	versionControl := cfg.VersionControl

	switch versionControl.Provider {
	case "github":
		if versionControl.Token == "" {
			return nil, fmt.Errorf("GitHub token is required for GitHub PR creation")
		}
		return versioncontrol.NewGitHubManager(ctx, versionControl.Token)
	default:
		return nil, fmt.Errorf("unsupported context_provider: %s", versionControl.Provider)
	}
}
