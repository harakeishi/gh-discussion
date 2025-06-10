package models

import "time"

// Discussion represents a GitHub discussion
type Discussion struct {
	ID                  string             `json:"id"`
	Number              int                `json:"number"`
	Title               string             `json:"title"`
	Body                string             `json:"body"`
	BodyText            string             `json:"bodyText"`
	BodyHTML            string             `json:"bodyHTML"`
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
	PublishedAt         *time.Time         `json:"publishedAt"`
	LastEditedAt        *time.Time         `json:"lastEditedAt"`
	Author              *User              `json:"author"`
	Category            *Category          `json:"category"`
	Repository          *Repository        `json:"repository"`
	URL                 string             `json:"url"`
	ResourcePath        string             `json:"resourcePath"`
	Locked              bool               `json:"locked"`
	ActiveLockReason    *string            `json:"activeLockReason"`
	AnswerChosenAt      *time.Time         `json:"answerChosenAt"`
	AnswerChosenBy      *User              `json:"answerChosenBy"`
	Answer              *Comment           `json:"answer"`
	IsAnswered          bool               `json:"isAnswered"`
	Comments            *CommentConnection `json:"comments"`
	Labels              *LabelConnection   `json:"labels"`
	ReactionGroups      []ReactionGroup    `json:"reactionGroups"`
	ViewerCanDelete     bool               `json:"viewerCanDelete"`
	ViewerCanReact      bool               `json:"viewerCanReact"`
	ViewerCanSubscribe  bool               `json:"viewerCanSubscribe"`
	ViewerCanUpdate     bool               `json:"viewerCanUpdate"`
	ViewerDidAuthor     bool               `json:"viewerDidAuthor"`
	ViewerSubscription  string             `json:"viewerSubscription"`
	AuthorAssociation   string             `json:"authorAssociation"`
	CreatedViaEmail     bool               `json:"createdViaEmail"`
	DatabaseID          int                `json:"databaseId"`
	Editor              *User              `json:"editor"`
	IncludesCreatedEdit bool               `json:"includesCreatedEdit"`
}

// User represents a GitHub user
type User struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	URL       string `json:"url"`
	AvatarURL string `json:"avatarUrl"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

// Category represents a discussion category
type Category struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Emoji        string    `json:"emoji"`
	EmojiHTML    string    `json:"emojiHTML"`
	IsAnswerable bool      `json:"isAnswerable"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Repository represents a GitHub repository
type Repository struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
	Owner         *User  `json:"owner"`
	URL           string `json:"url"`
	Description   string `json:"description"`
}

// Comment represents a discussion comment
type Comment struct {
	ID                      string             `json:"id"`
	Body                    string             `json:"body"`
	BodyText                string             `json:"bodyText"`
	BodyHTML                string             `json:"bodyHTML"`
	CreatedAt               time.Time          `json:"createdAt"`
	UpdatedAt               time.Time          `json:"updatedAt"`
	PublishedAt             *time.Time         `json:"publishedAt"`
	Author                  *User              `json:"author"`
	AuthorAssociation       string             `json:"authorAssociation"`
	IsAnswer                bool               `json:"isAnswer"`
	IsMinimized             bool               `json:"isMinimized"`
	MinimizedReason         *string            `json:"minimizedReason"`
	ReactionGroups          []ReactionGroup    `json:"reactionGroups"`
	Replies                 *CommentConnection `json:"replies"`
	ReplyTo                 *Comment           `json:"replyTo"`
	URL                     string             `json:"url"`
	ViewerCanMarkAsAnswer   bool               `json:"viewerCanMarkAsAnswer"`
	ViewerCanUnmarkAsAnswer bool               `json:"viewerCanUnmarkAsAnswer"`
}

// CommentConnection represents a paginated list of comments
type CommentConnection struct {
	Nodes      []Comment `json:"nodes"`
	PageInfo   PageInfo  `json:"pageInfo"`
	TotalCount int       `json:"totalCount"`
}

// Label represents a GitHub label
type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// LabelConnection represents a paginated list of labels
type LabelConnection struct {
	Nodes []Label `json:"nodes"`
}

// ReactionGroup represents a group of reactions
type ReactionGroup struct {
	Content string `json:"content"`
	Users   struct {
		TotalCount int `json:"totalCount"`
	} `json:"users"`
}

// PageInfo represents pagination information
type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

// DiscussionConnection represents a paginated list of discussions
type DiscussionConnection struct {
	Nodes    []Discussion `json:"nodes"`
	PageInfo PageInfo     `json:"pageInfo"`
}

// SearchResult represents search results
type SearchResult struct {
	Nodes    []Discussion `json:"nodes"`
	PageInfo PageInfo     `json:"pageInfo"`
}

// ListOptions represents options for listing discussions
type ListOptions struct {
	Owner    string
	Repo     string
	Author   string
	Search   string
	Category string
	Answered *bool
	Limit    int
	After    string
	Labels   []string
}

// ViewOptions represents options for viewing a discussion
type ViewOptions struct {
	Owner        string
	Repo         string
	Number       int
	ShowComments bool
}
