package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createOptions holds the options for the create command
type createOptions struct {
	repo     string
	title    string
	body     string
	category string
	web      bool
}

// NewCreateCmd creates the create command
func NewCreateCmd() *cobra.Command {
	opts := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new discussion",
		Long: `Create a new discussion in a repository.

This command will open an interactive prompt to gather the required information
for creating a discussion, including title, body, and category.`,
		Example: `  # Create a discussion interactively
  gh discussion create

  # Create a discussion with title and body
  gh discussion create --title "Discussion Title" --body "Discussion body"

  # Create a discussion in a specific category
  gh discussion create --category "General"

  # Create a discussion in a specific repository
  gh discussion create -R owner/repo

  # Open the discussion creation form in web browser
  gh discussion create -w`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(opts)
		},
	}

	// Repository options
	cmd.Flags().StringVarP(&opts.repo, "repo", "R", "", "Select another repository using the [HOST/]OWNER/REPO format")

	// Discussion options
	cmd.Flags().StringVar(&opts.title, "title", "", "Title for the discussion")
	cmd.Flags().StringVar(&opts.body, "body", "", "Body for the discussion")
	cmd.Flags().StringVar(&opts.category, "category", "", "Category for the discussion")

	// Output options
	cmd.Flags().BoolVarP(&opts.web, "web", "w", false, "Open the discussion creation form in the web browser")

	return cmd
}

// runCreate executes the create command
func runCreate(opts *createOptions) error {
	// Parse repository
	repo, err := parseRepository(opts.repo)
	if err != nil {
		return fmt.Errorf("failed to parse repository: %w", err)
	}

	// Handle web browser option
	if opts.web {
		return openInBrowser(fmt.Sprintf("https://github.com/%s/%s/discussions/new", repo.Owner, repo.Name))
	}

	// For now, just return a message indicating this feature is not yet implemented
	fmt.Println("Discussion creation via CLI is not yet implemented.")
	fmt.Printf("You can create a discussion in the web browser at: https://github.com/%s/%s/discussions/new\n", repo.Owner, repo.Name)

	return nil
}
