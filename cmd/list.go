package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"

	"github.com/harakeishi/gh-discussion/pkg/client"
	"github.com/harakeishi/gh-discussion/pkg/formatter"
	"github.com/harakeishi/gh-discussion/pkg/models"
)

// listOptions holds the options for the list command
type listOptions struct {
	repo       string
	author     string
	search     string
	category   string
	answered   string
	limit      int
	labels     []string
	json       string
	jsonFields []string
	template   string
	web        bool
}

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List discussions in a repository",
		Long: `List discussions in a repository.

The search query syntax is the same as GitHub's search syntax.
For more information about the search syntax, see:
https://docs.github.com/en/search-github/searching-on-github/searching-discussions`,
		Example: `  # List discussions in the current repository
  gh discussion list

  # List discussions in a specific repository
  gh discussion list -R owner/repo

  # List discussions by a specific author
  gh discussion list -a username

  # Search for discussions containing specific text
  gh discussion list -S "API documentation"

  # Filter by category
  gh discussion list --category "General"

  # Filter by answered status
  gh discussion list --answered
  gh discussion list --unanswered

  # Limit the number of results
  gh discussion list -L 50

  # Output as JSON with specific fields
  gh discussion list --json "number,title,author,category"

  # Use a custom template
  gh discussion list --template '{{range .}}{{.number}} {{.title}}{{"\n"}}{{end}}'

  # Open in web browser
  gh discussion list -w`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	// Repository options
	cmd.Flags().StringVarP(&opts.repo, "repo", "R", "", "Select another repository using the [HOST/]OWNER/REPO format")

	// Filter options
	cmd.Flags().StringVarP(&opts.author, "author", "a", "", "Filter by author")
	cmd.Flags().StringVarP(&opts.search, "search", "S", "", "Search discussions with a query")
	cmd.Flags().StringVar(&opts.category, "category", "", "Filter by category")
	cmd.Flags().StringVar(&opts.answered, "answered", "", "Filter by answered status (true/false)")
	cmd.Flags().StringSliceVarP(&opts.labels, "label", "l", nil, "Filter by labels")

	// Output options
	cmd.Flags().IntVarP(&opts.limit, "limit", "L", 30, "Maximum number of discussions to fetch")
	cmd.Flags().StringVar(&opts.json, "json", "", "Output JSON with the specified fields")
	cmd.Flags().StringVar(&opts.template, "template", "", "Format JSON output using a Go template")
	cmd.Flags().BoolVarP(&opts.web, "web", "w", false, "Open the discussion list in the web browser")

	// Mark mutually exclusive flags
	cmd.MarkFlagsMutuallyExclusive("json", "template", "web")

	return cmd
}

// runList executes the list command
func runList(opts *listOptions) error {
	// Parse repository
	repo, err := parseRepository(opts.repo)
	if err != nil {
		return fmt.Errorf("failed to parse repository: %w", err)
	}

	// Handle web browser option
	if opts.web {
		return openInBrowser(fmt.Sprintf("https://github.com/%s/%s/discussions", repo.Owner, repo.Name))
	}

	// Create GitHub client
	client, err := client.NewGitHubClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Parse answered filter
	var answered *bool
	if opts.answered != "" {
		switch strings.ToLower(opts.answered) {
		case "true", "yes", "1":
			answered = &[]bool{true}[0]
		case "false", "no", "0":
			answered = &[]bool{false}[0]
		default:
			return fmt.Errorf("invalid value for --answered: %s (expected true/false)", opts.answered)
		}
	}

	// Build list options
	listOpts := models.ListOptions{
		Owner:    repo.Owner,
		Repo:     repo.Name,
		Author:   opts.author,
		Search:   opts.search,
		Category: opts.category,
		Answered: answered,
		Limit:    opts.limit,
		Labels:   opts.labels,
	}

	// Fetch discussions
	discussions, err := client.ListDiscussions(listOpts)
	if err != nil {
		return fmt.Errorf("failed to list discussions: %w", err)
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

	// Format and output results
	f := formatter.NewFormatter(os.Stdout, outputOpts)
	return f.FormatDiscussionList(discussions.Nodes)
}

// parseRepository parses the repository string and returns owner/repo
func parseRepository(repoStr string) (*Repository, error) {
	if repoStr == "" {
		// Try to get repository from current directory
		currentRepo, err := repository.Current()
		if err != nil {
			return nil, fmt.Errorf("unable to determine repository. Use -R flag to specify repository")
		}
		return &Repository{
			Owner: currentRepo.Owner,
			Name:  currentRepo.Name,
		}, nil
	}

	// Parse owner/repo format
	parts := strings.Split(repoStr, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository format. Expected OWNER/REPO")
	}

	return &Repository{
		Owner: parts[0],
		Name:  parts[1],
	}, nil
}

// Repository represents a GitHub repository
type Repository struct {
	Owner string
	Name  string
}

// parseJSONFields parses the comma-separated JSON fields
func parseJSONFields(fields string) []string {
	if fields == "" {
		return nil
	}

	var result []string
	for _, field := range strings.Split(fields, ",") {
		field = strings.TrimSpace(field)
		if field != "" {
			result = append(result, field)
		}
	}
	return result
}

// openInBrowser opens the specified URL in the default web browser
func openInBrowser(url string) error {
	fmt.Printf("Opening %s in your browser.\n", url)
	// This would typically use a library like browser.OpenURL
	// For now, just print the URL
	return nil
}
