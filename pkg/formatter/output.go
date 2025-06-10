package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"git.pepabo.com/harachan/gh-discussion/pkg/models"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatTable    OutputFormat = "table"
	FormatJSON     OutputFormat = "json"
	FormatTemplate OutputFormat = "template"
)

// OutputOptions contains options for formatting output
type OutputOptions struct {
	Format    OutputFormat
	Fields    []string
	Template  string
	JQFilter  string
	ColorMode string
}

// Formatter handles output formatting
type Formatter struct {
	writer io.Writer
	opts   OutputOptions
}

// NewFormatter creates a new formatter
func NewFormatter(writer io.Writer, opts OutputOptions) *Formatter {
	return &Formatter{
		writer: writer,
		opts:   opts,
	}
}

// FormatDiscussionList formats a list of discussions
func (f *Formatter) FormatDiscussionList(discussions []models.Discussion) error {
	switch f.opts.Format {
	case FormatJSON:
		return f.formatDiscussionListJSON(discussions)
	case FormatTemplate:
		return f.formatDiscussionListTemplate(discussions)
	default:
		return f.formatDiscussionListTable(discussions)
	}
}

// FormatDiscussion formats a single discussion
func (f *Formatter) FormatDiscussion(discussion *models.Discussion) error {
	switch f.opts.Format {
	case FormatJSON:
		return f.formatDiscussionJSON(discussion)
	case FormatTemplate:
		return f.formatDiscussionTemplate(discussion)
	default:
		return f.formatDiscussionTable(discussion)
	}
}

// formatDiscussionListTable formats discussions as a table
func (f *Formatter) formatDiscussionListTable(discussions []models.Discussion) error {
	if len(discussions) == 0 {
		fmt.Fprintln(f.writer, "No discussions found")
		return nil
	}

	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Header
	fmt.Fprintln(w, "NUMBER\tTITLE\tAUTHOR\tCATEGORY\tANSWERED\tCOMMENTS\tUPDATED")

	for _, discussion := range discussions {
		author := ""
		if discussion.Author != nil {
			author = discussion.Author.Login
		}

		category := ""
		if discussion.Category != nil {
			category = discussion.Category.Name
		}

		answered := "No"
		if discussion.IsAnswered {
			answered = "Yes"
		}

		comments := "0"
		if discussion.Comments != nil {
			comments = strconv.Itoa(discussion.Comments.TotalCount)
		}

		updated := f.formatTime(discussion.UpdatedAt)

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
			discussion.Number,
			f.truncateString(discussion.Title, 50),
			author,
			category,
			answered,
			comments,
			updated,
		)
	}

	return nil
}

// formatDiscussionTable formats a single discussion as a table
func (f *Formatter) formatDiscussionTable(discussion *models.Discussion) error {
	fmt.Fprintf(f.writer, "Discussion #%d\n", discussion.Number)
	fmt.Fprintf(f.writer, "Title: %s\n", discussion.Title)

	if discussion.Author != nil {
		fmt.Fprintf(f.writer, "Author: %s\n", discussion.Author.Login)
	}

	if discussion.Category != nil {
		fmt.Fprintf(f.writer, "Category: %s\n", discussion.Category.Name)
	}

	if discussion.Repository != nil {
		fmt.Fprintf(f.writer, "Repository: %s\n", discussion.Repository.NameWithOwner)
	}

	fmt.Fprintf(f.writer, "Created: %s\n", f.formatTime(discussion.CreatedAt))
	fmt.Fprintf(f.writer, "Updated: %s\n", f.formatTime(discussion.UpdatedAt))
	fmt.Fprintf(f.writer, "Answered: %t\n", discussion.IsAnswered)

	if discussion.Comments != nil {
		fmt.Fprintf(f.writer, "Comments: %d\n", discussion.Comments.TotalCount)
	}

	fmt.Fprintf(f.writer, "URL: %s\n", discussion.URL)

	if discussion.Labels != nil && len(discussion.Labels.Nodes) > 0 {
		labels := make([]string, len(discussion.Labels.Nodes))
		for i, label := range discussion.Labels.Nodes {
			labels[i] = label.Name
		}
		fmt.Fprintf(f.writer, "Labels: %s\n", strings.Join(labels, ", "))
	}

	if discussion.Body != "" {
		fmt.Fprintf(f.writer, "\n%s\n", discussion.Body)
	}

	// Show comments if available
	if discussion.Comments != nil && len(discussion.Comments.Nodes) > 0 {
		fmt.Fprintf(f.writer, "\n--- Comments ---\n")
		for i, comment := range discussion.Comments.Nodes {
			fmt.Fprintf(f.writer, "\nComment #%d", i+1)
			if comment.IsAnswer {
				fmt.Fprintf(f.writer, " (Answer)")
			}
			fmt.Fprintf(f.writer, "\n")

			if comment.Author != nil {
				fmt.Fprintf(f.writer, "Author: %s\n", comment.Author.Login)
			}
			fmt.Fprintf(f.writer, "Created: %s\n", f.formatTime(comment.CreatedAt))
			fmt.Fprintf(f.writer, "\n%s\n", comment.Body)
		}
	}

	return nil
}

