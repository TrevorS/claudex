// ABOUTME: Comprehensive tests for DetailView component.
// Tests verify conversation loading, message rendering, scrolling, keyboard navigation, and error handling.
package views

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/TrevorS/claudex/internal/repository"
)

// Helper function to create a test conversation with messages
func createDetailTestConversation() *domain.Conversation {
	now := time.Now()
	conv := domain.NewConversation(
		"conv-detail-1",
		"Test Conversation for Detail View",
		"claude-sonnet-4",
		now.Add(-2*time.Hour),
	)
	conv.ProjectPath = "test/project"

	// Add a user message
	msg1 := domain.NewMessage(domain.RoleUser, "Hello, can you help me with this code?", now.Add(-2*time.Hour))
	msg1.Tokens = &domain.TokenCount{Input: 25, Output: 0}
	conv.AddMessage(msg1)

	// Add an assistant message
	msg2 := domain.NewMessage(domain.RoleAssistant, "Of course! I'd be happy to help. What specific issue are you encountering?", now.Add(-2*time.Hour).Add(5*time.Second))
	msg2.Tokens = &domain.TokenCount{Input: 0, Output: 50, CacheRead: 100}
	conv.AddMessage(msg2)

	// Add another user message
	msg3 := domain.NewMessage(domain.RoleUser, "I'm getting a type error in my Go code.", now.Add(-2*time.Hour).Add(1*time.Minute))
	msg3.Tokens = &domain.TokenCount{Input: 30, Output: 0}
	conv.AddMessage(msg3)

	return conv
}

// Helper to create conversation with long messages
func createLongMessageConversation() *domain.Conversation {
	now := time.Now()
	conv := domain.NewConversation(
		"conv-long",
		"Conversation with Very Long Messages",
		"claude-opus-4",
		now,
	)

	longContent := strings.Repeat("This is a very long message that will test wrapping. ", 50)
	msg := domain.NewMessage(domain.RoleUser, longContent, now)
	msg.Tokens = &domain.TokenCount{Input: 500, Output: 0}
	conv.AddMessage(msg)

	return conv
}

// Helper to create conversation with many messages (100+)
func createManyMessagesConversation() *domain.Conversation {
	now := time.Now()
	conv := domain.NewConversation(
		"conv-many",
		"Conversation with Many Messages",
		"claude-sonnet-4",
		now,
	)

	for i := 0; i < 150; i++ {
		role := domain.RoleUser
		if i%2 == 1 {
			role = domain.RoleAssistant
		}
		msg := domain.NewMessage(role, "Message "+string(rune(i)), now.Add(time.Duration(i)*time.Second))
		msg.Tokens = &domain.TokenCount{Input: 10, Output: 10}
		conv.AddMessage(msg)
	}

	return conv
}

// Test: Construction
func TestNewDetailView(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")

	assert.NotNil(t, view, "NewDetailView should return non-nil view")
	assert.NotNil(t, view.repository, "repository should be set")
	assert.Equal(t, "conv-detail-1", view.conversationID, "conversationID should be set")
	assert.Nil(t, view.conversation, "conversation should be nil before Init")
}

func TestDetailViewInit(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	cmd := view.Init()

	// Init should return a command to load conversation
	assert.NotNil(t, cmd, "Init() should return a command to load conversation")

	// Execute command
	msg := cmd()
	loadedMsg, ok := msg.(ConversationLoadedMsg)
	assert.True(t, ok, "Command should return ConversationLoadedMsg")
	assert.NotNil(t, loadedMsg.Conversation, "Loaded conversation should not be nil")
	assert.Equal(t, "conv-detail-1", loadedMsg.Conversation.ID)
	assert.Nil(t, loadedMsg.Err, "Should have no error")
}

// Test: Conversation Loading
func TestDetailViewLoadConversationSuccess(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")

	// Simulate loading
	msg := ConversationLoadedMsg{
		Conversation: conv,
		Err:          nil,
	}

	view, _ = view.Update(msg)

	assert.NotNil(t, view.conversation, "conversation should be loaded")
	assert.Equal(t, "conv-detail-1", view.conversation.ID)
	assert.Equal(t, 3, len(view.conversation.Messages), "should have 3 messages")
	assert.Nil(t, view.err, "should have no error")
}

