package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"git.pepabo.com/harachan/gh-discussion/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gh-discussion",
		Short: "GitHub CLI extension for managing discussions",
		Long: `A GitHub CLI extension for managing discussions.

This extension provides commands to list, view, and create discussions
in GitHub repositories, similar to how gh issue and gh pr work.`,
		Example: `  # List discussions in the current repository
  gh discussion list

  # View a specific discussion
  gh discussion view 123

  # Create a new discussion
  gh discussion create`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	rootCmd.AddCommand(cmd.NewListCmd())
	rootCmd.AddCommand(cmd.NewViewCmd())
	rootCmd.AddCommand(cmd.NewCreateCmd())

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
