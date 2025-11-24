// ABOUTME: conversation.go defines the Conversation domain entity representing a
// complete conversation with metadata and a collection of messages.
package domain

import (
	"time"
)

// Conversation represents a conversation with Claude.
type Conversation struct {
	ID                 string
	Title              string
	Model              string
	ProjectPath        string // Project path (e.g., "foo/bar" decoded from filesystem "-foo-bar")
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Messages           []*Message
	CachedMessageCount *int   // Optional cached message count from metadata (for lazy-loaded conversations)
	CachedTotalTokens  *int64 // Optional cached total tokens from metadata
}

// NewConversation creates a new Conversation with initialized message slice.
func NewConversation(id, title, model string, ts time.Time) *Conversation {
	return &Conversation{
		ID:        id,
		Title:     title,
		Model:     model,
		CreatedAt: ts,
		UpdatedAt: ts,
		Messages:  make([]*Message, 0),
	}
}

// AddMessage adds a message to the conversation and updates UpdatedAt.
func (c *Conversation) AddMessage(msg *Message) {
	c.Messages = append(c.Messages, msg)
	c.UpdatedAt = time.Now()
}

// MessageCount returns the total number of messages in the conversation.
// Returns cached count if available (for metadata-only conversations), otherwise counts loaded messages.
func (c *Conversation) MessageCount() int {
	if c.CachedMessageCount != nil {
		return *c.CachedMessageCount
	}
	return len(c.Messages)
}

// TotalTokens returns the sum of all tokens across all messages.
// Returns cached value if available (for metadata-only conversations), otherwise sums loaded messages.
func (c *Conversation) TotalTokens() int64 {
	if c.CachedTotalTokens != nil {
		return *c.CachedTotalTokens
	}
	total := int64(0)
	for _, msg := range c.Messages {
		total += msg.TotalTokens()
	}
	return total
}

// UserMessageCount returns the count of user messages only.
func (c *Conversation) UserMessageCount() int {
	count := 0
	for _, msg := range c.Messages {
		if msg.Role == RoleUser {
			count++
		}
	}
	return count
}

// LastMessage returns the most recent message or nil if no messages exist.
func (c *Conversation) LastMessage() *Message {
	if len(c.Messages) == 0 {
		return nil
	}
	return c.Messages[len(c.Messages)-1]
}

// Validate checks that required fields are set.
func (c *Conversation) Validate() error {
	if c.ID == "" {
		return NewAppError(ErrValidation, "conversation ID cannot be empty", nil)
	}
	if c.Title == "" {
		return NewAppError(ErrValidation, "conversation title cannot be empty", nil)
	}
	if c.Model == "" {
		return NewAppError(ErrValidation, "conversation model cannot be empty", nil)
	}
	return nil
}