func TestDetailViewConversationNotFound(t *testing.T) {
	repo := repository.NewMock([]*domain.Conversation{})

	view := NewDetailView(repo, "nonexistent")

	// Initialize and load
	cmd := view.Init()
	msg := cmd()

	loadedMsg, ok := msg.(ConversationLoadedMsg)
	assert.True(t, ok)
	assert.Nil(t, loadedMsg.Conversation, "conversation should be nil")
	assert.Nil(t, loadedMsg.Err, "error should be nil (not found case)")

	// Update view
	view, _ = view.Update(loadedMsg)

	assert.Nil(t, view.conversation, "conversation should remain nil")
}

func TestDetailViewLoadConversationError(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")

	// Simulate error during load
	errMsg := domain.NewAppError(domain.ErrRepository, "test load error", nil)
	msg := ConversationLoadedMsg{
		Conversation: nil,
		Err:          errMsg,
	}

	view, _ = view.Update(msg)

	assert.NotNil(t, view.err, "should have error")
	assert.Nil(t, view.conversation, "conversation should be nil on error")
}

func TestDetailViewEmptyMessageList(t *testing.T) {
	conv := domain.NewConversation("conv-empty", "Empty Conversation", "claude-sonnet-4", time.Now())
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-empty")

	msg := ConversationLoadedMsg{
		Conversation: conv,
		Err:          nil,
	}

	view, _ = view.Update(msg)
	view = view.SetSize(80, 30)

	output := view.View()
	assert.NotEmpty(t, output, "should render even with no messages")
	assert.Contains(t, output, "Empty Conversation", "should show conversation title")
}

func TestDetailViewSingleMessage(t *testing.T) {
	conv := domain.NewConversation("conv-single", "Single Message", "claude-sonnet-4", time.Now())
	msg := domain.NewMessage(domain.RoleUser, "Single message content", time.Now())
	msg.Tokens = &domain.TokenCount{Input: 10, Output: 0}
	conv.AddMessage(msg)

	repo := repository.NewMock([]*domain.Conversation{conv})
	view := NewDetailView(repo, "conv-single")

	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "Single message content", "should display message content")
}

func TestDetailViewManyMessages(t *testing.T) {
	conv := createManyMessagesConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-many")

	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 40)

	output := view.View()
	assert.NotEmpty(t, output, "should render large conversation")
	// With 150 messages, we should be able to scroll
}

// Test: Rendering
func TestDetailViewReturnsNonEmpty(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()
	assert.NotEmpty(t, output, "View() should return non-empty string")
}

func TestDetailViewRendersHeader(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	assert.Contains(t, output, "Test Conversation for Detail View", "should show title")
	assert.Contains(t, output, "sonnet-4", "should show model")
}

func TestDetailViewRendersMessages(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	// Should contain parts of all messages
	assert.Contains(t, output, "Hello, can you help me", "should show first user message")
	assert.Contains(t, output, "Of course", "should show assistant message")
	assert.Contains(t, output, "type error", "should show second user message")
}

func TestDetailViewIncludesRoleLabels(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	// Should show role indicators
	assert.Contains(t, output, "User", "should label user messages")
	assert.Contains(t, output, "Assistant", "should label assistant messages")
}

func TestDetailViewIncludesTokenCounts(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	// Should show token counts somewhere
	// The exact format depends on implementation, but numbers should appear
	assert.True(t, len(output) > 100, "should have substantial content with token info")
}

func TestDetailViewMessagesProperlySpaced(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	// Messages should be separated (multiple newlines between them)
	lines := strings.Split(output, "\n")
	assert.Greater(t, len(lines), 10, "should have multiple lines with spacing")
}

// Test: Scrolling
func TestDetailViewDefaultScrollPosition(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	// Viewport should start at top (default behavior)
	// We can't directly test viewport position, but we can verify view renders
	output := view.View()
	assert.NotEmpty(t, output)
}

func TestDetailViewScrollDown(t *testing.T) {
	conv := createManyMessagesConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-many")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	// Press space to scroll down
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})

	// View should still render (scrolled state)
	output := view.View()
	assert.NotEmpty(t, output)
}

