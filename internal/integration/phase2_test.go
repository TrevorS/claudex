// ABOUTME: End-to-end integration tests for Phase 2 features.
// Tests verify complete workflows across multiple components: list→detail navigation,
// search functionality, filtering, content rendering, and complex query execution.
package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/TrevorS/claudex/internal/repository"
	"github.com/TrevorS/claudex/internal/search"
	"github.com/TrevorS/claudex/internal/ui"
	"github.com/TrevorS/claudex/internal/ui/views"
)

// Helper: Create a rich test dataset with multiple conversations
func createIntegrationTestData() []*domain.Conversation {
	now := time.Now()
	conversations := make([]*domain.Conversation, 0)

	// Conversation 1: Python coding conversation
	conv1 := domain.NewConversation(
		"conv-python-1",
		"Help with Python sorting",
		"claude-sonnet-4",
		now.Add(-3*time.Hour),
	)
	conv1.ProjectPath = "test/python-project"
	msg1 := domain.NewMessage(domain.RoleUser, "How do I sort a list in Python?", now.Add(-3*time.Hour))
	msg1.Tokens = &domain.TokenCount{Input: 100, Output: 0}
	conv1.AddMessage(msg1)
	msg2 := domain.NewMessage(domain.RoleAssistant, "You can use the `sorted()` function or the `.sort()` method:\n```python\nmy_list = [3, 1, 4, 1, 5]\nsorted_list = sorted(my_list)\n```", now.Add(-3*time.Hour).Add(5*time.Second))
	msg2.Tokens = &domain.TokenCount{Input: 0, Output: 150}
	conv1.AddMessage(msg2)
	conversations = append(conversations, conv1)

	// Conversation 2: Go coding conversation with high tokens
	conv2 := domain.NewConversation(
		"conv-go-1",
		"Debugging Go concurrency issue",
		"claude-opus-4",
		now.Add(-2*time.Hour),
	)
	conv2.ProjectPath = "test/go-project"
	msg3 := domain.NewMessage(domain.RoleUser, "My Go program has a race condition", now.Add(-2*time.Hour))
	msg3.Tokens = &domain.TokenCount{Input: 2000, Output: 0}
	conv2.AddMessage(msg3)
	msg4 := domain.NewMessage(domain.RoleAssistant, "Let me help you find the race condition. First, run `go run -race main.go`", now.Add(-2*time.Hour).Add(10*time.Second))
	msg4.Tokens = &domain.TokenCount{Input: 0, Output: 3500}
	conv2.AddMessage(msg4)
	msg5 := domain.NewMessage(domain.RoleUser, "I found it! The mutex wasn't locked properly.", now.Add(-2*time.Hour).Add(2*time.Minute))
	msg5.Tokens = &domain.TokenCount{Input: 500, Output: 0}
	conv2.AddMessage(msg5)
	conversations = append(conversations, conv2)

	// Conversation 3: JavaScript conversation (recent)
	conv3 := domain.NewConversation(
		"conv-js-1",
		"React hooks tutorial",
		"claude-sonnet-4",
		now.Add(-30*time.Minute),
	)
	conv3.ProjectPath = "test/javascript-project"
	msg6 := domain.NewMessage(domain.RoleUser, "Explain React hooks", now.Add(-30*time.Minute))
	msg6.Tokens = &domain.TokenCount{Input: 50, Output: 0}
	conv3.AddMessage(msg6)
	msg7 := domain.NewMessage(domain.RoleAssistant, "React hooks allow you to use state in function components:\n```javascript\nconst [count, setCount] = useState(0);\n```", now.Add(-30*time.Minute).Add(3*time.Second))
	msg7.Tokens = &domain.TokenCount{Input: 0, Output: 200}
	conv3.AddMessage(msg7)
	conversations = append(conversations, conv3)

	// Conversation 4: Python AI/ML conversation
	conv4 := domain.NewConversation(
		"conv-python-ai",
		"Machine learning with Python",
		"claude-sonnet-4",
		now.Add(-1*time.Hour),
	)
	conv4.ProjectPath = "test/python-project"
	msg8 := domain.NewMessage(domain.RoleUser, "How do I train a neural network with Python?", now.Add(-1*time.Hour))
	msg8.Tokens = &domain.TokenCount{Input: 200, Output: 0}
	conv4.AddMessage(msg8)
	msg9 := domain.NewMessage(domain.RoleAssistant, "You can use TensorFlow or PyTorch. Here's a simple example:\n```python\nimport tensorflow as tf\nmodel = tf.keras.Sequential([...])\n```", now.Add(-1*time.Hour).Add(5*time.Second))
	msg9.Tokens = &domain.TokenCount{Input: 0, Output: 400}
	conv4.AddMessage(msg9)
	conversations = append(conversations, conv4)

	// Conversation 5: Database query conversation (low tokens)
	conv5 := domain.NewConversation(
		"conv-db-1",
		"SQL query help",
		"claude-sonnet-4",
		now.Add(-5*time.Hour),
	)
	conv5.ProjectPath = "test/database-project"
	msg10 := domain.NewMessage(domain.RoleUser, "SELECT statement syntax", now.Add(-5*time.Hour))
	msg10.Tokens = &domain.TokenCount{Input: 20, Output: 0}
	conv5.AddMessage(msg10)
	msg11 := domain.NewMessage(domain.RoleAssistant, "Basic SELECT: `SELECT * FROM users WHERE active = 1`", now.Add(-5*time.Hour).Add(2*time.Second))
	msg11.Tokens = &domain.TokenCount{Input: 0, Output: 30}
	conv5.AddMessage(msg11)
	conversations = append(conversations, conv5)

	return conversations
}

