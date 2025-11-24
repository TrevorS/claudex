// ABOUTME: Implements a caching wrapper around the Repository interface using LRU caching.
//
// CachedRepository wraps any Repository implementation and provides caching for:
//   - GetByID() - LRU cache for fully-loaded conversations (max 50, no TTL - immutable)
//   - List() - Single entry cache with 30s TTL (metadata-only conversations)
//   - GetStatistics() - Single entry cache with 5m TTL (expensive aggregation)
//   - Search() - Pass through without caching (results vary per query)
//
// Thread-safe for concurrent access using sync.RWMutex.
package repository

import (
	"context"
	"sync"
	"time"

	"github.com/TrevorS/claudex/internal/domain"
	lru "github.com/hashicorp/golang-lru/v2"
)

// CachedRepository wraps a Repository with LRU caching for performance.
type CachedRepository struct {
	wrapped Repository
	mu      sync.RWMutex

	// Individual conversation cache (LRU, max 50, no TTL - conversations are immutable)
	conversationCache *lru.Cache[string, *domain.Conversation]
	// In-flight request tracking for GetByID to prevent duplicate fetches
	conversationInflight map[string]*sync.Mutex

	// Full list cache (single entry, 30s TTL)
	listCache        []*domain.Conversation
	listCacheTime    time.Time
	listTTL          time.Duration
	listInflightLock sync.Mutex

	// Statistics cache (single entry, 5m TTL)
	statsCache        *domain.Statistics
	statsCacheTime    time.Time
	statsTTL          time.Duration
	statsInflightLock sync.Mutex
}

// NewCachedRepository creates a new caching wrapper around the given repository.
// maxConversations controls the LRU cache size for conversation objects.
func NewCachedRepository(wrapped Repository, maxConversations int) (*CachedRepository, error) {
	if wrapped == nil {
		return nil, &domain.AppError{
			Type:    domain.ErrValidation,
			Message: "wrapped repository cannot be nil",
		}
	}

	if maxConversations <= 0 {
		return nil, &domain.AppError{
			Type:    domain.ErrValidation,
			Message: "maxConversations must be positive",
		}
	}

	// Create LRU cache for conversations
	conversationCache, err := lru.New[string, *domain.Conversation](maxConversations)
	if err != nil {
		return nil, &domain.AppError{
			Type:    domain.ErrInternal,
			Message: "failed to create conversation cache",
			Err:     err,
		}
	}

	return &CachedRepository{
		wrapped:              wrapped,
		conversationCache:    conversationCache,
		conversationInflight: make(map[string]*sync.Mutex),
		listCache:            nil,
		listTTL:              30 * time.Second,
		statsCache:           nil,
		statsTTL:             5 * time.Minute,
	}, nil
}

// List returns all conversations with metadata only (no message bodies).
// Results are cached for 30 seconds to reduce filesystem access.
func (c *CachedRepository) List(ctx context.Context) ([]*domain.Conversation, error) {
	// Check cache first with read lock
	c.mu.RLock()
	if c.listCache != nil && time.Since(c.listCacheTime) < c.listTTL {
		cached := c.listCache
		c.mu.RUnlock()
		return cached, nil
	}
	c.mu.RUnlock()

	// Use inflight lock to prevent thundering herd
	c.listInflightLock.Lock()
	defer c.listInflightLock.Unlock()

	// Double-check cache after acquiring lock (another goroutine may have loaded it)
	c.mu.RLock()
	if c.listCache != nil && time.Since(c.listCacheTime) < c.listTTL {
		cached := c.listCache
		c.mu.RUnlock()
		return cached, nil
	}
	c.mu.RUnlock()

	// Cache miss or expired - fetch from wrapped repository
	convs, err := c.wrapped.List(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache with write lock
	c.mu.Lock()
	c.listCache = convs
	c.listCacheTime = time.Now()
	c.mu.Unlock()

	return convs, nil
}

// GetByID retrieves a conversation by ID.
// Conversations are cached in an LRU cache with no TTL (conversations are immutable).
func (c *CachedRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {
	// Check cache first with read lock
	c.mu.RLock()
	if conv, ok := c.conversationCache.Get(id); ok {
		c.mu.RUnlock()
		return conv, nil
	}
	c.mu.RUnlock()

	// Get or create inflight lock for this ID
	c.mu.Lock()
	inflightMu, exists := c.conversationInflight[id]
	if !exists {
		inflightMu = &sync.Mutex{}
		c.conversationInflight[id] = inflightMu
	}
	c.mu.Unlock()

	// Lock this specific conversation ID to prevent duplicate fetches
	inflightMu.Lock()
	defer inflightMu.Unlock()

	// Double-check cache after acquiring lock (another goroutine may have loaded it)
	c.mu.RLock()
	if conv, ok := c.conversationCache.Get(id); ok {
		c.mu.RUnlock()
		return conv, nil
	}
	c.mu.RUnlock()

	// Cache miss - fetch from wrapped repository
	conv, err := c.wrapped.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache with write lock (will evict LRU if full)
	c.mu.Lock()
	c.conversationCache.Add(id, conv)
	// Clean up inflight lock for this ID
	delete(c.conversationInflight, id)
	c.mu.Unlock()

	return conv, nil
}

// Search filters conversations by query string.
// This is passed through to the wrapped repository without caching since
// search results vary by query.
func (c *CachedRepository) Search(ctx context.Context, query string) ([]*domain.Conversation, error) {
	return c.wrapped.Search(ctx, query)
}

// GetStatistics returns aggregated statistics.
// Results are cached for 5 minutes to avoid expensive re-aggregation.
func (c *CachedRepository) GetStatistics(ctx context.Context) (*domain.Statistics, error) {
	// Check cache first with read lock
	c.mu.RLock()
	if c.statsCache != nil && time.Since(c.statsCacheTime) < c.statsTTL {
		cached := c.statsCache
		c.mu.RUnlock()
		return cached, nil
	}
	c.mu.RUnlock()

	// Use inflight lock to prevent thundering herd
	c.statsInflightLock.Lock()
	defer c.statsInflightLock.Unlock()

	// Double-check cache after acquiring lock (another goroutine may have loaded it)
	c.mu.RLock()
	if c.statsCache != nil && time.Since(c.statsCacheTime) < c.statsTTL {
		cached := c.statsCache
		c.mu.RUnlock()
		return cached, nil
	}
	c.mu.RUnlock()

	// Cache miss or expired - compute from wrapped repository
	stats, err := c.wrapped.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache with write lock
	c.mu.Lock()
	c.statsCache = stats
	c.statsCacheTime = time.Now()
	c.mu.Unlock()

	return stats, nil
}

// Invalidate clears the list and statistics caches.
// Conversation cache is NOT cleared since conversations are immutable after loading.
func (c *CachedRepository) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.listCache = nil
	c.statsCache = nil
	// Don't clear conversation cache - conversations are immutable
}