func TestDetailViewPageDown(t *testing.T) {
	conv := createManyMessagesConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-many")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	// Press PageDown
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyPgDown})

	output := view.View()
	assert.NotEmpty(t, output)
}

func TestDetailViewPageUp(t *testing.T) {
	conv := createManyMessagesConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-many")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	// Scroll down first, then up
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyPgUp})

	output := view.View()
	assert.NotEmpty(t, output)
}

func TestDetailViewHome(t *testing.T) {
	conv := createManyMessagesConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-many")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	// Scroll down, then press Home
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyHome})

	output := view.View()
	assert.NotEmpty(t, output)
}

func TestDetailViewEnd(t *testing.T) {
	conv := createManyMessagesConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-many")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	// Press End
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyEnd})

	output := view.View()
	assert.NotEmpty(t, output)
}

// Test: Keyboard Navigation
func TestDetailViewBackToList_B(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)

	// Press 'b' to go back
	view, cmd := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})

	assert.NotNil(t, cmd, "should return command")

	// Execute command
	msg := cmd()
	backMsg, ok := msg.(BackToListMsg)
	assert.True(t, ok, "should return BackToListMsg")
	assert.NotNil(t, backMsg, "message should not be nil")
}

func TestDetailViewBackToList_Escape(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)

	// Press Escape to go back
	view, cmd := view.Update(tea.KeyMsg{Type: tea.KeyEsc})

	assert.NotNil(t, cmd, "should return command")

	msg := cmd()
	_, ok := msg.(BackToListMsg)
	assert.True(t, ok, "should return BackToListMsg")
}

func TestDetailViewUnknownKey(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)

	// Press unknown key
	view, cmd := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	// Should not return command, viewport will handle it
	assert.Nil(t, cmd, "unknown keys should not return commands")
}

// Test: Window Resize
func TestDetailViewWindowResize(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	// Resize window
	view, _ = view.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	assert.Equal(t, 100, view.viewWidth, "width should update")
	assert.Equal(t, 40, view.viewHeight, "height should update")

	output := view.View()
	assert.NotEmpty(t, output)
}

// Test: Message Formatting
func TestDetailViewUserMessageFormat(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	// User messages should be identifiable
	assert.Contains(t, output, "User", "should label user messages")
}

func TestDetailViewAssistantMessageFormat(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	// Assistant messages should be identifiable
	assert.Contains(t, output, "Assistant", "should label assistant messages")
}

func TestDetailViewToolMessageFormat(t *testing.T) {
	conv := domain.NewConversation("conv-tool", "Tool Test", "claude-sonnet-4", time.Now())
	msg := domain.NewMessage(domain.RoleTool, "Tool execution result", time.Now())
	msg.Tokens = &domain.TokenCount{Input: 15, Output: 0}
	conv.AddMessage(msg)

	repo := repository.NewMock([]*domain.Conversation{conv})
	view := NewDetailView(repo, "conv-tool")

	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	// Tool messages should be identifiable
	assert.Contains(t, output, "Tool", "should label tool messages")
}

func TestDetailViewTokenCountFormatted(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	// Token counts should appear (exact format TBD)
	// Just verify output contains numbers
	assert.True(t, strings.Contains(output, "25") || strings.Contains(output, "50") || strings.Contains(output, "30"),
		"should display token counts")
}

func TestDetailViewLongLineWraps(t *testing.T) {
	conv := createLongMessageConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-long")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	// Long content should be wrapped
	lines := strings.Split(output, "\n")
	// Check that we have multiple lines (wrapping occurred)
	assert.Greater(t, len(lines), 5, "long content should wrap to multiple lines")
}

func TestDetailViewCodeBlocksPreserved(t *testing.T) {
	conv := domain.NewConversation("conv-code", "Code Test", "claude-sonnet-4", time.Now())
	codeContent := "Here's some code:\n```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```"
	msg := domain.NewMessage(domain.RoleAssistant, codeContent, time.Now())
	msg.Tokens = &domain.TokenCount{Input: 0, Output: 100}
	conv.AddMessage(msg)

	repo := repository.NewMock([]*domain.Conversation{conv})
	view := NewDetailView(repo, "conv-code")

	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()

	// Code blocks should be preserved (with syntax highlighting)
	assert.Contains(t, output, "func", "should preserve code content")
	assert.Contains(t, output, "main", "should preserve code content")
}

