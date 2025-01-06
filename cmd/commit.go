package cmd

import (
	"context"
	"fmt"

	"git-genius/config"
	"git-genius/internal/sdk"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message and commit changes",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the config path from the flag
		configPath, _ := cmd.Flags().GetString("config")

		// Initialize SDK
		genius := sdk.NewGeniusSDK()

		// Generate commit message
		message, err := genius.GenerateCommitMessage(context.Background(), configPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Generated Commit Message:\n%s\n", message)
	},
}

func init() {
	commitCmd.Flags().StringP("config", "c", config.DefaultConfigPath(), "Path to configuration file")
}