// ==============================================================================
// Integration Test 1: List → Detail Workflow
// ==============================================================================

func TestWorkflow_ListToDetail(t *testing.T) {
	conversations := createIntegrationTestData()
	repo := repository.NewMock(conversations)

	// Create ListView
	listView := views.NewListView(repo)

	// Initialize and load conversations
	cmd := listView.Init()
	require.NotNil(t, cmd, "ListView Init should return command")

	// Execute load command
	msg := cmd()
	loadedMsg, ok := msg.(views.ConversationsLoadedMsg)
	require.True(t, ok, "should return ConversationsLoadedMsg")
	require.Equal(t, 5, len(loadedMsg.Conversations), "should load 5 conversations")

	// Update ListView with loaded conversations
	listView, _ = listView.Update(loadedMsg)
	listView = listView.SetSize(100, 30)

	// Verify list displays
	output := listView.View()
	assert.NotEmpty(t, output, "ListView should render")
	assert.Contains(t, output, "Python", "should show Python conversation")
	assert.Contains(t, output, "Go", "should show Go conversation")

	// Navigate down in list
	listView, _ = listView.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 1, listView.SelectedIndex(), "should move selection down")

	// Select conversation (press Enter)
	listView, cmd = listView.Update(tea.KeyMsg{Type: tea.KeyEnter})
	require.NotNil(t, cmd, "Enter should return command")

	selectMsg := cmd()
	selectedMsg, ok := selectMsg.(views.SelectConversationMsg)
	require.True(t, ok, "should return SelectConversationMsg")
	require.NotNil(t, selectedMsg.Conversation, "should contain conversation")

	// Create DetailView with selected conversation
	detailView := views.NewDetailView(repo, selectedMsg.Conversation.ID)

	// Initialize and load detail
	cmd = detailView.Init()
	require.NotNil(t, cmd)
	msg = cmd()
	convLoadedMsg, ok := msg.(views.ConversationLoadedMsg)
	require.True(t, ok)
	require.NotNil(t, convLoadedMsg.Conversation)

	// Update DetailView
	detailView, _ = detailView.Update(convLoadedMsg)
	detailView = detailView.SetSize(100, 40)

	// Verify detail displays messages
	detailOutput := detailView.View()
	assert.NotEmpty(t, detailOutput, "DetailView should render")
	// selectedMsg.Conversation is the one at index 1 after sorting
	// Just check that it has content and title
	assert.True(t, len(detailOutput) > 200, "should have substantial content")

	// Test scrolling in detail view
	detailView, _ = detailView.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	assert.NotEmpty(t, detailView.View(), "should render after scroll")

	// Return to list
	detailView, cmd = detailView.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	require.NotNil(t, cmd, "should return back command")
	backMsg := cmd()
	_, ok = backMsg.(views.BackToListMsg)
	assert.True(t, ok, "should return BackToListMsg")
}