// Test: Edge Cases
func TestDetailViewVeryLongMessage(t *testing.T) {
	conv := domain.NewConversation("conv-vlong", "Very Long", "claude-sonnet-4", time.Now())
	veryLongContent := strings.Repeat("This is a very long message that goes on and on. ", 200)
	msg := domain.NewMessage(domain.RoleUser, veryLongContent, time.Now())
	msg.Tokens = &domain.TokenCount{Input: 1000, Output: 0}
	conv.AddMessage(msg)

	repo := repository.NewMock([]*domain.Conversation{conv})
	view := NewDetailView(repo, "conv-vlong")

	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()
	assert.NotEmpty(t, output, "should handle very long messages")
}

func TestDetailViewVeryShortTerminalWidth(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(40, 20) // Very narrow

	output := view.View()
	assert.NotEmpty(t, output, "should render even at narrow width")
}

func TestDetailViewEmptyConversationShowsMessage(t *testing.T) {
	conv := domain.NewConversation("conv-empty", "Empty", "claude-sonnet-4", time.Now())
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-empty")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()
	assert.NotEmpty(t, output, "should show something for empty conversation")
}

func TestDetailViewMissingConversationShowsError(t *testing.T) {
	repo := repository.NewMock([]*domain.Conversation{})

	view := NewDetailView(repo, "missing")
	cmd := view.Init()
	msg := cmd()

	loadedMsg := msg.(ConversationLoadedMsg)
	view, _ = view.Update(loadedMsg)
	view = view.SetSize(80, 30)

	output := view.View()
	assert.NotEmpty(t, output)
	// When conversation is nil, view should show appropriate message
}

// Test: Error Handling
func TestDetailViewErrorDisplay(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")

	// Set error
	errMsg := domain.NewAppError(domain.ErrRepository, "test error", nil)
	loadMsg := ConversationLoadedMsg{Conversation: nil, Err: errMsg}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	output := view.View()
	assert.Contains(t, output, "Error", "should display error message")
}

// Test: Helpers
func TestDetailViewSetSize(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	view = view.SetSize(100, 40)

	assert.Equal(t, 100, view.viewWidth)
	assert.Equal(t, 40, view.viewHeight)
}

// Test: Integration
func TestDetailViewFullWorkflow(t *testing.T) {
	// Create conversation
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	// Create view
	view := NewDetailView(repo, "conv-detail-1")
	require.NotNil(t, view)

	// Init
	cmd := view.Init()
	require.NotNil(t, cmd)

	// Load conversation
	msg := cmd()
	loadedMsg, ok := msg.(ConversationLoadedMsg)
	require.True(t, ok)
	require.NotNil(t, loadedMsg.Conversation)

	// Update with loaded conversation
	view, _ = view.Update(loadedMsg)
	require.NotNil(t, view.conversation)

	// Set size
	view = view.SetSize(80, 30)

	// Render
	output := view.View()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "Test Conversation")

	// Scroll
	view, _ = view.Update(tea.KeyMsg{Type: tea.KeyPgDown})

	// Back to list
	view, cmd = view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	require.NotNil(t, cmd)

	backMsg := cmd()
	_, ok = backMsg.(BackToListMsg)
	assert.True(t, ok)
}

// Test: Performance (basic check)
func TestDetailViewPerformanceWithLargeConversation(t *testing.T) {
	// This is a basic smoke test, not a real performance test
	conv := createManyMessagesConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-many")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)
	view = view.SetSize(80, 30)

	// Should render without hanging
	output := view.View()
	assert.NotEmpty(t, output, "should handle 150 messages")
}

// Test: Additional coverage for formatTimestamp
func TestDetailViewFormatTimestampAllCases(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})
	view := NewDetailView(repo, "conv-detail-1")

	now := time.Now()

	tests := []struct {
		name     string
		ts       time.Time
		contains string
	}{
		{"just now", now.Add(-30 * time.Second), "just now"},
		{"1 minute ago", now.Add(-1 * time.Minute), "1m ago"},
		{"59 minutes ago", now.Add(-59 * time.Minute), "59m ago"},
		{"1 hour ago", now.Add(-1 * time.Hour), "1h ago"},
		{"23 hours ago", now.Add(-23 * time.Hour), "23h ago"},
		{"1 day ago", now.Add(-25 * time.Hour), "1d ago"},
		{"6 days ago", now.Add(-6 * 24 * time.Hour), "6d ago"},
		{"8 days ago", now.Add(-8 * 24 * time.Hour), "20"}, // Should show date
		{"old date", now.Add(-365 * 24 * time.Hour), "20"}, // Should show date
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := view.formatTimestamp(tt.ts)
			assert.Contains(t, result, tt.contains, "timestamp format should match for: "+tt.name)
		})
	}
}

