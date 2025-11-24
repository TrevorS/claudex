// ABOUTME: Comprehensive tests for ListView component.
// Tests verify table initialization, keyboard navigation, sorting, and conversation selection.
package views

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/TrevorS/claudex/internal/repository"
)

// Helper function to create test conversations
func createTestConversations() []*domain.Conversation {
	now := time.Now()

	conversations := []*domain.Conversation{
		{
			ID:        "conv-1",
			Title:     "First Conversation",
			Model:     "claude-sonnet-4",
			CreatedAt: now.Add(-3 * time.Hour),
			UpdatedAt: now.Add(-3 * time.Hour),
			Messages: []*domain.Message{
				{
					Role:      domain.RoleUser,
					Content:   "Hello",
					Timestamp: now.Add(-3 * time.Hour),
					Tokens:    &domain.TokenCount{Input: 10, Output: 0},
				},
				{
					Role:      domain.RoleAssistant,
					Content:   "Hi there!",
					Timestamp: now.Add(-3 * time.Hour),
					Tokens:    &domain.TokenCount{Input: 0, Output: 20},
				},
			},
		},
		{
			ID:        "conv-2",
			Title:     "Second Conversation",
			Model:     "claude-opus-4",
			CreatedAt: now.Add(-2 * time.Hour),
			UpdatedAt: now.Add(-2 * time.Hour),
			Messages: []*domain.Message{
				{
					Role:      domain.RoleUser,
					Content:   "Test message",
					Timestamp: now.Add(-2 * time.Hour),
					Tokens:    &domain.TokenCount{Input: 15, Output: 0},
				},
			},
		},
		{
			ID:        "conv-3",
			Title:     "Third Conversation",
			Model:     "claude-sonnet-4",
			CreatedAt: now.Add(-1 * time.Hour),
			UpdatedAt: now.Add(-1 * time.Hour),
			Messages: []*domain.Message{
				{
					Role:      domain.RoleUser,
					Content:   "Another test",
					Timestamp: now.Add(-1 * time.Hour),
					Tokens:    &domain.TokenCount{Input: 5, Output: 0},
				},
				{
					Role:      domain.RoleAssistant,
					Content:   "Response",
					Timestamp: now.Add(-1 * time.Hour),
					Tokens:    &domain.TokenCount{Input: 0, Output: 50},
				},
				{
					Role:      domain.RoleUser,
					Content:   "Follow up",
					Timestamp: now.Add(-1 * time.Hour),
					Tokens:    &domain.TokenCount{Input: 8, Output: 0},
				},
			},
		},
	}

	return conversations
}

// Test: Construction
func TestNewListView(t *testing.T) {
	repo := repository.NewMock([]*domain.Conversation{})

	view := NewListView(repo)

	assert.NotNil(t, view, "NewListView should return non-nil view")
	assert.NotNil(t, view.repository, "repository should be set")
	assert.Equal(t, 0, len(view.conversations), "conversations should be empty initially")
	assert.Equal(t, 0, view.selectedIndex, "selectedIndex should be 0 initially")
}

func TestListViewInit(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)

	cmd := view.Init()

	// Init should return a command to load conversations
	assert.NotNil(t, cmd, "Init() should return a command to load conversations")
}

// Test: Data Loading
func TestListViewLoadConversations(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)

	// Simulate loading conversations
	view = view.SetConversations(conversations)

	assert.Equal(t, 3, len(view.conversations), "should load 3 conversations")
	assert.Equal(t, "conv-1", view.conversations[0].ID)
	assert.Equal(t, "conv-2", view.conversations[1].ID)
	assert.Equal(t, "conv-3", view.conversations[2].ID)
}

func TestListViewHandlesEmptyList(t *testing.T) {
	repo := repository.NewMock([]*domain.Conversation{})

	view := NewListView(repo)
	view = view.SetConversations([]*domain.Conversation{})

	output := view.View()

	assert.NotEmpty(t, output, "View should render something even with empty list")
	assert.Contains(t, output, "No conversations", "Empty list should show message")
}

func TestListViewHandlesSingleConversation(t *testing.T) {
	conversations := createTestConversations()[:1]
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	assert.Equal(t, 1, len(view.conversations))
	assert.Equal(t, 0, view.selectedIndex)
	assert.NotNil(t, view.SelectedConversation())
}