// formatDiscussionListJSON formats discussions as JSON
func (f *Formatter) formatDiscussionListJSON(discussions []models.Discussion) error {
	if len(f.opts.Fields) > 0 {
		filtered := f.filterFields(discussions, f.opts.Fields)
		return json.NewEncoder(f.writer).Encode(filtered)
	}
	return json.NewEncoder(f.writer).Encode(discussions)
}

// formatDiscussionJSON formats a single discussion as JSON
func (f *Formatter) formatDiscussionJSON(discussion *models.Discussion) error {
	if len(f.opts.Fields) > 0 {
		filtered := f.filterFields(discussion, f.opts.Fields)
		return json.NewEncoder(f.writer).Encode(filtered)
	}
	return json.NewEncoder(f.writer).Encode(discussion)
}

// formatDiscussionListTemplate formats discussions using a template
func (f *Formatter) formatDiscussionListTemplate(discussions []models.Discussion) error {
	tmpl, err := template.New("discussions").Parse(f.opts.Template)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl.Execute(f.writer, discussions)
}

// formatDiscussionTemplate formats a single discussion using a template
func (f *Formatter) formatDiscussionTemplate(discussion *models.Discussion) error {
	tmpl, err := template.New("discussion").Parse(f.opts.Template)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl.Execute(f.writer, discussion)
}

// filterFields filters the data to include only specified fields
func (f *Formatter) filterFields(data interface{}, fields []string) interface{} {
	// Convert to JSON and back to get a map representation
	jsonData, err := json.Marshal(data)
	if err != nil {
		return data
	}

	var result interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return data
	}

	return f.filterFieldsRecursive(result, fields)
}

// filterFieldsRecursive recursively filters fields from the data
func (f *Formatter) filterFieldsRecursive(data interface{}, fields []string) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		filtered := make(map[string]interface{})
		for _, field := range fields {
			if value, exists := v[field]; exists {
				filtered[field] = value
			}
		}
		return filtered
	case []interface{}:
		filtered := make([]interface{}, len(v))
		for i, item := range v {
			filtered[i] = f.filterFieldsRecursive(item, fields)
		}
		return filtered
	default:
		return data
	}
}

// formatTime formats a time value for display
func (f *Formatter) formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 2, 2006")
	}
}

// truncateString truncates a string to the specified length
func (f *Formatter) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// GetAvailableFields returns the available fields for JSON output
func GetAvailableFields() map[string][]string {
	return map[string][]string{
		"discussion": {
			"activeLockReason", "answer", "answerChosenAt", "answerChosenBy", "author", "authorAssociation",
			"body", "bodyHTML", "bodyText", "category", "comments", "createdAt", "createdViaEmail", "databaseId",
			"editor", "id", "includesCreatedEdit", "isAnswered", "lastEditedAt", "locked", "number",
			"publishedAt", "reactionGroups", "reactions", "repository", "resourcePath", "title", "updatedAt",
			"url", "userContentEdits", "viewerCanDelete", "viewerCanReact", "viewerCanSubscribe",
			"viewerCanUpdate", "viewerDidAuthor", "viewerSubscription",
		},
		"author": {
			"avatarUrl", "login", "url", "id", "name", "email",
		},
		"category": {
			"id", "name", "description", "emoji", "emojiHTML", "isAnswerable", "createdAt", "updatedAt",
		},
		"comments": {
			"author", "authorAssociation", "body", "bodyHTML", "bodyText", "createdAt", "id", "isAnswer",
			"isMinimized", "minimizedReason", "publishedAt", "reactionGroups", "replies", "replyTo",
			"updatedAt", "url", "viewerCanMarkAsAnswer", "viewerCanUnmarkAsAnswer",
		},
		"repository": {
			"id", "name", "nameWithOwner", "owner", "url", "description",
		},
	}
}
