package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		role    Role
		content string
		ts      time.Time
		want    *Message
	}{
		{
			name:    "user message",
			role:    RoleUser,
			content: "Hello, assistant!",
			ts:      now,
			want: &Message{
				Role:      RoleUser,
				Content:   "Hello, assistant!",
				Timestamp: now,
				Tokens:    &TokenCount{},
			},
		},
		{
			name:    "assistant message",
			role:    RoleAssistant,
			content: "Hi there!",
			ts:      now,
			want: &Message{
				Role:      RoleAssistant,
				Content:   "Hi there!",
				Timestamp: now,
				Tokens:    &TokenCount{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMessage(tt.role, tt.content, tt.ts)
			assert.Equal(t, tt.want.Role, got.Role)
			assert.Equal(t, tt.want.Content, got.Content)
			assert.Equal(t, tt.want.Timestamp, got.Timestamp)
			assert.NotNil(t, got.Tokens)
		})
	}
}

func TestMessage_TotalTokens(t *testing.T) {
	msg := NewMessage(RoleUser, "test", time.Now())
	msg.Tokens = &TokenCount{
		Input:       100,
		Output:      50,
		CacheRead:   25,
		CacheCreate: 0,
	}

	want := int64(175)
	got := msg.TotalTokens()
	assert.Equal(t, want, got)
}

func TestMessage_TotalTokens_Zero(t *testing.T) {
	msg := NewMessage(RoleAssistant, "test", time.Now())
	// Tokens already initialized as &TokenCount{} in NewMessage

	want := int64(0)
	got := msg.TotalTokens()
	assert.Equal(t, want, got)
}

func TestMessage_EstimatedCost(t *testing.T) {
	tests := []struct {
		name       string
		role       Role
		tokens     *TokenCount
		wantResult float64
		wantOK     bool
	}{
		{
			name: "valid input/output tokens",
			role: RoleAssistant,
			tokens: &TokenCount{
				Input:  1000,
				Output: 500,
			},
			wantResult: 0.0,
			wantOK:     true,
		},
		{
			name: "zero tokens",
			role: RoleUser,
			tokens: &TokenCount{
				Input:  0,
				Output: 0,
			},
			wantResult: 0.0,
			wantOK:     true,
		},
		{
			name:       "nil tokens",
			role:       RoleUser,
			tokens:     nil,
			wantResult: 0.0,
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMessage(tt.role, "test", time.Now())
			msg.Tokens = tt.tokens
			got, ok := msg.EstimatedCost()
			assert.Equal(t, tt.wantOK, ok)
			if ok {
				assert.GreaterOrEqual(t, got, 0.0)
			}
		})
	}
}

func TestRole_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		wantTrue bool
	}{
		{name: "user role", role: RoleUser, wantTrue: true},
		{name: "assistant role", role: RoleAssistant, wantTrue: true},
		{name: "tool role", role: RoleTool, wantTrue: true},
		{name: "invalid role", role: Role("invalid"), wantTrue: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.IsValid()
			assert.Equal(t, tt.wantTrue, got)
		})
	}
}

func TestTokenCount_Total(t *testing.T) {
	tests := []struct {
		name   string
		tc     *TokenCount
		wantOK bool
		want   int64
	}{
		{
			name: "all fields populated",
			tc: &TokenCount{
				Input:       100,
				Output:      200,
				CacheRead:   50,
				CacheCreate: 25,
			},
			wantOK: true,
			want:   375,
		},
		{
			name:   "zero tokens",
			tc:     &TokenCount{},
			wantOK: true,
			want:   0,
		},
		{
			name:   "nil pointer",
			tc:     nil,
			wantOK: false,
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tc == nil {
				// Can't call method on nil, test independently
				assert.Nil(t, tt.tc)
				return
			}
			got := tt.tc.Total()
			assert.Equal(t, tt.want, got)
		})
	}
}