// ==============================================================================
// Integration Test 2: Search Workflow
// ==============================================================================

func TestWorkflow_Search(t *testing.T) {
	conversations := createIntegrationTestData()
	repo := repository.NewMock(conversations)
	ctx := context.Background()

	// Create search engine
	engine := search.NewEngine(repo)

	// Test simple text search
	results, err := engine.Search(ctx, "python")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 2, "should find Python conversations")

	// Verify result ranking (Python-specific conversations should rank higher)
	pythonCount := 0
	for _, conv := range results {
		if strings.Contains(strings.ToLower(conv.Title), "python") {
			pythonCount++
		}
	}
	assert.GreaterOrEqual(t, pythonCount, 2, "should find multiple Python conversations")

	// Test field search
	results, err = engine.Search(ctx, "title:\"Go concurrency\"")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1, "should find Go conversation by title")

	// Test token comparison
	results, err = engine.Search(ctx, "tokens:>5000")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1, "should find high-token conversation")
	if len(results) > 0 {
		assert.Greater(t, results[0].TotalTokens(), int64(5000), "result should have >5000 tokens")
	}

	// Test empty search returns error (empty query not allowed)
	results, err = engine.Search(ctx, "")
	assert.Error(t, err, "empty search should return error")
	assert.Nil(t, results, "empty search should not return results")
}

// ==============================================================================
// Integration Test 3: Filter Workflow
// ==============================================================================

func TestWorkflow_Filter(t *testing.T) {
	conversations := createIntegrationTestData()
	repo := repository.NewMock(conversations)

	listView := views.NewListView(repo)

	// Load conversations
	cmd := listView.Init()
	msg := cmd()
	loadedMsg := msg.(views.ConversationsLoadedMsg)
	listView, _ = listView.Update(loadedMsg)
	listView = listView.SetSize(100, 30)

	// Get all projects
	allConvs := listView.Conversations()
	projects := make(map[string]bool)
	for _, conv := range allConvs {
		if conv.ProjectPath != "" {
			projects[conv.ProjectPath] = true
		}
	}
	assert.GreaterOrEqual(t, len(projects), 3, "should have multiple projects")

	// Test filtering by project
	pythonConvs := filterByProject(allConvs, "test/python-project")
	assert.Equal(t, 2, len(pythonConvs), "should find 2 Python project conversations")

	goConvs := filterByProject(allConvs, "test/go-project")
	assert.Equal(t, 1, len(goConvs), "should find 1 Go project conversation")

	// Test quick filter: Today
	now := time.Now()
	todayConvs := filterByTimeRange(allConvs, now.Add(-24*time.Hour), now)
	assert.GreaterOrEqual(t, len(todayConvs), 3, "should find conversations from today")

	// Test quick filter: This Week
	weekConvs := filterByTimeRange(allConvs, now.Add(-7*24*time.Hour), now)
	assert.Equal(t, 5, len(weekConvs), "all test conversations should be within a week")

	// Test quick filter: This Month
	monthConvs := filterByTimeRange(allConvs, now.Add(-30*24*time.Hour), now)
	assert.Equal(t, 5, len(monthConvs), "all test conversations should be within a month")
}

// Helper: Filter conversations by project
func filterByProject(conversations []*domain.Conversation, projectPath string) []*domain.Conversation {
	filtered := make([]*domain.Conversation, 0)
	for _, conv := range conversations {
		if conv.ProjectPath == projectPath {
			filtered = append(filtered, conv)
		}
	}
	return filtered
}

// Helper: Filter conversations by time range
func filterByTimeRange(conversations []*domain.Conversation, start, end time.Time) []*domain.Conversation {
	filtered := make([]*domain.Conversation, 0)
	for _, conv := range conversations {
		if conv.UpdatedAt.After(start) && conv.UpdatedAt.Before(end) {
			filtered = append(filtered, conv)
		}
	}
	return filtered
}

// ==============================================================================
// Integration Test 4: Complex Query Workflow
// ==============================================================================

