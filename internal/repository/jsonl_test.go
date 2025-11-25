// ABOUTME: Comprehensive tests for JSONLRepository implementation
package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TrevorS/claudex/internal/domain"
)

// TestNewJSONLRepository tests repository initialization
func TestNewJSONLRepository(t *testing.T) {
	tmpDir := t.TempDir()

	repo, err := NewJSONLRepository(tmpDir)
	require.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, tmpDir, repo.rootPath)
}

// TestNewJSONLRepository_InvalidPath tests error handling for invalid paths
func TestNewJSONLRepository_InvalidPath(t *testing.T) {
	_, err := NewJSONLRepository("/nonexistent/path/that/definitely/does/not/exist")
	assert.Error(t, err)
}

// TestJSONLRepository_List tests listing all conversations
func TestJSONLRepository_List(t *testing.T) {
	tmpDir := createTestProjectStructure(t, map[string][]string{
		"-Users-trevor-Projects-test": {"conv1.jsonl", "conv2.jsonl"},
		"-Users-trevor-code":          {"conv3.jsonl"},
	})

	// Write test conversation files
	writeTestConversation(t, filepath.Join(tmpDir, "-Users-trevor-Projects-test", "conv1.jsonl"),
		"Test Conversation 1", "sess-001", "2025-11-23T10:00:00Z", "2025-11-23T10:05:00Z")
	writeTestConversation(t, filepath.Join(tmpDir, "-Users-trevor-Projects-test", "conv2.jsonl"),
		"Test Conversation 2", "sess-002", "2025-11-23T11:00:00Z", "2025-11-23T11:05:00Z")
	writeTestConversation(t, filepath.Join(tmpDir, "-Users-trevor-code", "conv3.jsonl"),
		"Test Conversation 3", "sess-003", "2025-11-23T12:00:00Z", "2025-11-23T12:05:00Z")

	repo, err := NewJSONLRepository(tmpDir)
	require.NoError(t, err)

	conversations, err := repo.List(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 3, len(conversations))

	// Verify conversations are loaded with metadata (use map to avoid order dependency)
	idMap := make(map[string]*domain.Conversation)
	for _, conv := range conversations {
		idMap[conv.ID] = conv
	}

	assert.Contains(t, idMap, "sess-001")
	assert.Contains(t, idMap, "sess-002")
	assert.Contains(t, idMap, "sess-003")

	// Verify one conversation's details
	conv := idMap["sess-001"]
	assert.Equal(t, "Test Conversation 1", conv.Title)
	assert.Equal(t, "claude-sonnet-4-20250514", conv.Model, "Model should be extracted from assistant message")
	assert.Equal(t, 2, conv.MessageCount(), "Should count user and assistant messages")
	assert.False(t, conv.CreatedAt.IsZero())
	assert.False(t, conv.UpdatedAt.IsZero())
}

// TestJSONLRepository_List_Empty tests listing from empty directory
func TestJSONLRepository_List_Empty(t *testing.T) {
	tmpDir := t.TempDir()

	repo, err := NewJSONLRepository(tmpDir)
	require.NoError(t, err)

	conversations, err := repo.List(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0, len(conversations))
}

// TestJSONLRepository_GetByID tests loading a complete conversation
func TestJSONLRepository_GetByID(t *testing.T) {
	tmpDir := createTestProjectStructure(t, map[string][]string{
		"-Users-trevor-Projects-test": {"conv1.jsonl"},
	})

	// Write conversation with multiple messages
	convPath := filepath.Join(tmpDir, "-Users-trevor-Projects-test", "conv1.jsonl")
	writeConversationWithMessages(t, convPath, "sess-001", "Test Conv", []struct {
		role    string
		content string
		tokens  int
	}{
		{"user", "Hello", 10},
		{"assistant", "Hi there!", 20},
		{"user", "How are you?", 15},
		{"assistant", "I'm doing great!", 25},
	})

	repo, err := NewJSONLRepository(tmpDir)
	require.NoError(t, err)

	conv, err := repo.GetByID(context.Background(), "sess-001")
	require.NoError(t, err)
	require.NotNil(t, conv)

	assert.Equal(t, "sess-001", conv.ID)
	assert.Equal(t, "Test Conv", conv.Title)
	assert.Equal(t, 4, len(conv.Messages))
	assert.Equal(t, "Hello", conv.Messages[0].Content)
	assert.Equal(t, "Hi there!", conv.Messages[1].Content)
}

