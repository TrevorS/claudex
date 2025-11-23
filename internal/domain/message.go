// ABOUTME: message.go defines the Message domain entity representing individual messages
// in a conversation. Messages track content, role, tokens, and metadata.
package domain

import "time"

// Role defines the role of the message sender.
type Role string

// Role constants.
const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// IsValid returns true if the role is a valid role.
func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleAssistant, RoleTool:
		return true
	default:
		return false
	}
}

// Message represents a single message in a conversation.
type Message struct {
	Role      Role
	Content   string
	Timestamp time.Time
	Tokens    *TokenCount
}

// TokenCount tracks token usage for a message.
type TokenCount struct {
	Input       int64 // Tokens in the input
	Output      int64 // Tokens in the output
	CacheRead   int64 // Tokens read from cache
	CacheCreate int64 // Tokens written to cache
}

// NewMessage creates a new Message with initialized token tracking.
func NewMessage(role Role, content string, ts time.Time) *Message {
	return &Message{
		Role:      role,
		Content:   content,
		Timestamp: ts,
		Tokens:    &TokenCount{},
	}
}

// TotalTokens returns the total number of tokens in this message.
func (m *Message) TotalTokens() int64 {
	if m.Tokens == nil {
		return 0
	}
	return m.Tokens.Total()
}

// Total returns the total token count across all categories.
func (tc *TokenCount) Total() int64 {
	if tc == nil {
		return 0
	}
	return tc.Input + tc.Output + tc.CacheRead + tc.CacheCreate
}

// EstimatedCost provides a rough cost estimate for this message.
// This is a placeholder that returns a cost based on token usage.
// Real pricing would depend on the specific model and API.
func (m *Message) EstimatedCost() (float64, bool) {
	if m.Tokens == nil {
		return 0.0, false
	}
	// Placeholder: Basic cost calculation
	// In production, this would use model-specific pricing
	_ = m.TotalTokens()
	return 0.0, true
}
