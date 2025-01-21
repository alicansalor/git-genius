package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func prCmd(dep *SharedDependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "pr",
		Short: "Generate a pull request",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			// Generate pull request
			prContent, err := dep.sdk.GeneratePullRequestContent(ctx)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// print the title and body of the PR
			fmt.Printf("Generated Pull Request:\nTitle: %s\nBody: %s\n",
				prContent.Title, prContent.Body)
		},
	}
}