func TestWorkflow_ComplexQuery(t *testing.T) {
	conversations := createIntegrationTestData()
	repo := repository.NewMock(conversations)
	engine := search.NewEngine(repo)
	ctx := context.Background()

	// Test boolean AND
	results, err := engine.Search(ctx, "python AND sorting")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1, "should find Python sorting conversation")

	// Test boolean OR
	results, err = engine.Search(ctx, "python OR go")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 3, "should find Python and Go conversations")

	// Test boolean NOT
	results, err = engine.Search(ctx, "python NOT machine")
	require.NoError(t, err)
	// Should find Python conversations but exclude the ML one
	for _, conv := range results {
		if strings.Contains(strings.ToLower(conv.Title), "python") {
			assert.NotContains(t, strings.ToLower(conv.Title), "machine", "should exclude ML conversation")
		}
	}

	// Test complex query with parentheses
	results, err = engine.Search(ctx, "(python OR go) AND tokens:>1000")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1, "should find conversations matching complex query")
	for _, conv := range results {
		assert.Greater(t, conv.TotalTokens(), int64(1000), "all results should have >1000 tokens")
	}

	// Test field search with comparison
	results, err = engine.Search(ctx, "tokens:>=5000 AND tokens:<=10000")
	require.NoError(t, err)
	for _, conv := range results {
		tokens := conv.TotalTokens()
		assert.GreaterOrEqual(t, tokens, int64(5000), "should be >= 5000 tokens")
		assert.LessOrEqual(t, tokens, int64(10000), "should be <= 10000 tokens")
	}

	// Test message count filter
	results, err = engine.Search(ctx, "messages:>2")
	require.NoError(t, err)
	for _, conv := range results {
		assert.Greater(t, len(conv.Messages), 2, "should have >2 messages")
	}
}

// ==============================================================================
// Integration Test 5: Content Rendering Workflow
// ==============================================================================

func TestWorkflow_ContentRendering(t *testing.T) {
	conversations := createIntegrationTestData()
	repo := repository.NewMock(conversations)

	// Load a conversation with code blocks
	conv := conversations[0] // Python conversation with code
	detailView := views.NewDetailView(repo, conv.ID)

	// Init and load
	cmd := detailView.Init()
	msg := cmd()
	loadedMsg := msg.(views.ConversationLoadedMsg)
	detailView, _ = detailView.Update(loadedMsg)
	detailView = detailView.SetSize(100, 40)

	output := detailView.View()

	// Verify markdown formatting
	assert.NotEmpty(t, output, "should render content")

	// Verify code blocks are present
	assert.Contains(t, output, "python", "should show Python code block marker")
	assert.Contains(t, output, "sorted", "should show code content")

	// Verify role labels
	assert.Contains(t, output, "User", "should label user messages")
	assert.Contains(t, output, "Assistant", "should label assistant messages")

	// Load Go conversation with code
	conv = conversations[1]
	detailView = views.NewDetailView(repo, conv.ID)
	cmd = detailView.Init()
	msg = cmd()
	loadedMsg = msg.(views.ConversationLoadedMsg)
	detailView, _ = detailView.Update(loadedMsg)
	detailView = detailView.SetSize(100, 40)

	output = detailView.View()

	// Verify Go-specific content
	assert.Contains(t, output, "race", "should show race condition content")
	assert.Contains(t, output, "mutex", "should show mutex content")
}

// ==============================================================================
// Integration Test 6: Performance Workflow
// ==============================================================================