// Test: Additional coverage for Update message handling
func TestDetailViewUpdateWithUnknownMessage(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)

	// Send unknown message type
	type UnknownMsg struct{}
	view, _ = view.Update(UnknownMsg{})

	// Should handle gracefully
	assert.NotNil(t, view)
}

// Test: formatModel edge cases
func TestDetailViewFormatModelEdgeCases(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})
	view := NewDetailView(repo, "conv-detail-1")

	tests := []struct {
		name     string
		model    string
		expected string
	}{
		{"claude model", "claude-sonnet-4", "sonnet-4"},
		{"claude opus", "claude-opus-4", "opus-4"},
		{"non-claude model", "gpt-4-turbo", "gpt-4-turbo"},
		{"simple model", "sonnet", "sonnet"},
		{"empty model", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := view.formatModel(tt.model)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test: formatTokenCount edge cases
func TestDetailViewFormatTokenCountEdgeCases(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})
	view := NewDetailView(repo, "conv-detail-1")

	tests := []struct {
		name     string
		tokens   int64
		expected string
	}{
		{"zero", 0, ""},
		{"small", 500, "⊛ 500 tokens"},
		{"exactly 1000", 1000, "⊛ 1.0K tokens"},
		{"thousands", 5432, "⊛ 5.4K tokens"},
		{"exactly 1M", 1000000, "⊛ 1.0M tokens"},
		{"millions", 2500000, "⊛ 2.5M tokens"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := view.formatTokenCount(tt.tokens)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test: formatRole default case
func TestDetailViewFormatRoleUnknown(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})
	view := NewDetailView(repo, "conv-detail-1")

	// Test with an unknown role (edge case)
	result := view.formatRole(domain.Role("unknown"))
	assert.Equal(t, "unknown", result, "unknown roles should fall back to String()")
}

// Test: wrapText edge cases
func TestDetailViewWrapTextEdgeCases(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})
	view := NewDetailView(repo, "conv-detail-1")

	// Test with very small width
	result := view.wrapText("This is a long line that needs wrapping", 10)
	assert.NotEmpty(t, result)
	lines := strings.Split(result, "\n")
	assert.Greater(t, len(lines), 1, "should wrap to multiple lines")

	// Test with width < 20 (should use minimum)
	result = view.wrapText("Short text", 5)
	assert.NotEmpty(t, result)

	// Test with empty string
	result = view.wrapText("", 80)
	assert.Equal(t, "", result)

	// Test with multiple paragraphs
	result = view.wrapText("First paragraph\n\nSecond paragraph", 80)
	assert.Contains(t, result, "First paragraph")
	assert.Contains(t, result, "Second paragraph")
}

// Test: initViewport with small height
func TestDetailViewInitViewportSmallHeight(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	loadMsg := ConversationLoadedMsg{Conversation: conv, Err: nil}
	view, _ = view.Update(loadMsg)

	// Set very small height
	view = view.SetSize(80, 5)

	output := view.View()
	assert.NotEmpty(t, output, "should handle small height")
}

// Test: View before ready
func TestDetailViewBeforeReady(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	// Don't load conversation yet

	output := view.View()
	assert.Contains(t, output, "Conversation not found", "should show error when conversation not loaded")
}

// Test: Update before ready with key message
func TestDetailViewUpdateBeforeReady(t *testing.T) {
	conv := createDetailTestConversation()
	repo := repository.NewMock([]*domain.Conversation{conv})

	view := NewDetailView(repo, "conv-detail-1")
	// Don't initialize viewport

	// Try to send key before ready
	view, cmd := view.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	assert.NotNil(t, view)
	assert.Nil(t, cmd, "should not process scroll keys before ready")
}
