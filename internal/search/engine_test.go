// ABOUTME: engine_test.go contains comprehensive tests for the search engine,
// including query execution, relevance ranking, and error handling.
package search

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRepository implements a simple in-memory repository for testing.
type MockRepository struct {
	conversations []*domain.Conversation
}

func (r *MockRepository) List(ctx context.Context) ([]*domain.Conversation, error) {
	return r.conversations, nil
}

func (r *MockRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {
	for _, c := range r.conversations {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}

func (r *MockRepository) Search(ctx context.Context, query string) ([]*domain.Conversation, error) {
	// Not used by engine (engine implements its own search)
	return nil, nil
}

func (r *MockRepository) GetStatistics(ctx context.Context) (*domain.Statistics, error) {
	return nil, nil
}

// Helper to create test conversations
func makeTestConv(id, title, model string, content []string, tokens int64) *domain.Conversation {
	conv := domain.NewConversation(id, title, model, time.Now())
	for _, c := range content {
		msg := domain.NewMessage(domain.RoleUser, c, time.Now())
		msg.Tokens = &domain.TokenCount{Input: tokens / int64(len(content))}
		conv.AddMessage(msg)
	}
	return conv
}

func TestEngine_SimpleTextSearch(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Python Tutorial", "sonnet", []string{"Learn python basics"}, 100),
			makeTestConv("2", "Rust Guide", "sonnet", []string{"Learn rust programming"}, 100),
			makeTestConv("3", "Python Advanced", "sonnet", []string{"Advanced python topics"}, 100),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "python")

	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "1", results[0].ID)
	assert.Equal(t, "3", results[1].ID)
}

func TestEngine_FieldSearch(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "My Chat", "claude-sonnet-4", []string{"hello"}, 100),
			makeTestConv("2", "Another Chat", "claude-opus-3", []string{"hello"}, 100),
			makeTestConv("3", "Third Chat", "claude-sonnet-4", []string{"hello"}, 100),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "model:sonnet")

	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestEngine_BooleanAND(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Python Tutorial", "sonnet", []string{"Learn python and rust"}, 100),
			makeTestConv("2", "Python Guide", "sonnet", []string{"Only python here"}, 100),
			makeTestConv("3", "Rust Guide", "sonnet", []string{"Only rust here"}, 100),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "python AND rust")

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "1", results[0].ID)
}

func TestEngine_BooleanOR(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Python Tutorial", "sonnet", []string{"Learn python"}, 100),
			makeTestConv("2", "Rust Guide", "sonnet", []string{"Learn rust"}, 100),
			makeTestConv("3", "JavaScript Guide", "sonnet", []string{"Learn javascript"}, 100),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "python OR rust")

	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestEngine_BooleanNOT(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Python Tutorial", "sonnet", []string{"Learn python"}, 100),
			makeTestConv("2", "Python Archived", "sonnet", []string{"Old archived python content"}, 100),
			makeTestConv("3", "Python Active", "sonnet", []string{"New python content"}, 100),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "python AND NOT archived")

	require.NoError(t, err)
	assert.Len(t, results, 2)
	// Should not include ID "2" which has "archived"
	for _, r := range results {
		assert.NotEqual(t, "2", r.ID)
	}
}

func TestEngine_ComplexQuery(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "AI Chat", "claude-sonnet-4", []string{"Machine learning with python"}, 100),
			makeTestConv("2", "ML Discussion", "claude-sonnet-4", []string{"Deep learning basics"}, 100),
			makeTestConv("3", "Python Guide", "claude-opus-3", []string{"Python programming"}, 100),
			makeTestConv("4", "Old Chat", "claude-sonnet-4", []string{"Archived content"}, 100),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "(ai OR ml) AND model:sonnet AND NOT archived")

	require.NoError(t, err)
	assert.Len(t, results, 2)
	// Should match IDs 1 and 2 (have ai/ml, use sonnet, not archived)
}

