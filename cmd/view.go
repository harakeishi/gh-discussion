package main

import (
    "context"
    "errors"
    "fmt"
    "github.com/spf13/cobra"

    "github.com/example/gh-discussion/internal/api"
)


func newViewCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "view <url|id>",
        Short: "View a discussion",
        Args: func(cmd *cobra.Command, args []string) error {
            if len(args) < 1 {
                return errors.New("discussion URL or ID required")
            }
            return nil
        },
        RunE: func(cmd *cobra.Command, args []string) error {
            d, err := getDiscussion(args[0])
            if err != nil {
                return err
            }
            fmt.Printf("%s\nby %s\n%s\n", d.Title, d.Author, d.Body)
            return nil
        },
    }
    return cmd
}

// getDiscussion fetches discussion detail
func getDiscussion(idOrURL string) (api.DiscussionDetail, error) {
    client := api.NewClient()
    return client.GetDiscussion(context.Background(), idOrURL)
}

