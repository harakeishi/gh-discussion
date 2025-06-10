package api

import (
    "context"
)

// Client abstracts GitHub API interactions
// In real implementation this would hold authentication info and GraphQL client

type Client struct{}

func NewClient() *Client {
    return &Client{}
}

func (c *Client) SearchDiscussions(ctx context.Context, query string) ([]Discussion, error) {
    // TODO: implement actual GraphQL query
    return nil, nil
}

func (c *Client) GetDiscussion(ctx context.Context, idOrURL string) (DiscussionDetail, error) {
    // TODO: implement actual GraphQL query
    return DiscussionDetail{}, nil
}

// Reuse types from main package
