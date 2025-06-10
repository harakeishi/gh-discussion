package client

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/harakeishi/gh-discussion/pkg/models"
)

// GitHubClient wraps the GitHub GraphQL API client
type GitHubClient struct {
	client *api.GraphQLClient
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient() (*GitHubClient, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GraphQL client: %w", err)
	}

	return &GitHubClient{
		client: client,
	}, nil
}

// ListDiscussions retrieves a list of discussions based on the provided options
func (c *GitHubClient) ListDiscussions(opts models.ListOptions) (*models.DiscussionConnection, error) {
	// Use search API if search term or author filter is specified
	if opts.Search != "" || opts.Author != "" {
		return c.searchDiscussions(opts)
	}
	return c.listRepositoryDiscussions(opts)
}

// listRepositoryDiscussions lists discussions in a specific repository
func (c *GitHubClient) listRepositoryDiscussions(opts models.ListOptions) (*models.DiscussionConnection, error) {
	query := `
		query ListDiscussions($owner: String!, $repo: String!, $first: Int!, $after: String, $orderBy: DiscussionOrder, $categoryId: ID, $answered: Boolean) {
			repository(owner: $owner, name: $repo) {
				discussions(first: $first, after: $after, orderBy: $orderBy, categoryId: $categoryId, answered: $answered) {
					pageInfo {
						hasNextPage
						endCursor
					}
					nodes {
						id
						number
						title
						bodyText
						createdAt
						updatedAt
						author {
							login
							url
						}
						category {
							name
						}
						url
						answerChosenAt
						isAnswered
						comments(first: 0) {
							totalCount
						}
						labels(first: 10) {
							nodes {
								name
								color
							}
						}
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"owner": opts.Owner,
		"repo":  opts.Repo,
		"first": opts.Limit,
	}

	if opts.After != "" {
		variables["after"] = opts.After
	}

	// Set default ordering
	variables["orderBy"] = map[string]interface{}{
		"field":     "UPDATED_AT",
		"direction": "DESC",
	}

	// Add answered filter if specified
	if opts.Answered != nil {
		variables["answered"] = *opts.Answered
	}

	// Add category filter if specified
	if opts.Category != "" {
		categoryID, err := c.getCategoryID(opts.Owner, opts.Repo, opts.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to get category ID: %w", err)
		}
		if categoryID != "" {
			variables["categoryId"] = categoryID
		}
	}

	var response struct {
		Repository struct {
			Discussions models.DiscussionConnection `json:"discussions"`
		} `json:"repository"`
	}

	err := c.client.Do(query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list discussions: %w", err)
	}

	return &response.Repository.Discussions, nil
}

// searchDiscussions searches for discussions using GitHub's search API
func (c *GitHubClient) searchDiscussions(opts models.ListOptions) (*models.DiscussionConnection, error) {
	query := `
		query SearchDiscussions($query: String!, $first: Int!, $after: String) {
			search(type: DISCUSSION, query: $query, first: $first, after: $after) {
				pageInfo {
					hasNextPage
					endCursor
				}
				nodes {
					... on Discussion {
						id
						number
						title
						bodyText
						createdAt
						updatedAt
						author {
							login
							url
						}
						category {
							name
						}
						repository {
							nameWithOwner
						}
						url
						answerChosenAt
						isAnswered
						comments(first: 0) {
							totalCount
						}
						labels(first: 10) {
							nodes {
								name
								color
							}
						}
					}
				}
			}
		}`

	searchQuery := c.buildSearchQuery(opts)

	variables := map[string]interface{}{
		"query": searchQuery,
		"first": opts.Limit,
	}

	if opts.After != "" {
		variables["after"] = opts.After
	}

	var response struct {
		Search models.SearchResult `json:"search"`
	}

	err := c.client.Do(query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to search discussions: %w", err)
	}

	return &models.DiscussionConnection{
		Nodes:    response.Search.Nodes,
		PageInfo: response.Search.PageInfo,
	}, nil
}

// buildSearchQuery constructs a search query string based on the options
func (c *GitHubClient) buildSearchQuery(opts models.ListOptions) string {
	var parts []string

	// Add repository filter
	if opts.Owner != "" && opts.Repo != "" {
		parts = append(parts, fmt.Sprintf("repo:%s/%s", opts.Owner, opts.Repo))
	}

	// Add search term
	if opts.Search != "" {
		parts = append(parts, opts.Search)
	}

	// Add author filter
	if opts.Author != "" {
		parts = append(parts, fmt.Sprintf("author:%s", opts.Author))
	}

	// Add category filter
	if opts.Category != "" {
		parts = append(parts, fmt.Sprintf("category:\"%s\"", opts.Category))
	}

	// Add answered/unanswered filter
	if opts.Answered != nil {
		if *opts.Answered {
			parts = append(parts, "is:answered")
		} else {
			parts = append(parts, "is:unanswered")
		}
	}

	// Add label filters
	for _, label := range opts.Labels {
		parts = append(parts, fmt.Sprintf("label:\"%s\"", label))
	}

	return strings.Join(parts, " ")
}

// GetDiscussion retrieves a specific discussion by number
func (c *GitHubClient) GetDiscussion(opts models.ViewOptions) (*models.Discussion, error) {
	query := `
		query GetDiscussion($owner: String!, $repo: String!, $number: Int!, $includeComments: Boolean!) {
			repository(owner: $owner, name: $repo) {
				discussion(number: $number) {
					id
					number
					title
					body
					bodyText
					bodyHTML
					createdAt
					updatedAt
					publishedAt
					lastEditedAt
					author {
						login
						url
						avatarUrl
						... on User {
							name
							email
						}
					}
					category {
						name
						description
						emoji
						emojiHTML
						isAnswerable
					}
					repository {
						nameWithOwner
						url
						description
					}
					labels(first: 10) {
						nodes {
							name
							color
						}
					}
					url
					resourcePath
					locked
					activeLockReason
					answerChosenAt
					answerChosenBy {
						login
						url
						... on User {
							name
							email
						}
					}
					answer {
						id
						body
						bodyText
						createdAt
						author {
							login
							url
							... on User {
								name
								email
							}
						}
						isAnswer
					}
					isAnswered
					upvoteCount
					reactionGroups {
						content
						users {
							totalCount
						}
					}
					viewerCanDelete
					viewerCanReact
					viewerCanSubscribe
					viewerCanUpdate
					viewerDidAuthor
					viewerSubscription
					authorAssociation
					createdViaEmail
					databaseId
					editor {
						login
						url
						... on User {
							name
							email
						}
					}
					includesCreatedEdit
					comments(first: 100) @include(if: $includeComments) {
						totalCount
						pageInfo {
							hasNextPage
							endCursor
						}
						nodes {
							id
							body
							bodyText
							bodyHTML
							createdAt
							updatedAt
							publishedAt
							author {
								login
								url
								avatarUrl
								... on User {
									name
									email
								}
							}
							authorAssociation
							upvoteCount
							isAnswer
							isMinimized
							minimizedReason
							reactionGroups {
								content
								users {
									totalCount
								}
							}
							url
							viewerCanMarkAsAnswer
							viewerCanUnmarkAsAnswer
							replies(first: 50) {
								totalCount
								nodes {
									id
									body
									bodyText
									createdAt
									updatedAt
									author {
										login
										url
										avatarUrl
										... on User {
											name
											email
										}
									}
									authorAssociation
									isAnswer
									url
								}
							}
						}
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"owner":           opts.Owner,
		"repo":            opts.Repo,
		"number":          opts.Number,
		"includeComments": opts.ShowComments,
	}

	var response struct {
		Repository struct {
			Discussion *models.Discussion `json:"discussion"`
		} `json:"repository"`
	}

	err := c.client.Do(query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get discussion: %w", err)
	}

	if response.Repository.Discussion == nil {
		return nil, fmt.Errorf("discussion #%d not found", opts.Number)
	}

	return response.Repository.Discussion, nil
}

// GetRepositoryInfo retrieves basic repository information
func (c *GitHubClient) GetRepositoryInfo(owner, repo string) (*models.Repository, error) {
	query := `
		query GetRepository($owner: String!, $repo: String!) {
			repository(owner: $owner, name: $repo) {
				id
				name
				nameWithOwner
				owner {
					login
					url
				}
				url
				description
			}
		}`

	variables := map[string]interface{}{
		"owner": owner,
		"repo":  repo,
	}

	var response struct {
		Repository *models.Repository `json:"repository"`
	}

	err := c.client.Do(query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info: %w", err)
	}

	if response.Repository == nil {
		return nil, fmt.Errorf("repository %s/%s not found", owner, repo)
	}

	return response.Repository, nil
}

// GetDiscussionCategories retrieves available discussion categories for a repository
func (c *GitHubClient) GetDiscussionCategories(owner, repo string) ([]models.Category, error) {
	query := `
		query GetDiscussionCategories($owner: String!, $repo: String!) {
			repository(owner: $owner, name: $repo) {
				discussionCategories(first: 100) {
					nodes {
						id
						name
						description
						emoji
						emojiHTML
						isAnswerable
						createdAt
						updatedAt
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"owner": owner,
		"repo":  repo,
	}

	var response struct {
		Repository struct {
			DiscussionCategories struct {
				Nodes []models.Category `json:"nodes"`
			} `json:"discussionCategories"`
		} `json:"repository"`
	}

	err := c.client.Do(query, variables, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get discussion categories: %w", err)
	}

	return response.Repository.DiscussionCategories.Nodes, nil
}

// getCategoryID retrieves the ID of a discussion category by name
func (c *GitHubClient) getCategoryID(owner, repo, categoryName string) (string, error) {
	categories, err := c.GetDiscussionCategories(owner, repo)
	if err != nil {
		return "", err
	}

	for _, category := range categories {
		if category.Name == categoryName {
			return category.ID, nil
		}
	}

	return "", nil // Category not found, but not an error
}