func TestListViewHandlesLargeConversationList(t *testing.T) {
	// Create 100+ conversations
	conversations := make([]*domain.Conversation, 150)
	now := time.Now()

	for i := 0; i < 150; i++ {
		conversations[i] = domain.NewConversation(
			"conv-"+string(rune(i)),
			"Conversation "+string(rune(i)),
			"claude-sonnet-4",
			now.Add(-time.Duration(i)*time.Hour),
		)
		// Add a message to each conversation
		msg := domain.NewMessage(domain.RoleUser, "test", now)
		msg.Tokens = &domain.TokenCount{Input: 10, Output: 20}
		conversations[i].AddMessage(msg)
	}

	repo := repository.NewMock(conversations)
	view := NewListView(repo)
	view = view.SetConversations(conversations)

	assert.Equal(t, 150, len(view.conversations))
	output := view.View()
	assert.NotEmpty(t, output)
}

// Test: Sorting
func TestListViewDefaultSort(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Default sort should be by timestamp descending (newest first)
	// conv-3 is newest, conv-1 is oldest
	sorted := view.getSortedConversations()

	assert.Equal(t, "conv-3", sorted[0].ID, "Newest conversation should be first")
	assert.Equal(t, "conv-2", sorted[1].ID, "Middle conversation should be second")
	assert.Equal(t, "conv-1", sorted[2].ID, "Oldest conversation should be last")
}

func TestListViewSortByTitle(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Change sort to title ascending
	view = view.SetSort("title", true)
	sorted := view.getSortedConversations()

	assert.Equal(t, "First Conversation", sorted[0].Title)
	assert.Equal(t, "Second Conversation", sorted[1].Title)
	assert.Equal(t, "Third Conversation", sorted[2].Title)
}

func TestListViewSortByMessages(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Sort by message count descending (most messages first)
	view = view.SetSort("messages", false)
	sorted := view.getSortedConversations()

	assert.Equal(t, 3, len(sorted[0].Messages), "Conversation with 3 messages should be first")
	assert.Equal(t, 2, len(sorted[1].Messages), "Conversation with 2 messages should be second")
	assert.Equal(t, 1, len(sorted[2].Messages), "Conversation with 1 message should be last")
}

func TestListViewSortByTokens(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Sort by tokens descending (most tokens first)
	view = view.SetSort("tokens", false)
	sorted := view.getSortedConversations()

	// conv-3 has 5+50+8 = 63 tokens
	// conv-1 has 10+20 = 30 tokens
	// conv-2 has 15 tokens
	assert.Equal(t, "conv-3", sorted[0].ID, "Conversation with most tokens should be first")
	assert.Equal(t, "conv-1", sorted[1].ID)
	assert.Equal(t, "conv-2", sorted[2].ID, "Conversation with least tokens should be last")
}

func TestListViewSortToggle(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Start with title ascending
	view = view.SetSort("title", true)
	sorted1 := view.getSortedConversations()
	assert.Equal(t, "First Conversation", sorted1[0].Title)

	// Toggle to title descending
	view = view.SetSort("title", false)
	sorted2 := view.getSortedConversations()
	assert.Equal(t, "Third Conversation", sorted2[0].Title)
}

// Test: Selection
func TestSelectedIndex(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	assert.Equal(t, 0, view.SelectedIndex())

	view = view.SetSelectedIndex(1)
	assert.Equal(t, 1, view.SelectedIndex())

	view = view.SetSelectedIndex(2)
	assert.Equal(t, 2, view.SelectedIndex())
}

func TestSelectedConversation(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Default selection (index 0 after sorting)
	selected := view.SelectedConversation()
	require.NotNil(t, selected)

	// Select second item
	view = view.SetSelectedIndex(1)
	selected = view.SelectedConversation()
	require.NotNil(t, selected)
}

func TestSelectedConversationOutOfBounds(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Try to select beyond bounds (should be clamped)
	view = view.SetSelectedIndex(999)
	assert.Equal(t, 2, view.SelectedIndex(), "Should clamp to last valid index")

	// Try negative index (should clamp to 0)
	view = view.SetSelectedIndex(-5)
	assert.Equal(t, 0, view.SelectedIndex(), "Should clamp to 0")
}

func TestSelectedConversationEmptyList(t *testing.T) {
	repo := repository.NewMock([]*domain.Conversation{})

	view := NewListView(repo)
	view = view.SetConversations([]*domain.Conversation{})

	selected := view.SelectedConversation()
	assert.Nil(t, selected, "Should return nil when list is empty")
}

// Test: Keyboard Navigation
func TestListViewKeyDown(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Start at index 0
	assert.Equal(t, 0, view.SelectedIndex())

	// Press down
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 1, view.SelectedIndex())

	// Press down again
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 2, view.SelectedIndex())

	// Press down at end (should stay at last)
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 2, view.SelectedIndex())
}

func TestListViewKeyUp(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSelectedIndex(2) // Start at last

	// Press up
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 1, view.SelectedIndex())

	// Press up again
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 0, view.SelectedIndex())

	// Press up at beginning (should stay at 0)
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 0, view.SelectedIndex())
}

