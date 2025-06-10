package main

import (
    "os"
    "github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "gh-discussion",
        Short: "Interact with GitHub Discussions",
    }

    cmd.AddCommand(newSearchCmd())
    cmd.AddCommand(newViewCmd())
    return cmd
}

func main() {
    if err := newRootCmd().Execute(); err != nil {
        os.Exit(1)
    }
}