func TestEngine_Ranking(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Some Chat", "sonnet", []string{"python is mentioned here in content"}, 100),
			makeTestConv("2", "Python Tutorial", "sonnet", []string{"This is about something else"}, 100),
			makeTestConv("3", "Random Chat", "sonnet", []string{"Nothing relevant"}, 100),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "python")

	require.NoError(t, err)
	assert.Len(t, results, 2)
	// ID "2" should rank higher because "python" is in the title (3x weight)
	assert.Equal(t, "2", results[0].ID, "Title match should rank highest")
	assert.Equal(t, "1", results[1].ID, "Content match should rank lower")
}

func TestEngine_LimitResults(t *testing.T) {
	// Create 150 conversations
	convs := make([]*domain.Conversation, 150)
	for i := 0; i < 150; i++ {
		convs[i] = makeTestConv(
			string(rune(i)),
			"Test Chat",
			"sonnet",
			[]string{"test content"},
			100,
		)
	}

	repo := &MockRepository{conversations: convs}
	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "test")

	require.NoError(t, err)
	assert.LessOrEqual(t, len(results), 100, "Should limit to 100 results")
}

func TestEngine_EmptyQuery(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Chat", "sonnet", []string{"content"}, 100),
		},
	}

	engine := NewEngine(repo)
	_, err := engine.Search(context.Background(), "")

	assert.Error(t, err, "Empty query should return error")
}

func TestEngine_InvalidQuery(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Chat", "sonnet", []string{"content"}, 100),
		},
	}

	engine := NewEngine(repo)
	_, err := engine.Search(context.Background(), "(unclosed")

	assert.Error(t, err, "Invalid query should return error")
}

func TestEngine_NoResults(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Python Chat", "sonnet", []string{"python content"}, 100),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "javascript")

	require.NoError(t, err)
	assert.Len(t, results, 0, "Should return empty slice for no matches")
}

func TestEngine_TokenComparison(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Small Chat", "sonnet", []string{"content"}, 1000),
			makeTestConv("2", "Medium Chat", "sonnet", []string{"content"}, 5000),
			makeTestConv("3", "Large Chat", "sonnet", []string{"content"}, 10000),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "tokens:>4000")

	require.NoError(t, err)
	assert.Len(t, results, 2)
	// Should return IDs 2 and 3
}

func TestEngine_MultipleTermsImplicitAND(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Chat", "sonnet", []string{"python and rust programming"}, 100),
			makeTestConv("2", "Chat", "sonnet", []string{"only python here"}, 100),
			makeTestConv("3", "Chat", "sonnet", []string{"only rust here"}, 100),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "python rust")

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "1", results[0].ID)
}

func TestEngine_RankingWeights(t *testing.T) {
	repo := &MockRepository{
		conversations: []*domain.Conversation{
			makeTestConv("1", "Content match", "sonnet", []string{"test appears in content"}, 100),
			makeTestConv("2", "Field match", "test-model", []string{"nothing"}, 100),
			makeTestConv("3", "Test Title", "sonnet", []string{"nothing"}, 100),
		},
	}

	engine := NewEngine(repo)
	results, err := engine.Search(context.Background(), "test")

	require.NoError(t, err)
	assert.Len(t, results, 3)
	// Title (ID 3) should rank first (3x), field (ID 2) second (2x), content (ID 1) third (1x)
	assert.Equal(t, "3", results[0].ID, "Title match should rank highest")
	assert.Equal(t, "2", results[1].ID, "Field match should rank second")
	assert.Equal(t, "1", results[2].ID, "Content match should rank lowest")
}

func TestEngine_Performance(t *testing.T) {
	// Create 1000 conversations
	convs := make([]*domain.Conversation, 1000)
	for i := 0; i < 1000; i++ {
		content := "test content"
		if i%10 == 0 {
			content = "python programming tutorial with lots of content " + strings.Repeat("text ", 100)
		}
		convs[i] = makeTestConv(
			string(rune(i)),
			"Test Chat",
			"sonnet",
			[]string{content},
			100,
		)
	}

	repo := &MockRepository{conversations: convs}
	engine := NewEngine(repo)

	start := time.Now()
	results, err := engine.Search(context.Background(), "python")
	duration := time.Since(start)

	require.NoError(t, err)
	assert.Greater(t, len(results), 0)
	// Should complete in reasonable time (< 1 second for 1000 conversations)
	assert.Less(t, duration, 1*time.Second, "Search should be fast even with 1000 conversations")
}
