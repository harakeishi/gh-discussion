package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/harakeishi/gh-discussion/pkg/models"
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

// tableModel represents the bubbletea table model
type tableModel struct {
	table table.Model
}

func (m tableModel) Init() tea.Cmd { return nil }

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	return m.table.View() + "\n"
}

// formatDiscussionListTable formats discussions as a table
func (f *Formatter) formatDiscussionListTable(discussions []models.Discussion) error {
	if len(discussions) == 0 {
		fmt.Fprintln(f.writer, "No discussions found")
		return nil
	}

	// Define table columns
	columns := []table.Column{
		{Title: "NUMBER", Width: 8},
		{Title: "TITLE", Width: 60},
		{Title: "AUTHOR", Width: 15},
		{Title: "CATEGORY", Width: 15},
		{Title: "ANSWERED", Width: 10},
		{Title: "COMMENTS", Width: 10},
		{Title: "UPDATED", Width: 15},
	}

	// Prepare table rows
	var rows []table.Row
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

		rows = append(rows, table.Row{
			strconv.Itoa(discussion.Number),
			f.truncateString(discussion.Title, 60),
			f.truncateString(author, 15),
			f.truncateString(category, 15),
			answered,
			comments,
			updated,
		})
	}

	// Create table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+2),
	)

	// Apply styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		Bold(true).
		Foreground(lipgloss.Color("15"))
	s.Selected = s.Selected.
		Foreground(lipgloss.NoColor{}).
		Background(lipgloss.NoColor{}).
		Bold(false)
	s.Cell = s.Cell.
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false)
	t.SetStyles(s)

	// Create model and run
	m := tableModel{table: t}

	// For non-interactive output, just render the table view
	fmt.Fprint(f.writer, m.View())

	return nil
}

// formatDiscussionTable formats a single discussion as a table
func (f *Formatter) formatDiscussionTable(discussion *models.Discussion) error {
	// Title styling
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true)

	// Label styling for metadata
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)

	// Value styling for metadata
	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// URL styling
	urlStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	// Answer styling
	answerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	// Section separator
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("─────────────────────────────────────────────────────────────────────────")

	// Title
	fmt.Fprintf(f.writer, "%s\n", titleStyle.Render(fmt.Sprintf("Discussion #%d", discussion.Number)))
	fmt.Fprintf(f.writer, "%s\n\n", titleStyle.Render(discussion.Title))

	// Metadata
	if discussion.Author != nil {
		fmt.Fprintf(f.writer, "%s %s\n", labelStyle.Render("Author:"), valueStyle.Render(discussion.Author.Login))
	}

	if discussion.Category != nil {
		fmt.Fprintf(f.writer, "%s %s\n", labelStyle.Render("Category:"), valueStyle.Render(discussion.Category.Name))
	}

	if discussion.Repository != nil {
		fmt.Fprintf(f.writer, "%s %s\n", labelStyle.Render("Repository:"), valueStyle.Render(discussion.Repository.NameWithOwner))
	}

	fmt.Fprintf(f.writer, "%s %s\n", labelStyle.Render("Created:"), valueStyle.Render(f.formatTime(discussion.CreatedAt)))
	fmt.Fprintf(f.writer, "%s %s\n", labelStyle.Render("Updated:"), valueStyle.Render(f.formatTime(discussion.UpdatedAt)))

	fmt.Fprintf(f.writer, "%s ", labelStyle.Render("Answered:"))
	if discussion.IsAnswered {
		fmt.Fprintf(f.writer, "%s\n", answerStyle.Render("Yes"))
	} else {
		fmt.Fprintf(f.writer, "%s\n", valueStyle.Render("No"))
	}

	if discussion.Comments != nil {
		fmt.Fprintf(f.writer, "%s %s\n", labelStyle.Render("Comments:"), valueStyle.Render(strconv.Itoa(discussion.Comments.TotalCount)))
	}

	fmt.Fprintf(f.writer, "%s %s\n", labelStyle.Render("URL:"), urlStyle.Render(discussion.URL))

	if discussion.Labels != nil && len(discussion.Labels.Nodes) > 0 {
		labels := make([]string, len(discussion.Labels.Nodes))
		for i, label := range discussion.Labels.Nodes {
			labels[i] = label.Name
		}
		fmt.Fprintf(f.writer, "%s %s\n", labelStyle.Render("Labels:"), valueStyle.Render(strings.Join(labels, ", ")))
	}

	// Body section
	if discussion.Body != "" {
		fmt.Fprintf(f.writer, "\n%s\n", separator)
		fmt.Fprintf(f.writer, "\n%s\n", discussion.Body)
	}

	// Comments section
	if discussion.Comments != nil && len(discussion.Comments.Nodes) > 0 {
		commentsHeaderStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true)

		fmt.Fprintf(f.writer, "\n%s\n", separator)
		fmt.Fprintf(f.writer, "\n%s\n", commentsHeaderStyle.Render("Comments"))

		for i, comment := range discussion.Comments.Nodes {
			f.formatComment(comment, i+1, 0)
		}
	}

	return nil
}

// formatComment formats a single comment with its replies
func (f *Formatter) formatComment(comment models.Comment, number int, depth int) {
	indent := strings.Repeat("  ", depth)

	// Comment header styling
	commentHeaderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true)

	// Author styling
	authorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)

	// Time styling
	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Answer badge styling
	answerBadgeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("10")).
		Padding(0, 1).
		Bold(true)

	// Comment separator for depth > 0
	if depth > 0 {
		commentSeparator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("┌─")
		fmt.Fprintf(f.writer, "\n%s%s\n", indent, commentSeparator)
	} else {
		fmt.Fprintf(f.writer, "\n")
	}

	fmt.Fprintf(f.writer, "%s%s", indent, commentHeaderStyle.Render(fmt.Sprintf("Comment #%d", number)))
	if comment.IsAnswer {
		fmt.Fprintf(f.writer, " %s", answerBadgeStyle.Render("Answer"))
	}
	fmt.Fprintf(f.writer, "\n")

	if comment.Author != nil {
		fmt.Fprintf(f.writer, "%s%s • %s\n", indent, authorStyle.Render(comment.Author.Login), timeStyle.Render(f.formatTime(comment.CreatedAt)))
	} else {
		fmt.Fprintf(f.writer, "%s%s\n", indent, timeStyle.Render(f.formatTime(comment.CreatedAt)))
	}

	// Format the comment body with proper indentation
	bodyLines := strings.Split(comment.Body, "\n")
	fmt.Fprintf(f.writer, "\n")
	for _, line := range bodyLines {
		fmt.Fprintf(f.writer, "%s%s\n", indent, line)
	}

	// Show replies if available
	if comment.Replies != nil && len(comment.Replies.Nodes) > 0 {
		for i, reply := range comment.Replies.Nodes {
			f.formatComment(reply, i+1, depth+1)
		}
	}
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

// truncateString truncates a string to the specified length, handling Unicode properly
func (f *Formatter) truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return "..."
	}
	return string(runes[:maxLen-3]) + "..."
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
