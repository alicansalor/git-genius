package cmd

import (
	"context"
	"fmt"

	"git-genius/config"
	"git-genius/internal/sdk"

	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Generate a pull request",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the config path from the flag
		configPath, _ := cmd.Flags().GetString("config")

		// Initialize SDK
		genius := sdk.NewGeniusSDK()

		// Generate pull request
		prURL, err := genius.CreatePullRequest(context.Background(), configPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Pull Request Created:\n%s\n", prURL)
	},
}

func init() {
	prCmd.Flags().StringP("config", "c", config.DefaultConfigPath(), "Path to configuration file")
}
