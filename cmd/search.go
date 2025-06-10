package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"time"

	"github.com/harakeishi/gh-discussion/internal/api"
)

func newSearchCmd() *cobra.Command {
	var (
		fromStr  string
		toStr    string
		user     string
		keywords []string
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search GitHub Discussions",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromStr == "" && toStr == "" && user == "" && len(keywords) == 0 {
				return errors.New("at least one search condition must be specified")
			}

			var from, to *time.Time
			if fromStr != "" {
				t, err := time.Parse("2006-01-02", fromStr)
				if err != nil {
					return fmt.Errorf("invalid from date: %w", err)
				}
				from = &t
			}
			if toStr != "" {
				t, err := time.Parse("2006-01-02", toStr)
				if err != nil {
					return fmt.Errorf("invalid to date: %w", err)
				}
				to = &t
			}

			query := buildSearchQuery(from, to, user, repo, keywords)
			results, err := searchDiscussions(query)
			if err != nil {
				return err
			}
			for _, r := range results {
				fmt.Printf("%s\t%s\t%s\t%d\n", r.Title, r.URL, r.Author, r.Comments)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&fromStr, "from", "", "from date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&toStr, "to", "", "to date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&user, "user", "", "author username")
	cmd.Flags().StringSliceVarP(&keywords, "keyword", "k", nil, "search keywords")
	return cmd
}

func buildSearchQuery(from, to *time.Time, user, repo string, keywords []string) string {
	var parts []string
	if from != nil {
		parts = append(parts, fmt.Sprintf("created:%s..", from.Format("2006-01-02")))
	}
	if to != nil {
		if len(parts) > 0 && !strings.HasSuffix(parts[len(parts)-1], "..") {
			parts[len(parts)-1] += to.Format("2006-01-02")
		} else if len(parts) > 0 {
			parts[len(parts)-1] += to.Format("2006-01-02")
		} else {
			parts = append(parts, fmt.Sprintf("created:..%s", to.Format("2006-01-02")))
		}
	}
	if user != "" {
		parts = append(parts, "author:"+user)
	}
	if repo != "" {
		parts = append(parts, "repo:"+repo)
	}
	for _, kw := range keywords {
		parts = append(parts, kw)
	}
	return strings.Join(parts, " ")
}

// searchDiscussions performs a GraphQL search request
func searchDiscussions(query string) ([]api.Discussion, error) {
	client := api.NewClient()
	return client.SearchDiscussions(context.Background(), query)
}