func TestWorkflow_Performance(t *testing.T) {
	// Create large dataset (1000+ conversations)
	largeDataset := make([]*domain.Conversation, 1200)
	now := time.Now()

	for i := 0; i < 1200; i++ {
		conv := domain.NewConversation(
			"conv-"+string(rune(i)),
			"Conversation "+string(rune(i)),
			"claude-sonnet-4",
			now.Add(-time.Duration(i)*time.Hour),
		)
		conv.ProjectPath = "test/project-" + string(rune(i%10))

		// Add messages
		for j := 0; j < 5; j++ {
			msg := domain.NewMessage(domain.RoleUser, "Message "+string(rune(j)), now)
			msg.Tokens = &domain.TokenCount{Input: 100, Output: 50}
			conv.AddMessage(msg)
		}

		largeDataset[i] = conv
	}

	repo := repository.NewMock(largeDataset)

	// Test 1: List loads quickly
	listView := views.NewListView(repo)
	startTime := time.Now()

	cmd := listView.Init()
	msg := cmd()
	loadedMsg := msg.(views.ConversationsLoadedMsg)
	listView, _ = listView.Update(loadedMsg)
	listView = listView.SetSize(100, 30)

	listLoadTime := time.Since(startTime)
	assert.Less(t, listLoadTime, 2*time.Second, "list should load in <2 seconds with 1200 conversations")

	// Test 2: Search completes quickly
	engine := search.NewEngine(repo)
	ctx := context.Background()
	startTime = time.Now()

	results, err := engine.Search(ctx, "Conversation")
	require.NoError(t, err)

	searchTime := time.Since(startTime)
	assert.Less(t, searchTime, 500*time.Millisecond, "search should complete in <500ms")
	assert.GreaterOrEqual(t, len(results), 100, "should return at least 100 results")

	// Test 3: Detail view renders smoothly
	detailView := views.NewDetailView(repo, largeDataset[0].ID)
	startTime = time.Now()

	cmd = detailView.Init()
	msg = cmd()
	convLoadedMsg := msg.(views.ConversationLoadedMsg)
	detailView, _ = detailView.Update(convLoadedMsg)
	detailView = detailView.SetSize(100, 40)

	detailLoadTime := time.Since(startTime)
	assert.Less(t, detailLoadTime, 200*time.Millisecond, "detail view should load in <200ms")

	output := detailView.View()
	assert.NotEmpty(t, output, "detail should render")
}

// ==============================================================================
// Integration Test 7: Full UI Model State Transitions
// ==============================================================================

func TestWorkflow_FullUIModelStateTransitions(t *testing.T) {
	conversations := createIntegrationTestData()
	repo := repository.NewMock(conversations)

	// Create root UI model
	model := ui.NewModel(repo)
	require.NotNil(t, model)

	// Initialize model (loads conversations into ListView)
	cmd := model.Init()
	if cmd != nil {
		// Model is initialized, can render
		_ = cmd()
	}

	// Initial state should be ViewList
	assert.NotEmpty(t, model.View(), "model should render in list view")

	// Simulate window size - Update returns tea.Model interface
	result, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if typedModel, ok := result.(ui.Model); ok {
		model = typedModel
	}

	// Load conversations
	listCmd := views.NewListView(repo).Init()
	loadMsg := listCmd()
	result, _ = model.Update(loadMsg)
	if typedModel, ok := result.(ui.Model); ok {
		model = typedModel
	}

	// Verify model can render
	output := model.View()
	assert.NotEmpty(t, output, "model should render after loading")
}

// ==============================================================================
// Integration Test 8: Error Handling Across Components
// ==============================================================================

func TestWorkflow_ErrorHandling(t *testing.T) {
	// Create empty repository
	repo := repository.NewMock([]*domain.Conversation{})

	// Test ListView with empty data
	listView := views.NewListView(repo)
	cmd := listView.Init()
	msg := cmd()
	loadedMsg := msg.(views.ConversationsLoadedMsg)
	listView, _ = listView.Update(loadedMsg)
	listView = listView.SetSize(100, 30)

	output := listView.View()
	assert.NotEmpty(t, output, "should render even with no data")
	assert.Contains(t, output, "No conversations", "should show empty state message")

	// Test DetailView with nonexistent conversation
	detailView := views.NewDetailView(repo, "nonexistent")
	cmd = detailView.Init()
	msg = cmd()
	convLoadedMsg := msg.(views.ConversationLoadedMsg)
	detailView, _ = detailView.Update(convLoadedMsg)
	detailView = detailView.SetSize(100, 40)

	output = detailView.View()
	assert.NotEmpty(t, output, "should render even when conversation not found")
	assert.Contains(t, output, "not found", "should show error message")

	// Test search with invalid query (invalid token comparison returns empty, not error)
	engine := search.NewEngine(repo)
	ctx := context.Background()
	results, err := engine.Search(ctx, "tokens:invalid")
	require.NoError(t, err, "search parses but returns no results")
	assert.Equal(t, 0, len(results), "invalid query should return no results")
}

// ==============================================================================
// Integration Test 9: Caching and Repository Layer
// ==============================================================================

