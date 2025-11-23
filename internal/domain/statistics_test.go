package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var mockNow = time.Date(2025, 11, 23, 12, 0, 0, 0, time.UTC)

func TestNewStatistics(t *testing.T) {
	stats := NewStatistics()

	assert.NotNil(t, stats)
	assert.Equal(t, int64(0), stats.TotalTokens)
	assert.Equal(t, 0, stats.ConversationCount)
	assert.Equal(t, 0, stats.MessageCount)
	assert.Equal(t, 0.0, stats.TotalCost)
}

func TestStatistics_AddConversation(t *testing.T) {
	stats := NewStatistics()
	conv := NewConversation("conv-1", "Test", "claude-3-sonnet", mockNow)

	msg := NewMessage(RoleUser, "Hello", mockNow)
	msg.Tokens = &TokenCount{Input: 100, Output: 50}
	conv.AddMessage(msg)

	stats.AddConversation(conv)

	assert.Equal(t, 1, stats.ConversationCount)
	assert.Equal(t, 1, stats.MessageCount)
	assert.Equal(t, int64(150), stats.TotalTokens)
}

func TestStatistics_AddMultipleConversations(t *testing.T) {
	stats := NewStatistics()

	for i := 1; i <= 3; i++ {
		conv := NewConversation("conv-"+string(rune('0'+i)), "Test", "claude-3-sonnet", mockNow)
		for j := 0; j < 2; j++ {
			msg := NewMessage(RoleUser, "Msg", mockNow)
			msg.Tokens = &TokenCount{Input: 100}
			conv.AddMessage(msg)
		}
		stats.AddConversation(conv)
	}

	assert.Equal(t, 3, stats.ConversationCount)
	assert.Equal(t, 6, stats.MessageCount)
	assert.Equal(t, int64(600), stats.TotalTokens)
}

func TestStatistics_AverageMessagesPerConversation(t *testing.T) {
	tests := []struct {
		name     string
		convs    int
		msgsEach int
		want     float64
	}{
		{
			name:     "zero conversations",
			convs:    0,
			msgsEach: 0,
			want:     0.0,
		},
		{
			name:     "one conversation two messages",
			convs:    1,
			msgsEach: 2,
			want:     2.0,
		},
		{
			name:     "three conversations two messages each",
			convs:    3,
			msgsEach: 2,
			want:     2.0,
		},
		{
			name:     "three conversations three messages each",
			convs:    3,
			msgsEach: 3,
			want:     3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := NewStatistics()
			for i := 0; i < tt.convs; i++ {
				conv := NewConversation("conv", "Test", "claude-3-sonnet", mockNow)
				for j := 0; j < tt.msgsEach; j++ {
					conv.AddMessage(NewMessage(RoleUser, "msg", mockNow))
				}
				stats.AddConversation(conv)
			}

			got := stats.AverageMessagesPerConversation()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStatistics_AverageTokensPerMessage(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*Statistics)
		wantResult float64
		wantOK     bool
	}{
		{
			name: "zero messages",
			setup: func(s *Statistics) {
				// Don't add anything
			},
			wantResult: 0.0,
			wantOK:     false,
		},
		{
			name: "single message",
			setup: func(s *Statistics) {
				conv := NewConversation("conv-1", "Test", "claude-3-sonnet", mockNow)
				msg := NewMessage(RoleUser, "Hello", mockNow)
				msg.Tokens = &TokenCount{Input: 100}
				conv.AddMessage(msg)
				s.AddConversation(conv)
			},
			wantResult: 100.0,
			wantOK:     true,
		},
		{
			name: "multiple messages",
			setup: func(s *Statistics) {
				conv := NewConversation("conv-1", "Test", "claude-3-sonnet", mockNow)
				msg1 := NewMessage(RoleUser, "Hello", mockNow)
				msg1.Tokens = &TokenCount{Input: 100}
				msg2 := NewMessage(RoleAssistant, "Hi", mockNow)
				msg2.Tokens = &TokenCount{Output: 200}
				conv.AddMessage(msg1)
				conv.AddMessage(msg2)
				s.AddConversation(conv)
			},
			wantResult: 150.0,
			wantOK:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := NewStatistics()
			tt.setup(stats)

			got, ok := stats.AverageTokensPerMessage()
			assert.Equal(t, tt.wantOK, ok)
			if ok {
				assert.InDelta(t, tt.wantResult, got, 0.01)
			}
		})
	}
}

func TestStatistics_TrackModel(t *testing.T) {
	stats := NewStatistics()

	conv1 := NewConversation("conv-1", "Test", "claude-3-sonnet", mockNow)
	conv2 := NewConversation("conv-2", "Test", "claude-3-opus", mockNow)
	conv3 := NewConversation("conv-3", "Test", "claude-3-sonnet", mockNow)

	stats.AddConversation(conv1)
	stats.AddConversation(conv2)
	stats.AddConversation(conv3)

	assert.Equal(t, 2, stats.ModelDistribution["claude-3-sonnet"])
	assert.Equal(t, 1, stats.ModelDistribution["claude-3-opus"])
}

func TestStatistics_TrackRole(t *testing.T) {
	stats := NewStatistics()
	conv := NewConversation("conv-1", "Test", "claude-3-sonnet", mockNow)

	conv.AddMessage(NewMessage(RoleUser, "Q1", mockNow))
	conv.AddMessage(NewMessage(RoleAssistant, "A1", mockNow))
	conv.AddMessage(NewMessage(RoleUser, "Q2", mockNow))
	conv.AddMessage(NewMessage(RoleAssistant, "A2", mockNow))
	conv.AddMessage(NewMessage(RoleTool, "result", mockNow))

	stats.AddConversation(conv)

	assert.Equal(t, 2, stats.RoleDistribution[string(RoleUser)])
	assert.Equal(t, 2, stats.RoleDistribution[string(RoleAssistant)])
	assert.Equal(t, 1, stats.RoleDistribution[string(RoleTool)])
}

func TestStatistics_Empty(t *testing.T) {
	stats := NewStatistics()

	assert.Equal(t, 0, stats.ConversationCount)
	assert.Equal(t, 0, stats.MessageCount)
	assert.Equal(t, int64(0), stats.TotalTokens)
	assert.Equal(t, 0.0, stats.TotalCost)
	assert.Equal(t, 0.0, stats.AverageMessagesPerConversation())
	avg, ok := stats.AverageTokensPerMessage()
	assert.False(t, ok)
	assert.Equal(t, 0.0, avg)
}
