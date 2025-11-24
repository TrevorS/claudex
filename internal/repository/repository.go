// ABOUTME: Defines the Repository interface for conversation data access abstraction
//
// The repository implements a lazy loading strategy:
//   - List() scans the filesystem and returns metadata-only conversations (fast)
//   - GetByID() performs full JSONL parsing to load complete message data (slower, on-demand)
//
// This enables browsing 1000+ conversations efficiently without loading all message bodies at startup.
package repository

import (
	"context"

	"github.com/TrevorS/claudex/internal/domain"
)

// Repository defines the contract for accessing conversation data.
// Implementations must handle lazy loading: List() returns metadata only,
// while GetByID() loads complete conversation data.
type Repository interface {
	// List returns all conversations with metadata only (no message bodies).
	// This should be fast, only reading the first and last lines of each JSONL file.
	// Message bodies are loaded lazily via GetByID().
	List(ctx context.Context) ([]*domain.Conversation, error)

	// GetByID retrieves a specific conversation by ID and loads all messages.
	// This does the full JSONL parsing and message construction.
	GetByID(ctx context.Context, id string) (*domain.Conversation, error)

	// Search filters conversations by query string.
	// This may use full-text search or simple substring matching.
	// Returns conversations matching the query.
	Search(ctx context.Context, query string) ([]*domain.Conversation, error)

	// GetStatistics returns aggregated statistics across all conversations.
	// This includes total conversation count, message count, token usage, etc.
	GetStatistics(ctx context.Context) (*domain.Statistics, error)
}
