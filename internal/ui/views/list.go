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

// CursorMovedMsg is sent when the cursor moves to a different conversation.
type CursorMovedMsg struct {
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
		return v.moveSelectionWithCmd(-1)

	case tea.KeyDown:
		return v.moveSelectionWithCmd(1)

	case tea.KeyHome:
		if len(v.conversations) == 0 {
			return v, nil
		}
		newView := v.SetSelectedIndex(0)
		return newView, v.sendCursorMovedCmd(newView)

	case tea.KeyEnd:
		if len(v.conversations) == 0 {
			return v, nil
		}
		newView := v.SetSelectedIndex(len(v.conversations) - 1)
		return newView, v.sendCursorMovedCmd(newView)

	case tea.KeyPgUp:
		pageSize := v.calculatePageSize()
		return v.moveSelectionWithCmd(-pageSize)

	case tea.KeyPgDown:
		pageSize := v.calculatePageSize()
		return v.moveSelectionWithCmd(pageSize)

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
				return v.moveSelectionWithCmd(1)
			case 'k':
				return v.moveSelectionWithCmd(-1)
			}
		}
		return v, nil

	default:
		return v, nil
	}
}

// sendCursorMovedCmd returns a command to send CursorMovedMsg for the current selection.
func (v *ListView) sendCursorMovedCmd(newView *ListView) tea.Cmd {
	conv := newView.SelectedConversation()
	if conv != nil {
		return func() tea.Msg {
			return CursorMovedMsg{Conversation: conv}
		}
	}
	return nil
}

// moveSelection moves the selection by delta and clamps to valid bounds.
// Returns the new view and a command to notify about cursor movement.
func (v *ListView) moveSelectionWithCmd(delta int) (*ListView, tea.Cmd) {
	newView := v.SetSelectedIndex(v.selectedIndex + delta)

	// Send cursor moved message if selection actually changed
	if newView.selectedIndex != v.selectedIndex {
		conv := newView.SelectedConversation()
		if conv != nil {
			return newView, func() tea.Msg {
				return CursorMovedMsg{Conversation: conv}
			}
		}
	}
	return newView, nil
}

// moveSelection moves the selection by delta (legacy, doesn't send message).
func (v *ListView) moveSelection(delta int) *ListView {
	return v.SetSelectedIndex(v.selectedIndex + delta)
}

// calculatePageSize returns the number of items that fit on one page.
func (v *ListView) calculatePageSize() int {
	// Rough estimate: height minus header/footer
	pageSize := v.height - 5
	if pageSize < 1 {
		// L1: Use reasonable fallback based on height, not hardcoded 10
		pageSize = max(3, v.height/2)
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

// renderTable renders the conversation list with pagination.
func (v *ListView) renderTable() string {
	sorted := v.getSortedConversations()
	if len(sorted) == 0 {
		return "No conversations"
	}

	var b strings.Builder
	pageSize := v.calculatePageSize()
	startIdx := (v.selectedIndex / pageSize) * pageSize
	endIdx := startIdx + pageSize
	if endIdx > len(sorted) {
		endIdx = len(sorted)
	}

	// Calculate display width for title (leave room for count and date)
	// Format: "> 123 2024-01-15 Title..." = 2 + 3 + 1 + 11 + 1 = 18 chars overhead
	titleWidth := v.width - 18
	if titleWidth < 15 {
		titleWidth = 15
	}
	if titleWidth > 70 {
		titleWidth = 70
	}

	for i := startIdx; i < endIdx; i++ {
		conv := sorted[i]
		isSelected := i == v.selectedIndex

		// Format fields (compact)
		msgs := fmt.Sprintf("%3d", conv.MessageCount())
		date := v.formatTimestamp(conv.UpdatedAt)
		title := truncateString(conv.Title, titleWidth)

		// Build line with selection indicator (single spaces)
		prefix := "  "
		if isSelected {
			prefix = "> "
		}
		line := fmt.Sprintf("%s%s %s %s", prefix, msgs, date, title)

		b.WriteString(line)
		b.WriteString("\n")
	}

	// Pagination info
	totalPages := (len(sorted) + pageSize - 1) / pageSize
	currentPage := (v.selectedIndex / pageSize) + 1
	b.WriteString(fmt.Sprintf("\nPage %d/%d (%d total)", currentPage, totalPages, len(sorted)))

	return b.String()
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
	title := truncateString(conv.Title, 33)
	messages := fmt.Sprintf("%d", conv.MessageCount())
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

// formatTimestamp formats a timestamp for display with fixed width (11 chars).
func (v *ListView) formatTimestamp(ts time.Time) string {
	// Handle zero timestamps
	if ts.IsZero() {
		return fmt.Sprintf("%-11s", "Unknown")
	}

	now := time.Now()
	diff := now.Sub(ts)

	var result string
	switch {
	case diff < time.Minute:
		result = "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		result = fmt.Sprintf("%dm ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		result = fmt.Sprintf("%dh ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		result = fmt.Sprintf("%dd ago", days)
	default:
		result = ts.Format("2006-01-02")
	}

	// Pad to fixed width for column alignment
	return fmt.Sprintf("%-11s", result)
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

// truncateString safely truncates a string to maxLen runes, adding "..." if truncated.
// This handles multi-byte UTF-8 characters correctly.
func truncateString(s string, maxLen int) string {
	if maxLen <= 3 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
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
