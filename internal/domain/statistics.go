// ABOUTME: statistics.go defines the Statistics domain entity for aggregating
// metrics across conversations including tokens, costs, and message distribution.
package domain

// Statistics aggregates metrics across conversations.
type Statistics struct {
	ConversationCount int
	MessageCount      int
	TotalTokens       int64
	TotalCost         float64
	ModelDistribution map[string]int
	RoleDistribution  map[string]int
}

// NewStatistics creates a new Statistics with initialized maps.
func NewStatistics() *Statistics {
	return &Statistics{
		ConversationCount: 0,
		MessageCount:      0,
		TotalTokens:       0,
		TotalCost:         0.0,
		ModelDistribution: make(map[string]int),
		RoleDistribution:  make(map[string]int),
	}
}

// AddConversation adds metrics from a conversation to the statistics.
func (s *Statistics) AddConversation(conv *Conversation) {
	s.ConversationCount++
	s.MessageCount += conv.MessageCount()
	s.TotalTokens += conv.TotalTokens()

	// Track model distribution
	s.ModelDistribution[conv.Model]++

	// Track role distribution
	for _, msg := range conv.Messages {
		s.RoleDistribution[string(msg.Role)]++
	}
}

// AverageMessagesPerConversation returns the average number of messages
// per conversation, or 0.0 if no conversations exist.
func (s *Statistics) AverageMessagesPerConversation() float64 {
	if s.ConversationCount == 0 {
		return 0.0
	}
	return float64(s.MessageCount) / float64(s.ConversationCount)
}

// AverageTokensPerMessage returns the average tokens per message.
// It returns (0.0, false) if no messages exist.
func (s *Statistics) AverageTokensPerMessage() (float64, bool) {
	if s.MessageCount == 0 {
		return 0.0, false
	}
	return float64(s.TotalTokens) / float64(s.MessageCount), true
}
