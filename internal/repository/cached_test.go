// ABOUTME: Comprehensive tests for the CachedRepository wrapper implementation.
//
// Tests cover:
//   - Constructor validation
//   - LRU cache hit/miss behavior for GetByID()
//   - List caching with TTL (30s)
//   - Statistics caching with TTL (5m)
//   - Cache invalidation
//   - Thread safety under concurrent access
//   - LRU eviction when max size is reached
//   - Error propagation from wrapped repository
//   - Performance benchmarks
package repository

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepository is a simple mock for testing the cached wrapper.
type mockRepository struct {
	listFunc          func(ctx context.Context) ([]*domain.Conversation, error)
	getByIDFunc       func(ctx context.Context, id string) (*domain.Conversation, error)
	searchFunc        func(ctx context.Context, query string) ([]*domain.Conversation, error)
	getStatisticsFunc func(ctx context.Context) (*domain.Statistics, error)

	// Call counters for verification
	listCallCount       atomic.Int32
	getByIDCallCount    atomic.Int32
	searchCallCount     atomic.Int32
	statisticsCallCount atomic.Int32
}

func (m *mockRepository) List(ctx context.Context) ([]*domain.Conversation, error) {
	m.listCallCount.Add(1)
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return nil, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {
	m.getByIDCallCount.Add(1)
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) Search(ctx context.Context, query string) ([]*domain.Conversation, error) {
	m.searchCallCount.Add(1)
	if m.searchFunc != nil {
		return m.searchFunc(ctx, query)
	}
	return nil, nil
}

func (m *mockRepository) GetStatistics(ctx context.Context) (*domain.Statistics, error) {
	m.statisticsCallCount.Add(1)
	if m.getStatisticsFunc != nil {
		return m.getStatisticsFunc(ctx)
	}
	return nil, nil
}

func TestNewCachedRepository(t *testing.T) {
	t.Run("creates wrapper successfully", func(t *testing.T) {
		mock := &mockRepository{}
		cached, err := NewCachedRepository(mock, 100)

		require.NoError(t, err)
		assert.NotNil(t, cached)
		assert.Equal(t, mock, cached.wrapped)
	})

	t.Run("rejects nil wrapped repository", func(t *testing.T) {
		cached, err := NewCachedRepository(nil, 100)

		require.Error(t, err)
		assert.Nil(t, cached)

		appErr, ok := err.(*domain.AppError)
		require.True(t, ok, "expected domain.AppError")
		assert.Equal(t, domain.ErrValidation, appErr.Type)
		assert.Contains(t, appErr.Message, "cannot be nil")
	})

	t.Run("rejects zero maxConversations", func(t *testing.T) {
		mock := &mockRepository{}
		cached, err := NewCachedRepository(mock, 0)

		require.Error(t, err)
		assert.Nil(t, cached)

		appErr, ok := err.(*domain.AppError)
		require.True(t, ok, "expected domain.AppError")
		assert.Equal(t, domain.ErrValidation, appErr.Type)
		assert.Contains(t, appErr.Message, "must be positive")
	})

	t.Run("rejects negative maxConversations", func(t *testing.T) {
		mock := &mockRepository{}
		cached, err := NewCachedRepository(mock, -10)

		require.Error(t, err)
		assert.Nil(t, cached)

		appErr, ok := err.(*domain.AppError)
		require.True(t, ok, "expected domain.AppError")
		assert.Equal(t, domain.ErrValidation, appErr.Type)
		assert.Contains(t, appErr.Message, "must be positive")
	})
}

func TestCachedRepository_List(t *testing.T) {
	t.Run("first call fetches from wrapped repository", func(t *testing.T) {
		now := time.Now()
		conv := domain.NewConversation("test-123", "Test Conversation", "claude-sonnet-4", now)
		mock := &mockRepository{
			listFunc: func(ctx context.Context) ([]*domain.Conversation, error) {
				return []*domain.Conversation{conv}, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()
		result, err := cached.List(ctx)

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "test-123", result[0].ID)
		assert.Equal(t, int32(1), mock.listCallCount.Load())
	})

	t.Run("second call within TTL returns cached result", func(t *testing.T) {
		now := time.Now()
		conv := domain.NewConversation("test-123", "Test Conversation", "claude-sonnet-4", now)
		mock := &mockRepository{
			listFunc: func(ctx context.Context) ([]*domain.Conversation, error) {
				return []*domain.Conversation{conv}, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()

		// First call
		result1, err := cached.List(ctx)
		require.NoError(t, err)
		require.Len(t, result1, 1)
		assert.Equal(t, int32(1), mock.listCallCount.Load())

		// Second call immediately (within TTL)
		result2, err := cached.List(ctx)
		require.NoError(t, err)
		require.Len(t, result2, 1)

		// Should still only have called wrapped repository once
		assert.Equal(t, int32(1), mock.listCallCount.Load())

		// Results should be the same
		assert.Equal(t, result1[0].ID, result2[0].ID)
	})

	t.Run("call after TTL fetches fresh data", func(t *testing.T) {
		now := time.Now()
		conv1 := domain.NewConversation("test-123", "Test Conversation", "claude-sonnet-4", now)
		conv2 := domain.NewConversation("test-456", "Another Conversation", "claude-sonnet-4", now)

		callCount := 0
		mock := &mockRepository{
			listFunc: func(ctx context.Context) ([]*domain.Conversation, error) {
				callCount++
				if callCount == 1 {
					return []*domain.Conversation{conv1}, nil
				}
				return []*domain.Conversation{conv1, conv2}, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		// Override TTL to 1ms for testing
		cached.listTTL = 1 * time.Millisecond

		ctx := context.Background()

		// First call
		result1, err := cached.List(ctx)
		require.NoError(t, err)
		require.Len(t, result1, 1)

		// Wait for TTL to expire
		time.Sleep(2 * time.Millisecond)

		// Second call after TTL
		result2, err := cached.List(ctx)
		require.NoError(t, err)
		require.Len(t, result2, 2)

		// Should have called wrapped repository twice
		assert.Equal(t, int32(2), mock.listCallCount.Load())
	})

	t.Run("error propagation from wrapped repository", func(t *testing.T) {
		expectedErr := &domain.AppError{Type: domain.ErrRepository, Message: "disk error"}
		mock := &mockRepository{
			listFunc: func(ctx context.Context) ([]*domain.Conversation, error) {
				return nil, expectedErr
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()
		result, err := cached.List(ctx)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestCachedRepository_GetByID(t *testing.T) {
	t.Run("first call fetches from wrapped repository", func(t *testing.T) {
		now := time.Now()
		conv := domain.NewConversation("test-123", "Test Conversation", "claude-sonnet-4", now)
		mock := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {
				if id == "test-123" {
					return conv, nil
				}
				return nil, &domain.AppError{Type: domain.ErrNotFound, Message: "not found"}
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()
		result, err := cached.GetByID(ctx, "test-123")

		require.NoError(t, err)
		assert.Equal(t, "test-123", result.ID)
		assert.Equal(t, int32(1), mock.getByIDCallCount.Load())
	})

	t.Run("second call returns cached result", func(t *testing.T) {
		now := time.Now()
		conv := domain.NewConversation("test-123", "Test Conversation", "claude-sonnet-4", now)
		mock := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {
				return conv, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()

		// First call
		result1, err := cached.GetByID(ctx, "test-123")
		require.NoError(t, err)
		assert.Equal(t, int32(1), mock.getByIDCallCount.Load())

		// Second call
		result2, err := cached.GetByID(ctx, "test-123")
		require.NoError(t, err)

		// Should still only have called wrapped repository once
		assert.Equal(t, int32(1), mock.getByIDCallCount.Load())

		// Results should be the same object
		assert.Equal(t, result1.ID, result2.ID)
	})

	t.Run("different IDs are cached separately", func(t *testing.T) {
		now := time.Now()
		conv1 := domain.NewConversation("test-123", "Test 1", "claude-sonnet-4", now)
		conv2 := domain.NewConversation("test-456", "Test 2", "claude-sonnet-4", now)

		mock := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {
				if id == "test-123" {
					return conv1, nil
				}
				if id == "test-456" {
					return conv2, nil
				}
				return nil, &domain.AppError{Type: domain.ErrNotFound, Message: "not found"}
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()

		// Load both conversations
		result1, err := cached.GetByID(ctx, "test-123")
		require.NoError(t, err)
		assert.Equal(t, "test-123", result1.ID)

		result2, err := cached.GetByID(ctx, "test-456")
		require.NoError(t, err)
		assert.Equal(t, "test-456", result2.ID)

		// Each should have been fetched once
		assert.Equal(t, int32(2), mock.getByIDCallCount.Load())

		// Fetch them again - should come from cache
		result1b, err := cached.GetByID(ctx, "test-123")
		require.NoError(t, err)
		assert.Equal(t, "test-123", result1b.ID)

		result2b, err := cached.GetByID(ctx, "test-456")
		require.NoError(t, err)
		assert.Equal(t, "test-456", result2b.ID)

		// Still only 2 calls to wrapped repository
		assert.Equal(t, int32(2), mock.getByIDCallCount.Load())
	})

	t.Run("LRU eviction when cache is full", func(t *testing.T) {
		now := time.Now()

		mock := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {
				return domain.NewConversation(id, "Test "+id, "claude-sonnet-4", now), nil
			},
		}

		// Create cache with max size of 3
		cached, err := NewCachedRepository(mock, 3)
		require.NoError(t, err)

		ctx := context.Background()

		// Fill cache with 3 conversations
		_, err = cached.GetByID(ctx, "conv-1")
		require.NoError(t, err)
		_, err = cached.GetByID(ctx, "conv-2")
		require.NoError(t, err)
		_, err = cached.GetByID(ctx, "conv-3")
		require.NoError(t, err)

		assert.Equal(t, int32(3), mock.getByIDCallCount.Load())

		// Access conv-1 again to make it most recently used
		_, err = cached.GetByID(ctx, "conv-1")
		require.NoError(t, err)

		// Still 3 calls (cache hit)
		assert.Equal(t, int32(3), mock.getByIDCallCount.Load())

		// Add a 4th conversation - should evict conv-2 (least recently used)
		_, err = cached.GetByID(ctx, "conv-4")
		require.NoError(t, err)

		// Now 4 calls
		assert.Equal(t, int32(4), mock.getByIDCallCount.Load())

		// Access conv-1 again - should still be cached
		_, err = cached.GetByID(ctx, "conv-1")
		require.NoError(t, err)
		assert.Equal(t, int32(4), mock.getByIDCallCount.Load())

		// Access conv-2 - should be evicted, causing cache miss
		_, err = cached.GetByID(ctx, "conv-2")
		require.NoError(t, err)
		assert.Equal(t, int32(5), mock.getByIDCallCount.Load())
	})

	t.Run("error propagation from wrapped repository", func(t *testing.T) {
		expectedErr := &domain.AppError{Type: domain.ErrNotFound, Message: "conversation not found"}
		mock := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {
				return nil, expectedErr
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()
		result, err := cached.GetByID(ctx, "test-123")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("errors are not cached", func(t *testing.T) {
		now := time.Now()
		conv := domain.NewConversation("test-123", "Test", "claude-sonnet-4", now)

		callCount := 0
		mock := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {
				callCount++
				if callCount == 1 {
					return nil, &domain.AppError{Type: domain.ErrNotFound, Message: "not found"}
				}
				return conv, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()

		// First call returns error
		result1, err := cached.GetByID(ctx, "test-123")
		require.Error(t, err)
		assert.Nil(t, result1)

		// Second call should retry (error not cached)
		result2, err := cached.GetByID(ctx, "test-123")
		require.NoError(t, err)
		assert.NotNil(t, result2)

		// Should have called wrapped repository twice
		assert.Equal(t, int32(2), mock.getByIDCallCount.Load())
	})
}

func TestCachedRepository_Search(t *testing.T) {
	t.Run("passes through to wrapped repository", func(t *testing.T) {
		now := time.Now()
		conv := domain.NewConversation("test-123", "Test Conversation", "claude-sonnet-4", now)
		mock := &mockRepository{
			searchFunc: func(ctx context.Context, query string) ([]*domain.Conversation, error) {
				if query == "test" {
					return []*domain.Conversation{conv}, nil
				}
				return nil, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()
		result, err := cached.Search(ctx, "test")

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "test-123", result[0].ID)
		assert.Equal(t, int32(1), mock.searchCallCount.Load())
	})

	t.Run("does not cache results", func(t *testing.T) {
		now := time.Now()
		conv := domain.NewConversation("test-123", "Test Conversation", "claude-sonnet-4", now)
		mock := &mockRepository{
			searchFunc: func(ctx context.Context, query string) ([]*domain.Conversation, error) {
				return []*domain.Conversation{conv}, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()

		// First search
		result1, err := cached.Search(ctx, "test")
		require.NoError(t, err)
		require.Len(t, result1, 1)

		// Second search with same query
		result2, err := cached.Search(ctx, "test")
		require.NoError(t, err)
		require.Len(t, result2, 1)

		// Should have called wrapped repository twice (no caching)
		assert.Equal(t, int32(2), mock.searchCallCount.Load())
	})

	t.Run("error propagation from wrapped repository", func(t *testing.T) {
		expectedErr := &domain.AppError{Type: domain.ErrValidation, Message: "invalid query"}
		mock := &mockRepository{
			searchFunc: func(ctx context.Context, query string) ([]*domain.Conversation, error) {
				return nil, expectedErr
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()
		result, err := cached.Search(ctx, "test")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestCachedRepository_GetStatistics(t *testing.T) {
	t.Run("first call fetches from wrapped repository", func(t *testing.T) {
		stats := &domain.Statistics{
			ConversationCount: 5,
			MessageCount:      20,
		}
		mock := &mockRepository{
			getStatisticsFunc: func(ctx context.Context) (*domain.Statistics, error) {
				return stats, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()
		result, err := cached.GetStatistics(ctx)

		require.NoError(t, err)
		assert.Equal(t, 5, result.ConversationCount)
		assert.Equal(t, 20, result.MessageCount)
		assert.Equal(t, int32(1), mock.statisticsCallCount.Load())
	})

	t.Run("second call within TTL returns cached result", func(t *testing.T) {
		stats := &domain.Statistics{
			ConversationCount: 5,
			MessageCount:      20,
		}
		mock := &mockRepository{
			getStatisticsFunc: func(ctx context.Context) (*domain.Statistics, error) {
				return stats, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()

		// First call
		result1, err := cached.GetStatistics(ctx)
		require.NoError(t, err)
		assert.Equal(t, int32(1), mock.statisticsCallCount.Load())

		// Second call immediately (within TTL)
		result2, err := cached.GetStatistics(ctx)
		require.NoError(t, err)

		// Should still only have called wrapped repository once
		assert.Equal(t, int32(1), mock.statisticsCallCount.Load())

		// Results should be the same
		assert.Equal(t, result1.ConversationCount, result2.ConversationCount)
	})

	t.Run("call after TTL fetches fresh data", func(t *testing.T) {
		callCount := 0
		mock := &mockRepository{
			getStatisticsFunc: func(ctx context.Context) (*domain.Statistics, error) {
				callCount++
				return &domain.Statistics{
					ConversationCount: callCount * 5,
					MessageCount:      callCount * 20,
				}, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		// Override TTL to 1ms for testing
		cached.statsTTL = 1 * time.Millisecond

		ctx := context.Background()

		// First call
		result1, err := cached.GetStatistics(ctx)
		require.NoError(t, err)
		assert.Equal(t, 5, result1.ConversationCount)

		// Wait for TTL to expire
		time.Sleep(2 * time.Millisecond)

		// Second call after TTL
		result2, err := cached.GetStatistics(ctx)
		require.NoError(t, err)
		assert.Equal(t, 10, result2.ConversationCount)

		// Should have called wrapped repository twice
		assert.Equal(t, int32(2), mock.statisticsCallCount.Load())
	})

	t.Run("error propagation from wrapped repository", func(t *testing.T) {
		expectedErr := &domain.AppError{Type: domain.ErrRepository, Message: "failed to calculate stats"}
		mock := &mockRepository{
			getStatisticsFunc: func(ctx context.Context) (*domain.Statistics, error) {
				return nil, expectedErr
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()
		result, err := cached.GetStatistics(ctx)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestCachedRepository_Invalidate(t *testing.T) {
	t.Run("invalidates list cache", func(t *testing.T) {
		now := time.Now()
		conv := domain.NewConversation("test-123", "Test", "claude-sonnet-4", now)
		mock := &mockRepository{
			listFunc: func(ctx context.Context) ([]*domain.Conversation, error) {
				return []*domain.Conversation{conv}, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()

		// First call - cache it
		_, err = cached.List(ctx)
		require.NoError(t, err)
		assert.Equal(t, int32(1), mock.listCallCount.Load())

		// Second call - should be cached
		_, err = cached.List(ctx)
		require.NoError(t, err)
		assert.Equal(t, int32(1), mock.listCallCount.Load())

		// Invalidate cache
		cached.Invalidate()

		// Third call - should fetch fresh
		_, err = cached.List(ctx)
		require.NoError(t, err)
		assert.Equal(t, int32(2), mock.listCallCount.Load())
	})

	t.Run("invalidates statistics cache", func(t *testing.T) {
		stats := &domain.Statistics{ConversationCount: 5}
		mock := &mockRepository{
			getStatisticsFunc: func(ctx context.Context) (*domain.Statistics, error) {
				return stats, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()

		// First call - cache it
		_, err = cached.GetStatistics(ctx)
		require.NoError(t, err)
		assert.Equal(t, int32(1), mock.statisticsCallCount.Load())

		// Second call - should be cached
		_, err = cached.GetStatistics(ctx)
		require.NoError(t, err)
		assert.Equal(t, int32(1), mock.statisticsCallCount.Load())

		// Invalidate cache
		cached.Invalidate()

		// Third call - should fetch fresh
		_, err = cached.GetStatistics(ctx)
		require.NoError(t, err)
		assert.Equal(t, int32(2), mock.statisticsCallCount.Load())
	})

	t.Run("does not invalidate conversation cache", func(t *testing.T) {
		now := time.Now()
		conv := domain.NewConversation("test-123", "Test", "claude-sonnet-4", now)
		mock := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {
				return conv, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()

		// First call - cache it
		_, err = cached.GetByID(ctx, "test-123")
		require.NoError(t, err)
		assert.Equal(t, int32(1), mock.getByIDCallCount.Load())

		// Invalidate cache
		cached.Invalidate()

		// Second call - should still be cached (conversations are immutable)
		_, err = cached.GetByID(ctx, "test-123")
		require.NoError(t, err)
		assert.Equal(t, int32(1), mock.getByIDCallCount.Load())
	})
}

func TestCachedRepository_ConcurrentAccess(t *testing.T) {
	t.Run("concurrent GetByID calls are thread-safe", func(t *testing.T) {
		now := time.Now()
		mock := &mockRepository{
			getByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {
				// Simulate some work
				time.Sleep(1 * time.Millisecond)
				return domain.NewConversation(id, "Test "+id, "claude-sonnet-4", now), nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()
		const numGoroutines = 10
		const numCalls = 5

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < numCalls; j++ {
					id := fmt.Sprintf("conv-%d", j%3) // Use 3 different IDs
					conv, err := cached.GetByID(ctx, id)
					require.NoError(t, err)
					require.NotNil(t, conv)
					assert.Equal(t, id, conv.ID)
				}
			}(i)
		}

		wg.Wait()

		// With 3 unique IDs, should only fetch 3 times from wrapped repo
		assert.Equal(t, int32(3), mock.getByIDCallCount.Load())
	})

	t.Run("concurrent List calls are thread-safe", func(t *testing.T) {
		now := time.Now()
		conv := domain.NewConversation("test-123", "Test", "claude-sonnet-4", now)
		mock := &mockRepository{
			listFunc: func(ctx context.Context) ([]*domain.Conversation, error) {
				time.Sleep(1 * time.Millisecond)
				return []*domain.Conversation{conv}, nil
			},
		}

		cached, err := NewCachedRepository(mock, 100)
		require.NoError(t, err)

		ctx := context.Background()
		const numGoroutines = 10

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				result, err := cached.List(ctx)
				require.NoError(t, err)
				require.Len(t, result, 1)
			}()
		}

		wg.Wait()

		// Should only fetch once despite concurrent calls
		assert.Equal(t, int32(1), mock.listCallCount.Load())
	})
}

// Benchmark tests
func BenchmarkListCacheMiss(b *testing.B) {
	now := time.Now()
	conv := domain.NewConversation("test-123", "Test", "claude-sonnet-4", now)
	mock := &mockRepository{
		listFunc: func(ctx context.Context) ([]*domain.Conversation, error) {
			return []*domain.Conversation{conv}, nil
		},
	}

	cached, _ := NewCachedRepository(mock, 50)
	cached.listTTL = 1 * time.Millisecond // Force cache miss each time

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.Sleep(2 * time.Millisecond) // Expire cache
		_, _ = cached.List(ctx)
	}
}

func BenchmarkListCacheHit(b *testing.B) {
	now := time.Now()
	conv := domain.NewConversation("test-123", "Test", "claude-sonnet-4", now)
	mock := &mockRepository{
		listFunc: func(ctx context.Context) ([]*domain.Conversation, error) {
			return []*domain.Conversation{conv}, nil
		},
	}

	cached, _ := NewCachedRepository(mock, 50)
	ctx := context.Background()

	// Warm cache
	_, _ = cached.List(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cached.List(ctx)
	}
}

func BenchmarkGetByIDCacheMiss(b *testing.B) {
	now := time.Now()
	mock := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {
			return domain.NewConversation(id, "Test", "claude-sonnet-4", now), nil
		},
	}

	cached, _ := NewCachedRepository(mock, 50)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("conv-%d", i) // Different ID each time = cache miss
		_, _ = cached.GetByID(ctx, id)
	}
}

func BenchmarkGetByIDCacheHit(b *testing.B) {
	now := time.Now()
	conv := domain.NewConversation("test-123", "Test", "claude-sonnet-4", now)
	mock := &mockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {
			return conv, nil
		},
	}

	cached, _ := NewCachedRepository(mock, 50)
	ctx := context.Background()

	// Warm cache
	_, _ = cached.GetByID(ctx, "test-123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cached.GetByID(ctx, "test-123")
	}
}

func BenchmarkGetStatisticsCacheHit(b *testing.B) {
	stats := &domain.Statistics{ConversationCount: 5, MessageCount: 20}
	mock := &mockRepository{
		getStatisticsFunc: func(ctx context.Context) (*domain.Statistics, error) {
			return stats, nil
		},
	}

	cached, _ := NewCachedRepository(mock, 50)
	ctx := context.Background()

	// Warm cache
	_, _ = cached.GetStatistics(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cached.GetStatistics(ctx)
	}
}
