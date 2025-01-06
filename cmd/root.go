package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the base command
var RootCmd = &cobra.Command{
	Use:   "git-genius",
	Short: "Git Genius enhances your Git workflow with LLM-powered features",
	Long:  "Git Genius enhances your Git workflow with features like automated commit messages and PR generation using an LLM.",
}

// Execute adds all child commands to the root command
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	// Add subcommands
	RootCmd.AddCommand(commitCmd)
	RootCmd.AddCommand(prCmd)
}
