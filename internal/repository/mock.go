// ABOUTME: Mock implementation of Repository for testing without real file I/O
package repository

import (
	"context"

	"github.com/TrevorS/claudex/internal/domain"
)

// MockRepository is a test double that implements the Repository interface.
// It stores pre-loaded conversation data for testing without file system access.
type MockRepository struct {
	conversations []*domain.Conversation
	statistics    *domain.Statistics
}

// NewMock creates a new MockRepository with the provided conversations.
func NewMock(conversations []*domain.Conversation) *MockRepository {
	return &MockRepository{
		conversations: conversations,
	}
}

// NewMockWithStatistics creates a MockRepository with both conversations and pre-calculated statistics.
func NewMockWithStatistics(conversations []*domain.Conversation, stats *domain.Statistics) *MockRepository {
	return &MockRepository{
		conversations: conversations,
		statistics:    stats,
	}
}

// List returns all stored conversations.
func (m *MockRepository) List(ctx context.Context) ([]*domain.Conversation, error) {
	return m.conversations, nil
}

// GetByID finds and returns a conversation by ID, or nil if not found.
func (m *MockRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {
	for _, conv := range m.conversations {
		if conv.ID == id {
			return conv, nil
		}
	}
	return nil, nil
}

// Search filters conversations by simple substring matching on title.
func (m *MockRepository) Search(ctx context.Context, query string) ([]*domain.Conversation, error) {
	var results []*domain.Conversation
	for _, conv := range m.conversations {
		// Simple substring matching on title for testing
		if contains(conv.Title, query) {
			results = append(results, conv)
		}
	}
	return results, nil
}

// GetStatistics returns pre-calculated statistics or calculates them on demand.
func (m *MockRepository) GetStatistics(ctx context.Context) (*domain.Statistics, error) {
	if m.statistics != nil {
		return m.statistics, nil
	}

	// Calculate on demand if not provided
	stats := domain.NewStatistics()
	stats.ConversationCount = len(m.conversations)

	for _, conv := range m.conversations {
		stats.MessageCount += len(conv.Messages)
		stats.TotalTokens += conv.TotalTokens()
	}

	return stats, nil
}

// contains is a helper for substring matching (case-sensitive for simplicity)
func contains(haystack, needle string) bool {
	if needle == "" {
		return true
	}
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