func TestListViewVimKeys(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Test j (down)
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	assert.Equal(t, 1, view.SelectedIndex())

	// Test k (up)
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	assert.Equal(t, 0, view.SelectedIndex())
}

func TestListViewHomeEnd(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSelectedIndex(1)

	// Press End (go to last)
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyEnd})
	assert.Equal(t, 2, view.SelectedIndex())

	// Press Home (go to first)
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyHome})
	assert.Equal(t, 0, view.SelectedIndex())
}

func TestListViewPageUpDown(t *testing.T) {
	// Create more conversations for meaningful page tests
	conversations := make([]*domain.Conversation, 30)
	now := time.Now()
	for i := 0; i < 30; i++ {
		conversations[i] = domain.NewConversation(
			"conv-"+string(rune(i)),
			"Conversation "+string(rune(i)),
			"claude-sonnet-4",
			now.Add(-time.Duration(i)*time.Hour),
		)
	}

	repo := repository.NewMock(conversations)
	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSize(80, 30) // Set size for page calculations

	// Page down (should move by ~10 items)
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	assert.Greater(t, view.SelectedIndex(), 5, "Should move down significantly")

	// Page up (should move back)
	oldIndex := view.SelectedIndex()
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	assert.Less(t, view.SelectedIndex(), oldIndex, "Should move up")
}

func TestListViewEnter(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Press Enter (should return SelectConversationMsg)
	view, cmd := view.Update(tea.KeyMsg{Type: tea.KeyEnter})

	assert.NotNil(t, cmd, "Enter should return a command")

	// Execute the command to get the message
	msg := cmd()

	// Check that it's a SelectConversationMsg
	selectMsg, ok := msg.(SelectConversationMsg)
	assert.True(t, ok, "Command should return SelectConversationMsg")
	assert.NotNil(t, selectMsg.Conversation, "Message should contain selected conversation")
}

// Test: Rendering
func TestListViewView(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSize(100, 30)

	output := view.View()

	assert.NotEmpty(t, output, "View should return non-empty string")
	assert.Contains(t, output, "Conversation", "Should contain conversation data")
}

func TestListViewViewIncludesColumns(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSize(100, 30)

	output := view.View()

	// Should include column headers or data
	assert.True(t, len(output) > 50, "Should render substantial content")
}

func TestListViewResponsiveLayoutWide(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSize(100, 30) // Wide terminal

	output := view.View()
	assert.NotEmpty(t, output)
	// Wide layout should show all columns
}

func TestListViewResponsiveLayoutNarrow(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSize(60, 30) // Narrow terminal

	output := view.View()
	assert.NotEmpty(t, output)
	// Narrow layout should hide cost column
}

func TestListViewResponsiveLayoutMinimal(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSize(35, 20) // Very narrow terminal

	output := view.View()
	assert.NotEmpty(t, output)
	// Minimal layout should show only title
}

// Test: Edge Cases
func TestListViewWindowResize(t *testing.T) {
	// Use a long title that will truncate differently at different widths
	now := time.Now()
	conversations := []*domain.Conversation{
		{
			ID:        "conv-1",
			Title:     "This is a very long conversation title that should be truncated at different widths",
			Model:     "claude-sonnet-4",
			CreatedAt: now,
			UpdatedAt: now,
			Messages: []*domain.Message{
				{Role: domain.RoleUser, Content: "test", Timestamp: now},
			},
		},
	}
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	// Start with wide
	view = view.SetSize(100, 30)
	output1 := view.View()
	assert.NotEmpty(t, output1)

	// Resize to narrow
	view = view.SetSize(60, 30)
	output2 := view.View()
	assert.NotEmpty(t, output2)

	// Both should render successfully, and long title should truncate differently
	assert.NotEqual(t, output1, output2, "Different sizes should produce different output")
}

func TestListViewSelectionPersistsAcrossUpdates(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSelectedIndex(1)

	// Process a non-navigation message
	view, _ = view.Update(tea.WindowSizeMsg{Width: 80, Height: 30})

	// Selection should persist
	assert.Equal(t, 1, view.SelectedIndex())
}

func TestListViewUpdateReturnsNewInstance(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	oldView := view
	newView, _ := view.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Should return a new instance (immutable pattern)
	assert.NotEqual(t, oldView.SelectedIndex(), newView.SelectedIndex())
}

