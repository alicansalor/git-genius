package cmd

import (
	"context"
	"fmt"
	"git-genius/config"
	"git-genius/sdk"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command
var RootCmd = &cobra.Command{
	Use:                "git-genius",
	Short:              "Git Genius enhances your Git workflow with LLM-powered features",
	Long:               "Git Genius enhances your Git workflow with features like automated commit messages and PR generation using an LLM.",
	Args:               cobra.ArbitraryArgs, // Allow any arguments
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help() // Show help if no arguments are provided.
		}
		// Pass any unrecognized commands to Git.
		return runGitCommand(args)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		cmd.SetContext(ctx)

		if cmd.HasSubCommands() && len(args) == 0 {
			return nil
		}

		// if the command is not recognised then nothing to initialise
		if len(args) > 0 {
			subCmd, _, _ := cmd.Root().Find(args)
			if subCmd == nil || subCmd == cmd.Root() {
				return nil
			}
		}

		// Get the config path from the flag
		configPath, _ := cmd.Flags().GetString("config")
		// Initialize shared dependencies
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %v\n", err)
		}

		// Add issue ID to the config
		issueID, _ := cmd.Flags().GetString("issue")
		cfg.IssueID = issueID

		sharedDeps.sdk, err = sdk.NewGitGeniusSDK(ctx, cfg)
		if err != nil {
			return fmt.Errorf("failed to create SDK: %v", err)
		}

		return nil
	},
}

// SharedDependencies holds dependencies shared between subcommands
type SharedDependencies struct {
	sdk sdk.GitGenius
}

var sharedDeps SharedDependencies

// Execute runs the root command
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	// Define persistent flags for the root command
	RootCmd.PersistentFlags().String("config", "", "Path to the configuration file")
	RootCmd.PersistentFlags().String("issue", "", "Issue ID to associate with the operation")

	// Add subcommands and pass shared dependencies
	RootCmd.AddCommand(prCmd(&sharedDeps))
	RootCmd.AddCommand(commitCmd(&sharedDeps))
}

// runGitCommand forwards unrecognized commands to the Git CLI
func runGitCommand(args []string) error {
	fmt.Printf("Running fallback Git command with args: %v\n", args)

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute git command: %v", err)
	}
	return nil
}
