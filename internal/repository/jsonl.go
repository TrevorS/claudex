// ABOUTME: JSONLRepository implements Repository interface for conversations stored as JSONL files
package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/TrevorS/claudex/internal/domain"
)

// JSONLRepository implements Repository by reading conversations from JSONL files
// stored in the ~/.claude/projects/ directory structure.
type JSONLRepository struct {
	rootPath string
	mu       sync.RWMutex
	metadata map[string]*conversationMetadata // sessionID -> metadata
	loaded   bool
}

// conversationMetadata stores lightweight conversation information for lazy loading
type conversationMetadata struct {
	ID           string
	Title        string
	ProjectPath  string
	FilePath     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	MessageCount int
}

// conversationLine represents a single line in a Claude JSONL conversation file
type conversationLine struct {
	Type      string   `json:"type"`
	Summary   string   `json:"summary"`
	Message   *msgData `json:"message"`
	SessionID string   `json:"sessionId"`
	Timestamp string   `json:"timestamp"`
	UUID      string   `json:"uuid"`
	Version   string   `json:"version"`
}

// msgData represents message content and usage information
type msgData struct {
	Role    string     `json:"role"`
	Content string     `json:"content"`
	Usage   *usageData `json:"usage"`
}

// usageData represents token usage for a message
type usageData struct {
	InputTokens         int `json:"input_tokens"`
	OutputTokens        int `json:"output_tokens"`
	CacheReadTokens     int `json:"cache_read_tokens"`
	CacheCreationTokens int `json:"cache_creation_tokens"`
}

// NewJSONLRepository creates a new repository for the given root path.
// The root path should be the directory containing project subdirectories (e.g., ~/.claude/projects/)
func NewJSONLRepository(rootPath string) (*JSONLRepository, error) {
	// Verify the root path exists
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access repository root path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("repository root path is not a directory: %s", rootPath)
	}

	return &JSONLRepository{
		rootPath: rootPath,
		metadata: make(map[string]*conversationMetadata),
	}, nil
}

