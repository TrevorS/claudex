// ABOUTME: filter_test.go contains comprehensive tests for search filters,
// including text matching, field filtering, and boolean operations.
package search

import (
	"strings"
	"testing"
	"time"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Helper functions to create test conversations
func makeConversation(title, model string, messageContents []string, tokens int64) *domain.Conversation {
	conv := domain.NewConversation("test-id", title, model, time.Now())
	for _, content := range messageContents {
		msg := domain.NewMessage(domain.RoleUser, content, time.Now())
		msg.Tokens = &domain.TokenCount{Input: tokens / int64(len(messageContents))}
		conv.AddMessage(msg)
	}
	return conv
}

func TestNodeFilter_TermNode(t *testing.T) {
	tests := []struct {
		name    string
		term    string
		conv    *domain.Conversation
		matches bool
	}{
		{
			name:    "matches in title",
			term:    "authentication",
			conv:    makeConversation("Authentication Flow", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "matches in title case-insensitive",
			term:    "AUTH",
			conv:    makeConversation("Authentication Flow", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "matches in message content",
			term:    "python",
			conv:    makeConversation("My Chat", "sonnet", []string{"I love python programming"}, 100),
			matches: true,
		},
		{
			name:    "matches in model",
			term:    "sonnet",
			conv:    makeConversation("My Chat", "claude-sonnet-4", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "does not match",
			term:    "javascript",
			conv:    makeConversation("Python Chat", "sonnet", []string{"I love python"}, 100),
			matches: false,
		},
		{
			name:    "partial match in word",
			term:    "auth",
			conv:    makeConversation("Authentication", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewNodeFilter(&TermNode{Term: tt.term})
			result := filter.Match(tt.conv)
			assert.Equal(t, tt.matches, result)
		})
	}
}

func TestNodeFilter_FieldNode(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   string
		op      string
		conv    *domain.Conversation
		matches bool
	}{
		{
			name:    "title equals",
			field:   "title",
			value:   "My Chat",
			op:      "=",
			conv:    makeConversation("My Chat", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "title partial match",
			field:   "title",
			value:   "Chat",
			op:      "=",
			conv:    makeConversation("My Chat Room", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "title case-insensitive",
			field:   "title",
			value:   "chat",
			op:      "=",
			conv:    makeConversation("My Chat", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "content matches",
			field:   "content",
			value:   "python",
			op:      "=",
			conv:    makeConversation("Chat", "sonnet", []string{"I love python"}, 100),
			matches: true,
		},
		{
			name:    "model matches",
			field:   "model",
			value:   "sonnet",
			op:      "=",
			conv:    makeConversation("Chat", "claude-sonnet-4", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "tokens greater than (true)",
			field:   "tokens",
			value:   "50",
			op:      ">",
			conv:    makeConversation("Chat", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "tokens greater than (false)",
			field:   "tokens",
			value:   "200",
			op:      ">",
			conv:    makeConversation("Chat", "sonnet", []string{"hello"}, 100),
			matches: false,
		},
		{
			name:    "tokens less than (true)",
			field:   "tokens",
			value:   "200",
			op:      "<",
			conv:    makeConversation("Chat", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "tokens greater or equal",
			field:   "tokens",
			value:   "100",
			op:      ">=",
			conv:    makeConversation("Chat", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "tokens less or equal",
			field:   "tokens",
			value:   "100",
			op:      "<=",
			conv:    makeConversation("Chat", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
		{
			name:    "messages count",
			field:   "messages",
			value:   "2",
			op:      ">=",
			conv:    makeConversation("Chat", "sonnet", []string{"msg1", "msg2", "msg3"}, 100),
			matches: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewNodeFilter(&FieldNode{
				Field: tt.field,
				Value: tt.value,
				Op:    tt.op,
			})
			result := filter.Match(tt.conv)
			assert.Equal(t, tt.matches, result)
		})
	}
}

func TestNodeFilter_AndNode(t *testing.T) {
	tests := []struct {
		name    string
		left    Node
		right   Node
		conv    *domain.Conversation
		matches bool
	}{
		{
			name:    "both match",
			left:    &TermNode{Term: "python"},
			right:   &TermNode{Term: "rust"},
			conv:    makeConversation("Chat", "sonnet", []string{"python and rust"}, 100),
			matches: true,
		},
		{
			name:    "left matches, right does not",
			left:    &TermNode{Term: "python"},
			right:   &TermNode{Term: "javascript"},
			conv:    makeConversation("Chat", "sonnet", []string{"python programming"}, 100),
			matches: false,
		},
		{
			name:    "neither match",
			left:    &TermNode{Term: "java"},
			right:   &TermNode{Term: "javascript"},
			conv:    makeConversation("Chat", "sonnet", []string{"python programming"}, 100),
			matches: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewNodeFilter(&AndNode{Left: tt.left, Right: tt.right})
			result := filter.Match(tt.conv)
			assert.Equal(t, tt.matches, result)
		})
	}
}

func TestNodeFilter_OrNode(t *testing.T) {
	tests := []struct {
		name    string
		left    Node
		right   Node
		conv    *domain.Conversation
		matches bool
	}{
		{
			name:    "both match",
			left:    &TermNode{Term: "python"},
			right:   &TermNode{Term: "rust"},
			conv:    makeConversation("Chat", "sonnet", []string{"python and rust"}, 100),
			matches: true,
		},
		{
			name:    "left matches",
			left:    &TermNode{Term: "python"},
			right:   &TermNode{Term: "javascript"},
			conv:    makeConversation("Chat", "sonnet", []string{"python programming"}, 100),
			matches: true,
		},
		{
			name:    "right matches",
			left:    &TermNode{Term: "javascript"},
			right:   &TermNode{Term: "python"},
			conv:    makeConversation("Chat", "sonnet", []string{"python programming"}, 100),
			matches: true,
		},
		{
			name:    "neither match",
			left:    &TermNode{Term: "java"},
			right:   &TermNode{Term: "javascript"},
			conv:    makeConversation("Chat", "sonnet", []string{"python programming"}, 100),
			matches: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewNodeFilter(&OrNode{Left: tt.left, Right: tt.right})
			result := filter.Match(tt.conv)
			assert.Equal(t, tt.matches, result)
		})
	}
}

func TestNodeFilter_NotNode(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		conv    *domain.Conversation
		matches bool
	}{
		{
			name:    "negates match",
			node:    &TermNode{Term: "python"},
			conv:    makeConversation("Chat", "sonnet", []string{"python programming"}, 100),
			matches: false,
		},
		{
			name:    "negates non-match",
			node:    &TermNode{Term: "javascript"},
			conv:    makeConversation("Chat", "sonnet", []string{"python programming"}, 100),
			matches: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewNodeFilter(&NotNode{Node: tt.node})
			result := filter.Match(tt.conv)
			assert.Equal(t, tt.matches, result)
		})
	}
}

func TestNodeFilter_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		conv    *domain.Conversation
		matches bool
	}{
		{
			name: "(python OR rust) AND sonnet",
			node: &AndNode{
				Left: &OrNode{
					Left:  &TermNode{Term: "python"},
					Right: &TermNode{Term: "rust"},
				},
				Right: &FieldNode{Field: "model", Value: "sonnet", Op: "="},
			},
			conv:    makeConversation("Chat", "claude-sonnet-4", []string{"python code"}, 100),
			matches: true,
		},
		{
			name: "title:chat AND NOT archived",
			node: &AndNode{
				Left: &FieldNode{Field: "title", Value: "chat", Op: "="},
				Right: &NotNode{
					Node: &TermNode{Term: "archived"},
				},
			},
			conv:    makeConversation("My Chat", "sonnet", []string{"hello"}, 100),
			matches: true,
		},
		{
			name: "tokens:>5000 AND (python OR rust)",
			node: &AndNode{
				Left: &FieldNode{Field: "tokens", Value: "5000", Op: ">"},
				Right: &OrNode{
					Left:  &TermNode{Term: "python"},
					Right: &TermNode{Term: "rust"},
				},
			},
			conv:    makeConversation("Chat", "sonnet", []string{strings.Repeat("python ", 1000)}, 10000),
			matches: true,
		},
		{
			name: "tokens:>5000 AND (python OR rust) - no match",
			node: &AndNode{
				Left: &FieldNode{Field: "tokens", Value: "5000", Op: ">"},
				Right: &OrNode{
					Left:  &TermNode{Term: "python"},
					Right: &TermNode{Term: "rust"},
				},
			},
			conv:    makeConversation("Chat", "sonnet", []string{"javascript"}, 100),
			matches: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewNodeFilter(tt.node)
			result := filter.Match(tt.conv)
			assert.Equal(t, tt.matches, result)
		})
	}
}

func TestNodeFilter_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		conv    *domain.Conversation
		matches bool
	}{
		{
			name:    "empty term matches nothing",
			node:    &TermNode{Term: ""},
			conv:    makeConversation("Chat", "sonnet", []string{"hello"}, 100),
			matches: false,
		},
		{
			name:    "conversation with no messages",
			node:    &TermNode{Term: "python"},
			conv:    makeConversation("Chat", "sonnet", []string{}, 0),
			matches: false,
		},
		{
			name:    "unknown field returns false",
			node:    &FieldNode{Field: "unknown", Value: "test", Op: "="},
			conv:    makeConversation("Chat", "sonnet", []string{"hello"}, 100),
			matches: false,
		},
		{
			name:    "invalid numeric comparison returns false",
			node:    &FieldNode{Field: "tokens", Value: "notanumber", Op: ">"},
			conv:    makeConversation("Chat", "sonnet", []string{"hello"}, 100),
			matches: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewNodeFilter(tt.node)
			result := filter.Match(tt.conv)
			assert.Equal(t, tt.matches, result)
		})
	}
}
