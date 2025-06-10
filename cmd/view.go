package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/harakeishi/gh-discussion/pkg/client"
	"github.com/harakeishi/gh-discussion/pkg/formatter"
	"github.com/harakeishi/gh-discussion/pkg/models"
)

// viewOptions holds the options for the view command
type viewOptions struct {
	repo     string
	comments bool
	json     string
	template string
	web      bool
}

// NewViewCmd creates the view command
func NewViewCmd() *cobra.Command {
	opts := &viewOptions{}

	cmd := &cobra.Command{
		Use:   "view {<number> | <url>}",
		Short: "Display the title, body, and other information about a discussion",
		Long: `Display the title, body, and other information about a discussion.

With '--comments', view discussion comments.
With '--web', open the discussion in a web browser instead.`,
		Example: `  # View discussion #123 in the current repository
  gh discussion view 123

  # View discussion #123 in a specific repository
  gh discussion view 123 -R owner/repo

  # View discussion by URL
  gh discussion view https://github.com/owner/repo/discussions/123

  # View discussion with comments
  gh discussion view 123 -c

  # View discussion as JSON
  gh discussion view 123 --json

  # View specific fields as JSON
  gh discussion view 123 --json "title,body,author,comments"

  # Use a custom template
  gh discussion view 123 --template '{{.title}} by {{.author.login}}'

  # Open in web browser
  gh discussion view 123 -w`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runView(opts, args[0])
		},
	}

	// Repository options
	cmd.Flags().StringVarP(&opts.repo, "repo", "R", "", "Select another repository using the [HOST/]OWNER/REPO format")

	// Display options
	cmd.Flags().BoolVarP(&opts.comments, "comments", "c", false, "View discussion comments")

	// Output options
	cmd.Flags().StringVar(&opts.json, "json", "", "Output JSON with the specified fields")
	cmd.Flags().StringVar(&opts.template, "template", "", "Format JSON output using a Go template")
	cmd.Flags().BoolVarP(&opts.web, "web", "w", false, "Open the discussion in the web browser")

	// Mark mutually exclusive flags
	cmd.MarkFlagsMutuallyExclusive("json", "template", "web")

	return cmd
}

// runView executes the view command
func runView(opts *viewOptions, discussionArg string) error {
	// Parse discussion argument (number or URL)
	repo, number, err := parseDiscussionArg(discussionArg, opts.repo)
	if err != nil {
		return fmt.Errorf("failed to parse discussion argument: %w", err)
	}

	// Handle web browser option
	if opts.web {
		return openInBrowser(fmt.Sprintf("https://github.com/%s/%s/discussions/%d", repo.Owner, repo.Name, number))
	}

	// Create GitHub client
	client, err := client.NewGitHubClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Build view options
	viewOpts := models.ViewOptions{
		Owner:        repo.Owner,
		Repo:         repo.Name,
		Number:       number,
		ShowComments: opts.comments,
	}

	// Fetch discussion
	discussion, err := client.GetDiscussion(viewOpts)
	if err != nil {
		return fmt.Errorf("failed to get discussion: %w", err)
	}

	// Determine output format
	outputOpts := formatter.OutputOptions{
		Format: formatter.FormatTable,
	}

	if opts.json != "" {
		outputOpts.Format = formatter.FormatJSON
		if opts.json != "true" && opts.json != "1" {
			outputOpts.Fields = parseJSONFields(opts.json)
		}
	} else if opts.template != "" {
		outputOpts.Format = formatter.FormatTemplate
		outputOpts.Template = opts.template
	}

	// Format and output result
	f := formatter.NewFormatter(os.Stdout, outputOpts)
	return f.FormatDiscussion(discussion)
}

// parseDiscussionArg parses the discussion argument which can be a number or URL
func parseDiscussionArg(arg, repoStr string) (*Repository, int, error) {
	// Check if it's a URL
	if strings.HasPrefix(arg, "https://github.com/") {
		return parseDiscussionURL(arg)
	}

	// Parse as discussion number
	number, err := strconv.Atoi(arg)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid discussion number: %s", arg)
	}

	// Parse repository
	repo, err := parseRepository(repoStr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse repository: %w", err)
	}

	return repo, number, nil
}

// parseDiscussionURL parses a GitHub discussion URL
func parseDiscussionURL(url string) (*Repository, int, error) {
	// Expected format: https://github.com/owner/repo/discussions/123
	url = strings.TrimPrefix(url, "https://github.com/")
	parts := strings.Split(url, "/")

	if len(parts) != 4 || parts[2] != "discussions" {
		return nil, 0, fmt.Errorf("invalid discussion URL format")
	}

	owner := parts[0]
	repo := parts[1]

	number, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, 0, fmt.Errorf("invalid discussion number in URL: %s", parts[3])
	}

	return &Repository{
		Owner: owner,
		Name:  repo,
	}, number, nil
}
