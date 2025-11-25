// ABOUTME: filter.go implements the Filter interface that matches AST nodes against
// conversations, providing search functionality for terms, fields, and boolean operations.
package search

import (
	"strconv"
	"strings"

	"github.com/TrevorS/claudex/internal/domain"
)

// Filter defines the interface for matching conversations against search criteria.
type Filter interface {
	Match(c *domain.Conversation) bool
}

// NodeFilter wraps an AST node and implements the Filter interface.
type NodeFilter struct {
	node Node
}

// NewNodeFilter creates a filter from an AST node.
func NewNodeFilter(node Node) *NodeFilter {
	return &NodeFilter{node: node}
}

// Match evaluates the filter against a conversation.
func (f *NodeFilter) Match(c *domain.Conversation) bool {
	return f.matchNode(f.node, c)
}

// matchNode recursively evaluates an AST node against a conversation.
func (f *NodeFilter) matchNode(node Node, c *domain.Conversation) bool {
	switch n := node.(type) {
	case *TermNode:
		return f.matchTerm(n.Term, c)
	case *FieldNode:
		return f.matchField(n.Field, n.Value, n.Op, c)
	case *AndNode:
		return f.matchNode(n.Left, c) && f.matchNode(n.Right, c)
	case *OrNode:
		return f.matchNode(n.Left, c) || f.matchNode(n.Right, c)
	case *NotNode:
		return !f.matchNode(n.Node, c)
	default:
		return false
	}
}

// matchTerm searches for a term in all searchable fields (title, content, model).
func (f *NodeFilter) matchTerm(term string, c *domain.Conversation) bool {
	if term == "" {
		return false
	}

	termLower := strings.ToLower(term)

	// Search in title
	if strings.Contains(strings.ToLower(c.Title), termLower) {
		return true
	}

	// Search in model
	if strings.Contains(strings.ToLower(c.Model), termLower) {
		return true
	}

	// Search in all message content
	for _, msg := range c.Messages {
		if strings.Contains(strings.ToLower(msg.Content), termLower) {
			return true
		}
	}

	return false
}

// matchField matches a specific field with a value and operator.
func (f *NodeFilter) matchField(field, value, op string, c *domain.Conversation) bool {
	fieldLower := strings.ToLower(field)
	valueLower := strings.ToLower(value)

	switch fieldLower {
	case "title":
		return f.compareString(c.Title, valueLower, op)

	case "content":
		// Search across all message content
		for _, msg := range c.Messages {
			if f.compareString(msg.Content, valueLower, op) {
				return true
			}
		}
		return false

	case "model":
		return f.compareString(c.Model, valueLower, op)

	case "tokens":
		total := c.TotalTokens()
		targetVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return false
		}
		return f.compareNumeric(total, targetVal, op)

	case "messages":
		count := int64(c.MessageCount())
		targetVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return false
		}
		return f.compareNumeric(count, targetVal, op)

	case "date", "created":
		// TODO: Implement date comparison in future enhancement
		return false

	default:
		return false
	}
}

// compareString performs case-insensitive substring matching for strings.
func (f *NodeFilter) compareString(haystack, needle, op string) bool {
	haystackLower := strings.ToLower(haystack)
	needleLower := strings.ToLower(needle)

	switch op {
	case "=":
		return strings.Contains(haystackLower, needleLower)
	default:
		// For strings, only equality/contains makes sense
		return strings.Contains(haystackLower, needleLower)
	}
}

// compareNumeric performs numeric comparisons.
func (f *NodeFilter) compareNumeric(actual, target int64, op string) bool {
	switch op {
	case "=":
		return actual == target
	case ">":
		return actual > target
	case "<":
		return actual < target
	case ">=":
		return actual >= target
	case "<=":
		return actual <= target
	default:
		return false
	}
}
