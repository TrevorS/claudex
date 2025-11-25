// ABOUTME: Performance benchmarks for ListView rendering with various dataset sizes.
package views

import (
	"testing"
	"time"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/TrevorS/claudex/internal/repository"
)

// BenchmarkListViewRender_Small benchmarks rendering with a small dataset (10 conversations)
func BenchmarkListViewRender_Small(b *testing.B) {
	conversations := createBenchmarkConversations(10)
	repo := repository.NewMock(conversations)
	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSize(100, 30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.View()
	}
}

// BenchmarkListViewRender_Medium benchmarks rendering with a medium dataset (100 conversations)
func BenchmarkListViewRender_Medium(b *testing.B) {
	conversations := createBenchmarkConversations(100)
	repo := repository.NewMock(conversations)
	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSize(100, 30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.View()
	}
}

// BenchmarkListViewRender_Large benchmarks rendering with a large dataset (1000 conversations)
func BenchmarkListViewRender_Large(b *testing.B) {
	conversations := createBenchmarkConversations(1000)
	repo := repository.NewMock(conversations)
	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSize(100, 30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = view.View()
	}
}

// BenchmarkListViewUpdate_Navigation benchmarks keyboard navigation updates
func BenchmarkListViewUpdate_Navigation(b *testing.B) {
	conversations := createBenchmarkConversations(100)
	repo := repository.NewMock(conversations)
	view := NewListView(repo)
	view = view.SetConversations(conversations)
	view = view.SetSize(100, 30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Alternate up/down navigation
		if i%2 == 0 {
			view = view.moveSelection(1) // move down
		} else {
			view = view.moveSelection(-1) // move up
		}
	}
}

// BenchmarkListViewSorting benchmarks conversation sorting operations
func BenchmarkListViewSorting(b *testing.B) {
	conversations := createBenchmarkConversations(500)
	repo := repository.NewMock(conversations)
	view := NewListView(repo)
	view = view.SetConversations(conversations)

	sortFields := []string{"timestamp", "title", "messages", "tokens"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field := sortFields[i%len(sortFields)]
		view = view.SetSort(field, i%2 == 0)
		_ = view.getSortedConversations()
	}
}

// Helper function to create benchmark conversations
func createBenchmarkConversations(count int) []*domain.Conversation {
	conversations := make([]*domain.Conversation, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		conv := domain.NewConversation(
			"conv-"+string(rune(i)),
			"Benchmark Conversation "+string(rune(i)),
			"claude-sonnet-4",
			now.Add(-time.Duration(i)*time.Hour),
		)

		// Add a few messages
		for j := 0; j < 3; j++ {
			msg := domain.NewMessage(domain.RoleUser, "Message "+string(rune(j)), now)
			msg.Tokens = &domain.TokenCount{Input: 100, Output: 50}
			conv.AddMessage(msg)
		}

		conversations[i] = conv
	}

	return conversations
}