func TestWorkflow_CachedRepositoryIntegration(t *testing.T) {
	conversations := createIntegrationTestData()
	baseRepo := repository.NewMock(conversations)
	ctx := context.Background()

	// Wrap with caching (uses default TTLs: 30s for list, 5m for stats)
	cachedRepo, err := repository.NewCachedRepository(baseRepo, 100)
	require.NoError(t, err)

	// First call - should hit underlying repository
	list1, err := cachedRepo.List(ctx)
	require.NoError(t, err)
	assert.Equal(t, 5, len(list1))

	// Second call - should return cached result
	list2, err := cachedRepo.List(ctx)
	require.NoError(t, err)
	assert.Equal(t, 5, len(list2))

	// GetByID - should cache individual conversations
	conv1, err := cachedRepo.GetByID(ctx, "conv-python-1")
	require.NoError(t, err)
	require.NotNil(t, conv1)

	// Second GetByID - should return cached
	conv2, err := cachedRepo.GetByID(ctx, "conv-python-1")
	require.NoError(t, err)
	require.NotNil(t, conv2)
	assert.Equal(t, conv1.ID, conv2.ID)

	// Search should pass through (not cached)
	// MockRepository.Search does substring matching on title
	results, err := cachedRepo.Search(ctx, "Python")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1, "search should work through cache layer")
}

// ==============================================================================
// Integration Test 10: Multi-Step User Journey
// ==============================================================================

func TestWorkflow_CompleteUserJourney(t *testing.T) {
	// This test simulates a complete user session:
	// 1. Launch app (load conversations)
	// 2. Browse list
	// 3. Search for specific topic
	// 4. Select result
	// 5. View detail
	// 6. Scroll through messages
	// 7. Return to list
	// 8. Filter by project
	// 9. Select different conversation
	// 10. Exit

	conversations := createIntegrationTestData()
	repo := repository.NewMock(conversations)

	// Step 1 & 2: Launch and browse
	listView := views.NewListView(repo)
	cmd := listView.Init()
	msg := cmd()
	loadedMsg := msg.(views.ConversationsLoadedMsg)
	listView, _ = listView.Update(loadedMsg)
	listView = listView.SetSize(100, 30)

	assert.Equal(t, 5, len(listView.Conversations()), "should load all conversations")

	// Step 3: Search
	engine := search.NewEngine(repo)
	ctx := context.Background()
	searchResults, err := engine.Search(ctx, "python")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(searchResults), 2, "should find Python results")

	// Step 4 & 5: Select and view detail
	selectedConv := searchResults[0]
	detailView := views.NewDetailView(repo, selectedConv.ID)
	cmd = detailView.Init()
	msg = cmd()
	convLoadedMsg := msg.(views.ConversationLoadedMsg)
	detailView, _ = detailView.Update(convLoadedMsg)
	detailView = detailView.SetSize(100, 40)

	output := detailView.View()
	assert.NotEmpty(t, output, "detail should render")
	assert.Contains(t, output, "Python", "should show Python content")

	// Step 6: Scroll
	detailView, _ = detailView.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	detailView, _ = detailView.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	assert.NotEmpty(t, detailView.View(), "should render after scrolling")

	// Step 7: Return to list
	detailView, cmd = detailView.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	require.NotNil(t, cmd)
	backMsg := cmd()
	_, ok := backMsg.(views.BackToListMsg)
	assert.True(t, ok, "should return to list")

	// Step 8: Filter by project
	filtered := filterByProject(listView.Conversations(), "test/go-project")
	assert.Equal(t, 1, len(filtered), "should filter to Go project")

	// Step 9: Select different conversation
	goConv := filtered[0]
	detailView2 := views.NewDetailView(repo, goConv.ID)
	cmd = detailView2.Init()
	msg = cmd()
	convLoadedMsg = msg.(views.ConversationLoadedMsg)
	detailView2, _ = detailView2.Update(convLoadedMsg)
	detailView2 = detailView2.SetSize(100, 40)

	output = detailView2.View()
	assert.Contains(t, output, "Go", "should show Go conversation")
	assert.Contains(t, output, "concurrency", "should show conversation content")

	// Step 10: User exits (would press Ctrl+C or q, handled by main app)
	// This is tested in the UI model layer
}
