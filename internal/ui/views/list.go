// ABOUTME: ListView component displays a table of conversations with sortable columns
// for title, model, message count, tokens, and timestamps. Supports keyboard navigation,
// sorting, and responsive layout adaptation.
package views

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/TrevorS/claudex/internal/repository"
)

// ListView displays a table of conversations with metadata.
type ListView struct {
	repository    repository.Repository
	table         table.Model
	conversations []*domain.Conversation
	selectedIndex int
	sortColumn    string
	sortAsc       bool
	width         int
	height        int
	err           *domain.AppError
}

// ConversationsLoadedMsg is sent when conversations are loaded from the repository.
type ConversationsLoadedMsg struct {
	Conversations []*domain.Conversation
	Err           error
}

// SelectConversationMsg is sent when a conversation is selected.
type SelectConversationMsg struct {
	Conversation *domain.Conversation
}

// NewListView creates a new ListView with the given repository.
func NewListView(repo repository.Repository) *ListView {
	return &ListView{
		repository:    repo,
		conversations: make([]*domain.Conversation, 0),
		selectedIndex: 0,
		sortColumn:    "timestamp",
		sortAsc:       false, // Default: newest first
		width:         80,
		height:        24,
	}
}

// Init initializes the ListView component and loads conversations.
func (v *ListView) Init() tea.Cmd {
	return v.loadConversationsCmd()
}

// loadConversationsCmd returns a command that loads conversations from the repository.
func (v *ListView) loadConversationsCmd() tea.Cmd {
	return func() tea.Msg {
		conversations, err := v.repository.List(context.Background())
		return ConversationsLoadedMsg{
			Conversations: conversations,
			Err:           err,
		}
	}
}

// Update handles messages and updates the ListView state.
func (v *ListView) Update(msg tea.Msg) (*ListView, tea.Cmd) {
	switch msg := msg.(type) {

	case ConversationsLoadedMsg:
		if msg.Err != nil {
			newView := *v
			if appErr, ok := msg.Err.(*domain.AppError); ok {
				newView.err = appErr
			} else {
				newView.err = domain.NewAppError(domain.ErrRepository, "failed to load conversations", msg.Err)
			}
			return &newView, nil
		}
		return v.SetConversations(msg.Conversations), nil

	case tea.WindowSizeMsg:
		return v.SetSize(msg.Width, msg.Height), nil

	case tea.KeyMsg:
		return v.handleKeyPress(msg)

	default:
		return v, nil
	}
}

// handleKeyPress handles keyboard input.
func (v *ListView) handleKeyPress(msg tea.KeyMsg) (*ListView, tea.Cmd) {
	switch msg.Type {

	case tea.KeyUp:
		return v.moveSelection(-1), nil

	case tea.KeyDown:
		return v.moveSelection(1), nil

	case tea.KeyHome:
		return v.SetSelectedIndex(0), nil

	case tea.KeyEnd:
		return v.SetSelectedIndex(len(v.conversations) - 1), nil

	case tea.KeyPgUp:
		pageSize := v.calculatePageSize()
		return v.moveSelection(-pageSize), nil

	case tea.KeyPgDown:
		pageSize := v.calculatePageSize()
		return v.moveSelection(pageSize), nil

	case tea.KeyEnter:
		selected := v.SelectedConversation()
		if selected != nil {
			return v, func() tea.Msg {
				return SelectConversationMsg{Conversation: selected}
			}
		}
		return v, nil

	case tea.KeyRunes:
		if len(msg.Runes) > 0 {
			switch msg.Runes[0] {
			case 'j':
				return v.moveSelection(1), nil
			case 'k':
				return v.moveSelection(-1), nil
			}
		}
		return v, nil

	default:
		return v, nil
	}
}

// moveSelection moves the selection by delta and clamps to valid bounds.
func (v *ListView) moveSelection(delta int) *ListView {
	newIndex := v.selectedIndex + delta
	return v.SetSelectedIndex(newIndex)
}

// calculatePageSize returns the number of items that fit on one page.
func (v *ListView) calculatePageSize() int {
	// Rough estimate: height minus header/footer
	pageSize := v.height - 5
	if pageSize < 1 {
		pageSize = 10
	}
	return pageSize
}

// View renders the ListView as a string.
func (v *ListView) View() string {
	if v.err != nil {
		return fmt.Sprintf("Error loading conversations: %s", v.err.Error())
	}

	if len(v.conversations) == 0 {
		return v.renderEmptyState()
	}

	return v.renderTable()
}

// renderEmptyState displays a message when there are no conversations.
func (v *ListView) renderEmptyState() string {
	return "No conversations found.\n\nPress Ctrl+C to exit."
}

// renderTable renders the conversation table.
func (v *ListView) renderTable() string {
	sorted := v.getSortedConversations()

	// Build table rows
	rows := make([]table.Row, 0, len(sorted))
	for i, conv := range sorted {
		row := v.buildRow(conv, i)
		rows = append(rows, row)
	}

	// Create columns based on width
	columns := v.buildColumns()

	// Create table model
	tbl := table.New(columns).
		WithRows(rows).
		Focused(true).
		WithPageSize(v.calculatePageSize()).
		WithTargetWidth(v.width)

	// Highlight selected row
	if v.selectedIndex >= 0 && v.selectedIndex < len(rows) {
		tbl = tbl.WithHighlightedRow(v.selectedIndex)
	}

	return tbl.View()
}

