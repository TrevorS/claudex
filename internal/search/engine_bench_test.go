// ABOUTME: Performance benchmarks for search engine with various query types and dataset sizes.
package search

import (
	"context"
	"testing"
	"time"

	"github.com/TrevorS/claudex/internal/domain"
)

// BenchmarkSearch_SimpleText_Small benchmarks simple text search on small dataset (100 conversations)
func BenchmarkSearch_SimpleText_Small(b *testing.B) {
	repo := createBenchmarkRepo(100)
	engine := NewEngine(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Search(ctx, "python")
	}
}

// BenchmarkSearch_SimpleText_Medium benchmarks simple text search on medium dataset (500 conversations)
func BenchmarkSearch_SimpleText_Medium(b *testing.B) {
	repo := createBenchmarkRepo(500)
	engine := NewEngine(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Search(ctx, "python")
	}
}

// BenchmarkSearch_SimpleText_Large benchmarks simple text search on large dataset (1000 conversations)
func BenchmarkSearch_SimpleText_Large(b *testing.B) {
	repo := createBenchmarkRepo(1000)
	engine := NewEngine(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Search(ctx, "python")
	}
}

// BenchmarkSearch_FieldSearch benchmarks field-specific searches
func BenchmarkSearch_FieldSearch(b *testing.B) {
	repo := createBenchmarkRepo(500)
	engine := NewEngine(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Search(ctx, "title:python")
	}
}

// BenchmarkSearch_BooleanAND benchmarks boolean AND queries
func BenchmarkSearch_BooleanAND(b *testing.B) {
	repo := createBenchmarkRepo(500)
	engine := NewEngine(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Search(ctx, "python AND tutorial")
	}
}

// BenchmarkSearch_BooleanOR benchmarks boolean OR queries
func BenchmarkSearch_BooleanOR(b *testing.B) {
	repo := createBenchmarkRepo(500)
	engine := NewEngine(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Search(ctx, "python OR rust")
	}
}

// BenchmarkSearch_ComplexQuery benchmarks complex queries with multiple operators
func BenchmarkSearch_ComplexQuery(b *testing.B) {
	repo := createBenchmarkRepo(500)
	engine := NewEngine(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Search(ctx, "(python OR rust) AND tokens:>1000")
	}
}

// BenchmarkSearch_TokenComparison benchmarks token comparison queries
func BenchmarkSearch_TokenComparison(b *testing.B) {
	repo := createBenchmarkRepo(500)
	engine := NewEngine(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Search(ctx, "tokens:>5000")
	}
}

// BenchmarkParser_SimpleQuery benchmarks query parsing for simple queries
func BenchmarkParser_SimpleQuery(b *testing.B) {
	query := "python programming"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(query)
		_, _ = parser.Parse()
	}
}

// BenchmarkParser_ComplexQuery benchmarks query parsing for complex queries
func BenchmarkParser_ComplexQuery(b *testing.B) {
	query := "(python OR rust) AND (tokens:>1000 AND tokens:<10000) NOT archived"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(query)
		_, _ = parser.Parse()
	}
}

// BenchmarkFilter_TermMatch benchmarks term filtering
func BenchmarkFilter_TermMatch(b *testing.B) {
	conversations := createBenchmarkConversations(1000)
	node := &TermNode{Term: "python"}
	filter := NewNodeFilter(node)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, conv := range conversations {
			_ = filter.Match(conv)
		}
	}
}

// BenchmarkFilter_FieldMatch benchmarks field filtering
func BenchmarkFilter_FieldMatch(b *testing.B) {
	conversations := createBenchmarkConversations(1000)
	node := &FieldNode{Field: "title", Op: ":", Value: "python"}
	filter := NewNodeFilter(node)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, conv := range conversations {
			_ = filter.Match(conv)
		}
	}
}

// BenchmarkRanking benchmarks relevance ranking
func BenchmarkRanking(b *testing.B) {
	conversations := createBenchmarkConversations(1000)
	repo := &MockRepository{conversations: conversations}
	engine := NewEngine(repo)
	query := "python tutorial"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, conv := range conversations {
			_ = engine.rankResult(conv, query)
		}
	}
}

// Helper: Create benchmark repository with specified number of conversations
func createBenchmarkRepo(count int) *MockRepository {
	return &MockRepository{
		conversations: createBenchmarkConversations(count),
	}
}

// Helper: Create benchmark conversations
func createBenchmarkConversations(count int) []*domain.Conversation {
	conversations := make([]*domain.Conversation, count)
	now := time.Now()

	topics := []string{"python", "rust", "go", "javascript", "tutorial", "guide", "help"}
	models := []string{"claude-sonnet-4", "claude-opus-4", "claude-haiku-4"}

	for i := 0; i < count; i++ {
		topic := topics[i%len(topics)]
		model := models[i%len(models)]

		conv := domain.NewConversation(
			"conv-"+string(rune(i)),
			"Benchmark "+topic+" conversation "+string(rune(i)),
			model,
			now.Add(-time.Duration(i)*time.Hour),
		)

		// Add messages with varied content and tokens
		for j := 0; j < 5; j++ {
			role := domain.RoleUser
			if j%2 == 1 {
				role = domain.RoleAssistant
			}

			content := "Message about " + topic + " with some additional content"
			msg := domain.NewMessage(role, content, now.Add(-time.Duration(i)*time.Hour+time.Duration(j)*time.Minute))
			msg.Tokens = &domain.TokenCount{
				Input:  int64((i*10 + j*50) % 1000),
				Output: int64((i*15 + j*75) % 2000),
			}
			conv.AddMessage(msg)
		}

		conversations[i] = conv
	}

	return conversations
}