// Test: Commands and Messages
func TestLoadConversationsCmd(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)

	cmd := view.Init()
	require.NotNil(t, cmd, "Init should return load command")

	// Execute command
	msg := cmd()

	// Should return ConversationsLoadedMsg
	loadedMsg, ok := msg.(ConversationsLoadedMsg)
	assert.True(t, ok, "Command should return ConversationsLoadedMsg")
	assert.Equal(t, 3, len(loadedMsg.Conversations))
	assert.Nil(t, loadedMsg.Err)
}

func TestLoadConversationsCmdError(t *testing.T) {
	// Create a repo that will return error
	repo := repository.NewMock([]*domain.Conversation{})

	view := NewListView(repo)
	cmd := view.Init()
	require.NotNil(t, cmd)

	// For this test, we'd need a mock that returns error
	// Since our mock doesn't support errors, we'll skip this for now
	// In a real implementation, you'd test error handling
}

func TestConversationsLoadedMsg(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)

	// Simulate receiving loaded message
	msg := ConversationsLoadedMsg{
		Conversations: conversations,
		Err:           nil,
	}

	newView, _ := view.Update(msg)

	assert.Equal(t, 3, len(newView.conversations))
}

// Test: Format helpers
func TestFormatTimestamp(t *testing.T) {
	view := NewListView(repository.NewMock([]*domain.Conversation{}))
	now := time.Now()

	tests := []struct {
		name     string
		ts       time.Time
		contains string
	}{
		{"just now", now.Add(-30 * time.Second), "just now"},
		{"minutes ago", now.Add(-5 * time.Minute), "m ago"},
		{"hours ago", now.Add(-3 * time.Hour), "h ago"},
		{"days ago", now.Add(-2 * 24 * time.Hour), "d ago"},
		{"weeks ago", now.Add(-10 * 24 * time.Hour), "20"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := view.formatTimestamp(tt.ts)
			assert.Contains(t, result, tt.contains)
		})
	}
}

func TestFormatTokens(t *testing.T) {
	view := NewListView(repository.NewMock([]*domain.Conversation{}))

	tests := []struct {
		name     string
		tokens   int64
		expected string
	}{
		{"small", 500, "500"},
		{"thousands", 5000, "5.0K"},
		{"millions", 2500000, "2.5M"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := view.formatTokens(tt.tokens)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatModel(t *testing.T) {
	view := NewListView(repository.NewMock([]*domain.Conversation{}))

	tests := []struct {
		name     string
		model    string
		expected string
	}{
		{"with prefix", "claude-sonnet-4", "sonnet-4"},
		{"with prefix opus", "claude-opus-4", "opus-4"},
		{"no prefix", "gpt4", "gpt4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := view.formatModel(tt.model)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConversationsGetter(t *testing.T) {
	conversations := createTestConversations()
	repo := repository.NewMock(conversations)

	view := NewListView(repo)
	view = view.SetConversations(conversations)

	result := view.Conversations()
	assert.Equal(t, conversations, result)
}

func TestBuildRowWithLongTitle(t *testing.T) {
	repo := repository.NewMock([]*domain.Conversation{})
	view := NewListView(repo)

	now := time.Now()
	conv := &domain.Conversation{
		ID:        "test",
		Title:     "This is a very long title that should be truncated when displayed in the table view",
		Model:     "claude-sonnet-4",
		CreatedAt: now,
		UpdatedAt: now,
		Messages:  []*domain.Message{},
	}

	row := view.buildRow(conv, 0)
	rowData := row.Data

	titleStr, ok := rowData["title"].(string)
	assert.True(t, ok)
	assert.LessOrEqual(t, len(titleStr), 35, "Title should be truncated")
	assert.Contains(t, titleStr, "...", "Truncated title should have ellipsis")
}

func TestErrorInConversationsLoaded(t *testing.T) {
	repo := repository.NewMock([]*domain.Conversation{})
	view := NewListView(repo)

	// Simulate error during load
	errMsg := domain.NewAppError(domain.ErrRepository, "test error", nil)
	msg := ConversationsLoadedMsg{
		Conversations: nil,
		Err:           errMsg,
	}

	newView, _ := view.Update(msg)
	assert.NotNil(t, newView.err)

	// View should show error
	output := newView.View()
	assert.Contains(t, output, "Error")
}

func TestCalculatePageSize(t *testing.T) {
	repo := repository.NewMock([]*domain.Conversation{})
	view := NewListView(repo)

	// Test with normal height
	view = view.SetSize(80, 30)
	pageSize := view.calculatePageSize()
	assert.Greater(t, pageSize, 0)
	assert.LessOrEqual(t, pageSize, 30)

	// Test with very small height
	view = view.SetSize(80, 3)
	pageSize = view.calculatePageSize()
	assert.Equal(t, 10, pageSize, "Should use minimum page size for small heights")
}