// buildColumns creates table columns based on terminal width.
func (v *ListView) buildColumns() []table.Column {
	if v.width < 40 {
		// Minimal: just title
		return []table.Column{
			table.NewColumn("title", "Title", 35),
		}
	}

	if v.width < 80 {
		// Narrow: no cost column
		return []table.Column{
			table.NewColumn("timestamp", "Time", 16),
			table.NewColumn("title", "Title", 30),
			table.NewColumn("messages", "Msgs", 6),
			table.NewColumn("tokens", "Tokens", 8),
		}
	}

	// Wide: all columns
	return []table.Column{
		table.NewColumn("timestamp", "Time", 18),
		table.NewColumn("title", "Title", 35),
		table.NewColumn("messages", "Msgs", 6),
		table.NewColumn("tokens", "Tokens", 10),
		table.NewColumn("model", "Model", 18),
	}
}

// buildRow creates a table row from a conversation.
func (v *ListView) buildRow(conv *domain.Conversation, index int) table.Row {
	timestamp := v.formatTimestamp(conv.UpdatedAt)
	title := conv.Title
	if len(title) > 33 {
		title = title[:30] + "..."
	}
	messages := fmt.Sprintf("%d", len(conv.Messages))
	tokens := v.formatTokens(conv.TotalTokens())
	model := v.formatModel(conv.Model)

	rowData := table.RowData{
		"timestamp": timestamp,
		"title":     title,
		"messages":  messages,
		"tokens":    tokens,
		"model":     model,
	}

	return table.NewRow(rowData)
}

// formatTimestamp formats a timestamp for display.
func (v *ListView) formatTimestamp(ts time.Time) string {
	now := time.Now()
	diff := now.Sub(ts)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	default:
		return ts.Format("2006-01-02")
	}
}

// formatTokens formats token count for display.
func (v *ListView) formatTokens(tokens int64) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	}
	if tokens < 1000000 {
		return fmt.Sprintf("%.1fK", float64(tokens)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(tokens)/1000000)
}

// formatModel formats model name for display.
func (v *ListView) formatModel(model string) string {
	// Simplify model names (e.g., "claude-sonnet-4" -> "sonnet-4")
	parts := strings.Split(model, "-")
	if len(parts) > 1 {
		return strings.Join(parts[1:], "-")
	}
	return model
}

// getSortedConversations returns conversations sorted by current sort settings.
func (v *ListView) getSortedConversations() []*domain.Conversation {
	sorted := make([]*domain.Conversation, len(v.conversations))
	copy(sorted, v.conversations)

	sort.Slice(sorted, func(i, j int) bool {
		a, b := sorted[i], sorted[j]
		var less bool

		switch v.sortColumn {
		case "title":
			less = a.Title < b.Title
		case "messages":
			less = len(a.Messages) < len(b.Messages)
		case "tokens":
			less = a.TotalTokens() < b.TotalTokens()
		case "timestamp":
			less = a.UpdatedAt.Before(b.UpdatedAt)
		default:
			less = a.UpdatedAt.Before(b.UpdatedAt)
		}

		if v.sortAsc {
			return less
		}
		return !less
	})

	return sorted
}

// Conversations returns the current list of conversations.
func (v *ListView) Conversations() []*domain.Conversation {
	return v.conversations
}

// SetConversations sets the conversations list and returns a new ListView.
func (v *ListView) SetConversations(conversations []*domain.Conversation) *ListView {
	newView := *v
	newView.conversations = conversations
	// Clamp selection to valid range
	if newView.selectedIndex >= len(conversations) {
		newView.selectedIndex = len(conversations) - 1
	}
	if newView.selectedIndex < 0 && len(conversations) > 0 {
		newView.selectedIndex = 0
	}
	return &newView
}

// SelectedIndex returns the currently selected conversation index.
func (v *ListView) SelectedIndex() int {
	return v.selectedIndex
}

// SetSelectedIndex sets the selected conversation index with bounds checking.
func (v *ListView) SetSelectedIndex(index int) *ListView {
	newView := *v

	// Clamp to valid range
	if len(v.conversations) == 0 {
		newView.selectedIndex = 0
		return &newView
	}

	if index < 0 {
		newView.selectedIndex = 0
	} else if index >= len(v.conversations) {
		newView.selectedIndex = len(v.conversations) - 1
	} else {
		newView.selectedIndex = index
	}

	return &newView
}

// SelectedConversation returns the currently selected conversation, or nil if none selected.
func (v *ListView) SelectedConversation() *domain.Conversation {
	sorted := v.getSortedConversations()
	if v.selectedIndex < 0 || v.selectedIndex >= len(sorted) {
		return nil
	}
	return sorted[v.selectedIndex]
}

// SetSort sets the sort column and direction, returning a new ListView.
func (v *ListView) SetSort(column string, ascending bool) *ListView {
	newView := *v
	newView.sortColumn = column
	newView.sortAsc = ascending
	return &newView
}

// SetSize sets the terminal dimensions, returning a new ListView.
func (v *ListView) SetSize(width, height int) *ListView {
	newView := *v
	newView.width = width
	newView.height = height
	return &newView
}