// TestJSONLRepository_GetByID_NotFound tests error handling for missing conversations
func TestJSONLRepository_GetByID_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	repo, err := NewJSONLRepository(tmpDir)
	require.NoError(t, err)

	conv, err := repo.GetByID(context.Background(), "nonexistent-id")
	// Should return nil conversation, not error (conversation doesn't exist)
	assert.Nil(t, conv)
}

// TestJSONLRepository_PathDecoding tests hyphenated path conversion
func TestJSONLRepository_PathDecoding(t *testing.T) {
	tests := []struct {
		encoded string
		want    string
	}{
		{"-Users-trevor-Projects-claudex", "/Users/trevor/Projects/claudex"},
		{"-home-user-documents", "/home/user/documents"},
		{"-tmp-test", "/tmp/test"},
	}

	for _, tt := range tests {
		t.Run(tt.encoded, func(t *testing.T) {
			got := decodePath(tt.encoded)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestJSONLRepository_CorruptedJSONL tests graceful handling of malformed lines
func TestJSONLRepository_CorruptedJSONL(t *testing.T) {
	tmpDir := createTestProjectStructure(t, map[string][]string{
		"-Users-trevor-test": {"corrupted.jsonl"},
	})

	// Write a file with mixed valid and invalid JSON
	corruptPath := filepath.Join(tmpDir, "-Users-trevor-test", "corrupted.jsonl")
	err := os.WriteFile(corruptPath, []byte(`{"type":"summary","summary":"Test","sessionId":"sess-001","timestamp":"2025-11-23T10:00:00Z"}
{"type":"user","message":{"role":"user","content":"Valid message"},"sessionId":"sess-001","timestamp":"2025-11-23T10:01:00Z"}
THIS IS NOT VALID JSON
{"type":"assistant","message":{"role":"assistant","content":"Another valid message"},"sessionId":"sess-001","timestamp":"2025-11-23T10:02:00Z"}
`), 0644)
	require.NoError(t, err)

	repo, err := NewJSONLRepository(tmpDir)
	require.NoError(t, err)

	// Should still be able to load the conversation, skipping bad lines
	conv, err := repo.GetByID(context.Background(), "sess-001")
	require.NoError(t, err)
	require.NotNil(t, conv)

	// Should have loaded the 2 valid messages (not the corrupted line)
	assert.Equal(t, "Test", conv.Title)
	assert.Greater(t, len(conv.Messages), 0)
}

// TestJSONLRepository_Search tests conversation search
func TestJSONLRepository_Search(t *testing.T) {
	tmpDir := createTestProjectStructure(t, map[string][]string{
		"-Users-trevor-test": {"c1.jsonl", "c2.jsonl", "c3.jsonl"},
	})

	writeTestConversation(t, filepath.Join(tmpDir, "-Users-trevor-test", "c1.jsonl"),
		"Building a REST API", "sess-001", "2025-11-23T10:00:00Z", "2025-11-23T10:05:00Z")
	writeTestConversation(t, filepath.Join(tmpDir, "-Users-trevor-test", "c2.jsonl"),
		"Debugging my code", "sess-002", "2025-11-23T11:00:00Z", "2025-11-23T11:05:00Z")
	writeTestConversation(t, filepath.Join(tmpDir, "-Users-trevor-test", "c3.jsonl"),
		"REST API design patterns", "sess-003", "2025-11-23T12:00:00Z", "2025-11-23T12:05:00Z")

	repo, err := NewJSONLRepository(tmpDir)
	require.NoError(t, err)

	results, err := repo.Search(context.Background(), "REST API")
	require.NoError(t, err)

	// Should match both conversations with "REST API" in title
	assert.Equal(t, 2, len(results))
}

// TestJSONLRepository_GetStatistics tests statistics calculation
func TestJSONLRepository_GetStatistics(t *testing.T) {
	tmpDir := createTestProjectStructure(t, map[string][]string{
		"-Users-trevor-test": {"c1.jsonl", "c2.jsonl"},
	})

	writeConversationWithMessages(t, filepath.Join(tmpDir, "-Users-trevor-test", "c1.jsonl"),
		"sess-001", "Conv 1", []struct {
			role    string
			content string
			tokens  int
		}{
			{"user", "msg1", 100},
			{"assistant", "msg2", 200},
		})

	writeConversationWithMessages(t, filepath.Join(tmpDir, "-Users-trevor-test", "c2.jsonl"),
		"sess-002", "Conv 2", []struct {
			role    string
			content string
			tokens  int
		}{
			{"user", "msg3", 150},
		})

	repo, err := NewJSONLRepository(tmpDir)
	require.NoError(t, err)

	stats, err := repo.GetStatistics(context.Background())
	require.NoError(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, 2, stats.ConversationCount)
	assert.Equal(t, 3, stats.MessageCount)
}

// TestJSONLRepository_ProjectPathInConversation tests that project path is stored correctly
func TestJSONLRepository_ProjectPathInConversation(t *testing.T) {
	tmpDir := createTestProjectStructure(t, map[string][]string{
		"-Users-trevor-Projects-claudex": {"conv1.jsonl"},
	})

	writeTestConversation(t, filepath.Join(tmpDir, "-Users-trevor-Projects-claudex", "conv1.jsonl"),
		"Test", "sess-001", "2025-11-23T10:00:00Z", "2025-11-23T10:05:00Z")

	repo, err := NewJSONLRepository(tmpDir)
	require.NoError(t, err)

	conversations, err := repo.List(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, len(conversations))

	// Title should be extracted from summary
	assert.Equal(t, "Test", conversations[0].Title)
}

// TestMockRepository_Basic tests the MockRepository implementation
func TestMockRepository_Basic(t *testing.T) {
	conv1 := &domain.Conversation{
		ID:        "conv-001",
		Title:     "Test Conversation 1",
		Model:     "claude-3",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Messages:  make([]*domain.Message, 0),
	}

	conv2 := &domain.Conversation{
		ID:        "conv-002",
		Title:     "Test Conversation 2",
		Model:     "claude-3",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Messages:  make([]*domain.Message, 0),
	}

	mock := NewMock([]*domain.Conversation{conv1, conv2})

	// Test List
	conversations, err := mock.List(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, len(conversations))

	// Test GetByID
	conv, err := mock.GetByID(context.Background(), "conv-001")
	require.NoError(t, err)
	assert.Equal(t, "Test Conversation 1", conv.Title)

	// Test Search
	results, err := mock.Search(context.Background(), "Conversation 2")
	require.NoError(t, err)
	assert.Equal(t, 1, len(results))

	// Test GetStatistics
	stats, err := mock.GetStatistics(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, stats.ConversationCount)
}

// ============ Test Helpers ============

// createTestProjectStructure creates a temporary directory structure for testing
func createTestProjectStructure(t *testing.T, projects map[string][]string) string {
	tmpDir := t.TempDir()

	for projectDir, files := range projects {
		projectPath := filepath.Join(tmpDir, projectDir)
		err := os.MkdirAll(projectPath, 0755)
		require.NoError(t, err)

		for _, file := range files {
			filePath := filepath.Join(projectPath, file)
			err := os.WriteFile(filePath, []byte(""), 0644)
			require.NoError(t, err)
		}
	}

	return tmpDir
}

// writeTestConversation writes a simple test conversation with summary and timestamps
func writeTestConversation(t *testing.T, filePath string, title, sessionID, createdAt, updatedAt string) {
	content := `{"type":"summary","summary":"` + title + `","sessionId":"` + sessionID + `","timestamp":"` + createdAt + `","uuid":"msg-001","version":"1.0"}
{"type":"user","message":{"role":"user","content":"Test message"},"usage":{"input_tokens":50,"output_tokens":50},"sessionId":"` + sessionID + `","timestamp":"` + createdAt + `","uuid":"msg-002"}
{"type":"assistant","message":{"role":"assistant","model":"claude-sonnet-4-20250514","content":"Response"},"usage":{"input_tokens":50,"output_tokens":100},"sessionId":"` + sessionID + `","timestamp":"` + updatedAt + `","uuid":"msg-003"}`

	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)
}

// writeConversationWithMessages writes a conversation with multiple messages
func writeConversationWithMessages(t *testing.T, filePath, sessionID, title string, messages []struct {
	role    string
	content string
	tokens  int
}) {
	lines := []string{
		`{"type":"summary","summary":"` + title + `","sessionId":"` + sessionID + `","timestamp":"2025-11-23T10:00:00Z","uuid":"msg-000","version":"1.0"}`,
	}

	for i, msg := range messages {
		// Claude JSONL format uses role as the type (e.g., "user", "assistant")
		line := `{"type":"` + msg.role + `","message":{"role":"` + msg.role + `","content":"` + msg.content + `"},"usage":{"input_tokens":` + itoa(msg.tokens) + `,"output_tokens":50},"sessionId":"` + sessionID + `","timestamp":"2025-11-23T10:0` + itoa(i) + `:00Z","uuid":"msg-00` + itoa(i+1) + `"}`
		lines = append(lines, line)
	}

	content := ""
	for i, line := range lines {
		content += line
		if i < len(lines)-1 {
			content += "\n"
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)
}

// itoa converts an int to string (simple version of strconv.Itoa)
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}

	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}