// List returns all conversations with metadata only (no message bodies).
// This scans the directory structure and extracts metadata from first/last lines of each JSONL file.
func (r *JSONLRepository) List(ctx context.Context) ([]*domain.Conversation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Load metadata if not already loaded
	if !r.loaded {
		if err := r.discoverProjects(); err != nil {
			return nil, fmt.Errorf("failed to discover projects: %w", err)
		}
		r.loaded = true
	}

	// Build conversations from metadata
	conversations := make([]*domain.Conversation, 0, len(r.metadata))
	for _, meta := range r.metadata {
		// Cache message count from metadata so List() can show accurate counts
		messageCount := meta.MessageCount
		zero := int64(0) // Placeholder for token count (not available from metadata alone)
		conv := &domain.Conversation{
			ID:                 meta.ID,
			Title:              meta.Title,
			Model:              "unknown", // Model info not stored in metadata, will be updated on full load
			ProjectPath:        meta.ProjectPath,
			CreatedAt:          meta.CreatedAt,
			UpdatedAt:          meta.UpdatedAt,
			Messages:           []*domain.Message{}, // Empty for metadata-only load
			CachedMessageCount: &messageCount,       // Use metadata message count
			CachedTotalTokens:  &zero,               // Tokens not available from metadata
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// GetByID loads a complete conversation including all messages.
func (r *JSONLRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {
	r.mu.RLock()
	meta, exists := r.metadata[id]
	r.mu.RUnlock()

	if !exists {
		// If not in metadata, try to discover first
		r.mu.Lock()
		if !r.loaded {
			if err := r.discoverProjects(); err != nil {
				r.mu.Unlock()
				return nil, fmt.Errorf("failed to discover projects: %w", err)
			}
			r.loaded = true
		}
		meta, exists = r.metadata[id]
		r.mu.Unlock()

		if !exists {
			return nil, nil // Conversation not found
		}
	}

	// Load the full conversation from the JSONL file
	return r.loadConversation(meta)
}

// Search filters conversations by substring matching on title.
func (r *JSONLRepository) Search(ctx context.Context, query string) ([]*domain.Conversation, error) {
	conversations, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var results []*domain.Conversation
	lowerQuery := strings.ToLower(query)

	for _, conv := range conversations {
		if strings.Contains(strings.ToLower(conv.Title), lowerQuery) {
			results = append(results, conv)
		}
	}

	return results, nil
}

// GetStatistics returns aggregated statistics across all conversations.
func (r *JSONLRepository) GetStatistics(ctx context.Context) (*domain.Statistics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Ensure metadata is loaded
	if !r.loaded {
		// Create temp copy to load metadata
		r.mu.RUnlock()
		_, err := r.List(ctx)
		r.mu.RLock()
		if err != nil {
			return nil, err
		}
	}

	stats := domain.NewStatistics()

	// Sum up message counts from metadata
	for _, meta := range r.metadata {
		stats.ConversationCount++
		stats.MessageCount += meta.MessageCount
		// Note: Counts from metadata only, not from full conversation load
		// For accurate token counts, would need to load all conversations fully
	}

	return stats, nil
}

// ============ Private Methods ============

// discoverProjects scans the root directory and extracts metadata from all JSONL files.
func (r *JSONLRepository) discoverProjects() error {
	entries, err := os.ReadDir(r.rootPath)
	if err != nil {
		return fmt.Errorf("failed to read root directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectName := entry.Name()
		projectPath := filepath.Join(r.rootPath, projectName)
		decodedPath := decodePath(projectName)

		// Read all JSONL files in this project directory
		if err := r.discoverConversations(projectPath, decodedPath); err != nil {
			// Log error but continue discovering other projects
			continue
		}
	}

	return nil
}

// discoverConversations scans a project directory and extracts metadata from JSONL files.
func (r *JSONLRepository) discoverConversations(projectPath, decodedPath string) error {
	entries, err := os.ReadDir(projectPath)
	if err != nil {
		return fmt.Errorf("failed to read project directory %s: %w", projectPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		filePath := filepath.Join(projectPath, entry.Name())

		// Extract metadata from first and last lines
		meta, err := r.extractMetadata(filePath, decodedPath)
		if err != nil {
			// Skip files with errors
			continue
		}

		if meta != nil {
			r.metadata[meta.ID] = meta
		}
	}

	return nil
}

// extractMetadata reads the first and last lines of a JSONL file to extract metadata.
// Returns nil metadata if the file is empty.
func (r *JSONLRepository) extractMetadata(filePath, projectPath string) (*conversationMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read first line for session ID and title
	var firstLine, lastLine conversationLine
	var hasLines bool

	if scanner.Scan() {
		hasLines = true
		text := scanner.Text()
		if err := json.Unmarshal([]byte(text), &firstLine); err != nil {
			return nil, fmt.Errorf("failed to parse first line of %s: %w", filePath, err)
		}
	}

	// Read all lines to find the last one
	for scanner.Scan() {
		text := scanner.Text()
		if err := json.Unmarshal([]byte(text), &lastLine); err != nil {
			// Skip malformed lines
			continue
		}
	}

	if !hasLines {
		return nil, nil // Empty file
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	// Extract metadata
	sessionID := firstLine.SessionID
	if sessionID == "" {
		// Use filename stem as fallback ID
		sessionID = strings.TrimSuffix(filepath.Base(filePath), ".jsonl")
	}

	createdAt, _ := time.Parse(time.RFC3339Nano, firstLine.Timestamp)
	updatedAt, _ := time.Parse(time.RFC3339Nano, lastLine.Timestamp)
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}

	title := firstLine.Summary
	if title == "" {
		title = "Untitled Conversation"
	}

	// Count message lines (skip summary lines)
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	messageCount := 0
	for scanner.Scan() {
		text := scanner.Text()
		var line conversationLine
		if err := json.Unmarshal([]byte(text), &line); err != nil {
			continue
		}
		if line.Type == "message" {
			messageCount++
		}
	}

	return &conversationMetadata{
		ID:           sessionID,
		Title:        title,
		ProjectPath:  projectPath,
		FilePath:     filePath,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		MessageCount: messageCount,
	}, nil
}

// loadConversation fully loads a conversation including all messages.
func (r *JSONLRepository) loadConversation(meta *conversationMetadata) (*domain.Conversation, error) {
	file, err := os.Open(meta.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open conversation file %s: %w", meta.FilePath, err)
	}
	defer file.Close()

	conv := &domain.Conversation{
		ID:          meta.ID,
		Title:       meta.Title,
		Model:       "unknown", // Model info not available from file metadata
		ProjectPath: meta.ProjectPath,
		CreatedAt:   meta.CreatedAt,
		UpdatedAt:   meta.UpdatedAt,
		Messages:    make([]*domain.Message, 0),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var cl conversationLine
		if err := json.Unmarshal([]byte(line), &cl); err != nil {
			// Skip malformed lines, continue parsing
			continue
		}

		// Only add message-type lines (skip summaries)
		if cl.Type == "message" && cl.Message != nil {
			var role domain.Role
			switch cl.Message.Role {
			case "user":
				role = domain.RoleUser
			case "assistant":
				role = domain.RoleAssistant
			case "tool":
				role = domain.RoleTool
			default:
				continue // Skip unknown roles
			}

			timestamp, _ := time.Parse(time.RFC3339Nano, cl.Timestamp)

			// Calculate token count from usage
			var tokenCount *domain.TokenCount
			if cl.Message.Usage != nil {
				tokenCount = &domain.TokenCount{
					Input:       int64(cl.Message.Usage.InputTokens),
					Output:      int64(cl.Message.Usage.OutputTokens),
					CacheRead:   int64(cl.Message.Usage.CacheReadTokens),
					CacheCreate: int64(cl.Message.Usage.CacheCreationTokens),
				}
			}

			msg := &domain.Message{
				Role:      role,
				Content:   cl.Message.Content,
				Timestamp: timestamp,
				Tokens:    tokenCount,
			}

			conv.Messages = append(conv.Messages, msg)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading conversation file %s: %w", meta.FilePath, err)
	}

	return conv, nil
}

// decodePath converts hyphenated directory names back to real paths.
// Example: "-Users-trevor-Projects-claudex" -> "/Users/trevor/Projects/claudex"
// Special case: double hyphen "--" becomes a single "-" (for paths with dots like ".claude")
func decodePath(encoded string) string {
	if len(encoded) == 0 {
		return ""
	}

	// Remove leading hyphen if present
	if encoded[0] == '-' {
		encoded = encoded[1:]
	}

	// Replace double hyphens with a placeholder to preserve them as single hyphens
	encoded = strings.ReplaceAll(encoded, "--", "\x00")

	// Replace remaining hyphens with slashes
	decoded := strings.ReplaceAll(encoded, "-", "/")

	// Replace placeholder with single hyphen
	decoded = strings.ReplaceAll(decoded, "\x00", "-")

	// Ensure it starts with /
	if !strings.HasPrefix(decoded, "/") {
		decoded = "/" + decoded
	}

	return decoded
}
