package main

import (
	"github.com/spf13/cobra"
	"os"
)

var repo string

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gh-discussion",
		Short: "Interact with GitHub Discussions",
	}

	cmd.PersistentFlags().StringVar(&repo, "repo", "", "target repository (owner/repo)")

	cmd.AddCommand(newSearchCmd())
	cmd.AddCommand(newViewCmd())
	return cmd
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
