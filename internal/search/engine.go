// ABOUTME: engine.go implements the search engine that orchestrates query parsing,
// filtering, and relevance ranking to provide search results.
package search

import (
	"context"
	"sort"
	"strings"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/TrevorS/claudex/internal/repository"
)

// Engine orchestrates search operations.
type Engine struct {
	repository repository.Repository
}

// NewEngine creates a new search engine.
func NewEngine(repo repository.Repository) *Engine {
	return &Engine{repository: repo}
}

// SearchResult represents a search result with relevance score.
type SearchResult struct {
	Conversation *domain.Conversation
	Score        float64
}

// Search executes a search query and returns ranked results.
func (e *Engine) Search(ctx context.Context, query string) ([]*domain.Conversation, error) {
	// Parse query
	parser := NewParser(query)
	node, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	// Get all conversations
	allConvs, err := e.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	// Apply filter and rank
	filter := NewNodeFilter(node)
	var results []SearchResult

	for _, conv := range allConvs {
		if filter.Match(conv) {
			score := e.rankResult(conv, query)
			results = append(results, SearchResult{
				Conversation: conv,
				Score:        score,
			})
		}
	}

	// Sort by score (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Limit to top 100
	if len(results) > 100 {
		results = results[:100]
	}

	// Extract conversations
	conversations := make([]*domain.Conversation, len(results))
	for i, r := range results {
		conversations[i] = r.Conversation
	}

	return conversations, nil
}

// rankResult calculates relevance score for a conversation.
// Title matches: 3x weight
// Field matches (model): 2x weight
// Content matches: 1x weight
func (e *Engine) rankResult(c *domain.Conversation, query string) float64 {
	score := 0.0

	// Extract search terms from query (simple approach: lowercase words)
	terms := extractTerms(query)

	for _, term := range terms {
		termLower := strings.ToLower(term)

		// Title matches (3x weight)
		titleMatches := strings.Count(strings.ToLower(c.Title), termLower)
		score += float64(titleMatches) * 3.0

		// Model matches (2x weight)
		modelMatches := strings.Count(strings.ToLower(c.Model), termLower)
		score += float64(modelMatches) * 2.0

		// Content matches (1x weight)
		for _, msg := range c.Messages {
			contentMatches := strings.Count(strings.ToLower(msg.Content), termLower)
			score += float64(contentMatches) * 1.0
		}
	}

	// Base score of 1.0 for any match
	if score == 0 {
		score = 1.0
	}

	return score
}

// extractTerms extracts search terms from a query string.
// This is a simplified approach that extracts words, ignoring operators and syntax.
func extractTerms(query string) []string {
	var terms []string
	words := strings.Fields(query)

	for _, word := range words {
		// Skip operators and special characters
		word = strings.Trim(word, "()\"")
		if word == "" {
			continue
		}
		upper := strings.ToUpper(word)
		if upper == "AND" || upper == "OR" || upper == "NOT" {
			continue
		}
		// Skip field names (contains colon)
		if strings.Contains(word, ":") {
			parts := strings.Split(word, ":")
			// Include the value part after colon
			if len(parts) > 1 && parts[1] != "" {
				// Skip operators like >5000
				val := strings.TrimLeft(parts[1], "><=>")
				if val != "" {
					terms = append(terms, val)
				}
			}
			continue
		}
		terms = append(terms, word)
	}

	return terms
}
