package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConversation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		id    string
		title string
		model string
		ts    time.Time
		want  *Conversation
	}{
		{
			name:  "basic conversation",
			id:    "conv-123",
			title: "Planning Project",
			model: "claude-3-sonnet",
			ts:    now,
			want: &Conversation{
				ID:        "conv-123",
				Title:     "Planning Project",
				Model:     "claude-3-sonnet",
				CreatedAt: now,
				UpdatedAt: now,
				Messages:  make([]*Message, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConversation(tt.id, tt.title, tt.model, tt.ts)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Title, got.Title)
			assert.Equal(t, tt.want.Model, got.Model)
			assert.Equal(t, tt.want.CreatedAt, got.CreatedAt)
			assert.Equal(t, tt.want.UpdatedAt, got.UpdatedAt)
			assert.NotNil(t, got.Messages)
		})
	}
}

func TestConversation_AddMessage(t *testing.T) {
	conv := NewConversation("conv-123", "Test", "claude-3-sonnet", time.Now())
	msg := NewMessage(RoleUser, "Hello", time.Now())

	conv.AddMessage(msg)

	assert.Len(t, conv.Messages, 1)
	assert.Equal(t, msg, conv.Messages[0])
}

func TestConversation_AddMultipleMessages(t *testing.T) {
	conv := NewConversation("conv-123", "Test", "claude-3-sonnet", time.Now())

	msg1 := NewMessage(RoleUser, "Hello", time.Now())
	msg2 := NewMessage(RoleAssistant, "Hi there!", time.Now().Add(1*time.Second))
	msg3 := NewMessage(RoleUser, "How are you?", time.Now().Add(2*time.Second))

	conv.AddMessage(msg1)
	conv.AddMessage(msg2)
	conv.AddMessage(msg3)

	assert.Len(t, conv.Messages, 3)
	assert.Equal(t, RoleUser, conv.Messages[0].Role)
	assert.Equal(t, RoleAssistant, conv.Messages[1].Role)
	assert.Equal(t, RoleUser, conv.Messages[2].Role)
}

func TestConversation_MessageCount(t *testing.T) {
	tests := []struct {
		name     string
		msgCount int
		want     int
	}{
		{
			name:     "empty conversation",
			msgCount: 0,
			want:     0,
		},
		{
			name:     "single message",
			msgCount: 1,
			want:     1,
		},
		{
			name:     "multiple messages",
			msgCount: 10,
			want:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := NewConversation("conv-123", "Test", "claude-3-sonnet", time.Now())
			for i := 0; i < tt.msgCount; i++ {
				conv.AddMessage(NewMessage(RoleUser, "msg", time.Now()))
			}

			got := conv.MessageCount()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConversation_TotalTokens(t *testing.T) {
	conv := NewConversation("conv-123", "Test", "claude-3-sonnet", time.Now())

	msg1 := NewMessage(RoleUser, "Hello", time.Now())
	msg1.Tokens = &TokenCount{Input: 10, Output: 0}

	msg2 := NewMessage(RoleAssistant, "Hi", time.Now())
	msg2.Tokens = &TokenCount{Input: 0, Output: 20}

	conv.AddMessage(msg1)
	conv.AddMessage(msg2)

	got := conv.TotalTokens()
	assert.Equal(t, int64(30), got)
}

func TestConversation_TotalTokens_NoMessages(t *testing.T) {
	conv := NewConversation("conv-123", "Test", "claude-3-sonnet", time.Now())

	got := conv.TotalTokens()
	assert.Equal(t, int64(0), got)
}

func TestConversation_UpdatedAtChanges(t *testing.T) {
	conv := NewConversation("conv-123", "Test", "claude-3-sonnet", time.Now())
	originalUpdated := conv.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	msg := NewMessage(RoleUser, "Hello", time.Now())
	conv.AddMessage(msg)

	assert.True(t, conv.UpdatedAt.After(originalUpdated))
}

func TestConversation_Validate(t *testing.T) {
	tests := []struct {
		name    string
		conv    *Conversation
		wantErr bool
	}{
		{
			name:    "valid conversation",
			conv:    NewConversation("conv-123", "Test Title", "claude-3-sonnet", time.Now()),
			wantErr: false,
		},
		{
			name: "missing ID",
			conv: &Conversation{
				ID:        "",
				Title:     "Test",
				Model:     "claude-3-sonnet",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing title",
			conv: &Conversation{
				ID:        "conv-123",
				Title:     "",
				Model:     "claude-3-sonnet",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing model",
			conv: &Conversation{
				ID:        "conv-123",
				Title:     "Test",
				Model:     "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.conv.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConversation_UserMessageCount(t *testing.T) {
	conv := NewConversation("conv-123", "Test", "claude-3-sonnet", time.Now())

	conv.AddMessage(NewMessage(RoleUser, "Q1", time.Now()))
	conv.AddMessage(NewMessage(RoleAssistant, "A1", time.Now()))
	conv.AddMessage(NewMessage(RoleUser, "Q2", time.Now()))
	conv.AddMessage(NewMessage(RoleAssistant, "A2", time.Now()))
	conv.AddMessage(NewMessage(RoleTool, "result", time.Now()))

	got := conv.UserMessageCount()
	assert.Equal(t, 2, got)
}

func TestConversation_LastMessage(t *testing.T) {
	conv := NewConversation("conv-123", "Test", "claude-3-sonnet", time.Now())

	lastMsg := NewMessage(RoleAssistant, "Final response", time.Now())
	conv.AddMessage(NewMessage(RoleUser, "First", time.Now()))
	conv.AddMessage(lastMsg)

	got := conv.LastMessage()
	require.NotNil(t, got)
	assert.Equal(t, lastMsg.Content, got.Content)
}

func TestConversation_LastMessage_Empty(t *testing.T) {
	conv := NewConversation("conv-123", "Test", "claude-3-sonnet", time.Now())

	got := conv.LastMessage()
	assert.Nil(t, got)
}
