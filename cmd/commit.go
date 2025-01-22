package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func commitCmd(dep *SharedDependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "smart-commit",
		Short: "Generate a commit message and commit changes",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			// Generate commit message
			commitMessage, err := dep.sdk.GenerateCommitMessage(ctx)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Show the generated commit message
			fmt.Println("\nGenerated Commit Message:")
			fmt.Println("-------------------------")
			fmt.Println(commitMessage)
			fmt.Println("-------------------------")

			// Prompt user for input
			fmt.Println("What would you like to do?")
			fmt.Println("[Y] Accept and Commit")
			fmt.Println("[N] Cancel")

			var choice string
			fmt.Print("Enter your choice (Y/N): ")
			fmt.Scanln(&choice)

			switch strings.ToLower(choice) {
			case "y":
				// Proceed with Git's editor for commit
				performGitCommitWithEditor(commitMessage)
			case "n":
				// Cancel the operation
				fmt.Println("Commit canceled.")
			default:
				fmt.Println("Invalid choice. Commit canceled.")
			}
		},
	}

	// flags
	cmd.PersistentFlags().String("issue", "", "Issue ID to associate with the operation")

	return cmd
}

func performGitCommitWithEditor(commitMessage string) error {
	// Write the message to a temporary file
	tempFile, err := os.CreateTemp("", "git-commit-msg-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(commitMessage)
	if err != nil {
		return fmt.Errorf("failed to write commit message to file: %w", err)
	}
	tempFile.Close()

	// Open Git's editor for editing and commit
	cmd := exec.Command("git", "commit", "--edit", "--file", tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		// If Git aborts due to an empty commit message, notify the user
		fmt.Println("Commit canceled by user.")
		return nil
	}

	fmt.Println("Commit completed successfully.")
	return nil
}
