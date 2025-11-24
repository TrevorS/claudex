// ABOUTME: Metadata/stats panel component for displaying conversation details
package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/TrevorS/claudex/internal/domain"
)

// Metadata displays conversation metadata and statistics in the right pane.
type Metadata struct {
	conversation *domain.Conversation
}

// NewMetadata creates a new Metadata component.
func NewMetadata() *Metadata {
	return &Metadata{}
}

// SetConversation sets the conversation to display.
func (m *Metadata) SetConversation(conv *domain.Conversation) *Metadata {
	newM := *m
	newM.conversation = conv
	return &newM
}

// View renders the metadata panel.
func (m *Metadata) View() string {
	if m.conversation == nil {
		return m.renderEmpty()
	}

	return m.renderConversationMetadata()
}

// renderEmpty returns an empty state message.
func (m *Metadata) renderEmpty() string {
	return "Metadata\n─────────\n\nSelect a conversation\nto see details"
}

// renderConversationMetadata renders the metadata for the current conversation.
func (m *Metadata) renderConversationMetadata() string {
	var content string

	// Title
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("51"))
	content += titleStyle.Render(truncateString(m.conversation.Title, 30)) + "\n"
	content += "─────────────────────────────\n\n"

	// Basic info
	content += m.renderInfoLine("Date:", m.formatDate(m.conversation.CreatedAt))
	content += m.renderInfoLine("Messages:", fmt.Sprintf("%d", m.conversation.MessageCount()))
	content += m.renderInfoLine("Tokens:", m.formatTokens(int(m.conversation.TotalTokens())))
	content += m.renderInfoLine("Model:", m.conversation.Model)

	// Statistics
	content += "\n\nStatistics\n───────────\n"
	content += m.renderInfoLine("User Msgs:", fmt.Sprintf("%d", m.conversation.UserMessageCount()))
	content += m.renderInfoLine("Avg Msg:", m.formatAvgTokens())

	// Token breakdown
	if totalTokens := m.conversation.TotalTokens(); totalTokens > 0 {
		inputTokens, outputTokens, _, _ := m.getTokenBreakdown()
		content += "\n\nToken Breakdown\n────────────────\n"
		content += m.renderInfoLine("Input:", m.formatTokens(int(inputTokens)))
		content += m.renderInfoLine("Output:", m.formatTokens(int(outputTokens)))
	}

	return content
}

// renderInfoLine renders a key-value info line.
func (m *Metadata) renderInfoLine(key, value string) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Dim gray
		Width(12)
	return keyStyle.Render(key) + " " + value + "\n"
}

// formatDate formats a timestamp for display.
func (m *Metadata) formatDate(t time.Time) string {
	return t.Format("Jan 2, 2006 03:04 PM")
}

// formatTokens formats token count with K/M suffix.
func (m *Metadata) formatTokens(tokens int) string {
	if tokens >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(tokens)/1_000_000)
	}
	if tokens >= 1_000 {
		return fmt.Sprintf("%.1fK", float64(tokens)/1_000)
	}
	return fmt.Sprintf("%d", tokens)
}

// formatAvgTokens calculates and formats average tokens per message.
func (m *Metadata) formatAvgTokens() string {
	count := int64(m.conversation.MessageCount())
	if count == 0 {
		return "0"
	}
	avg := m.conversation.TotalTokens() / count
	return fmt.Sprintf("%d", avg)
}

// getTokenBreakdown extracts token count breakdown from messages.
func (m *Metadata) getTokenBreakdown() (input, output, cacheRead, cacheCreate int64) {
	for _, msg := range m.conversation.Messages {
		if msg.Tokens != nil {
			input += msg.Tokens.Input
			output += msg.Tokens.Output
			cacheRead += msg.Tokens.CacheRead
			cacheCreate += msg.Tokens.CacheCreate
		}
	}
	return
}

// truncateString truncates a string to maxLen with ellipsis if needed.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
