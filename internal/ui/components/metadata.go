// ABOUTME: Metadata/stats panel component for displaying conversation details
package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/TrevorS/claudex/internal/domain"
)

// Metadata displays conversation metadata and statistics in the right pane.
type Metadata struct {
	conversation *domain.Conversation
	width        int
}

// NewMetadata creates a new Metadata component.
func NewMetadata() *Metadata {
	return &Metadata{
		width: 30, // Default width
	}
}

// SetSize sets the available width for the metadata panel.
func (m *Metadata) SetSize(width, height int) *Metadata {
	newM := *m
	newM.width = width
	return &newM
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
	title := "Metadata"
	return title + "\n" + underline(len(title)) + "\n\nSelect a conversation\nto see details"
}

// renderConversationMetadata renders the metadata for the current conversation.
func (m *Metadata) renderConversationMetadata() string {
	var content string

	// Title - truncate to available width
	maxTitleLen := m.width - 2
	if maxTitleLen < 10 {
		maxTitleLen = 10
	}
	titleText := truncateString(m.conversation.Title, maxTitleLen)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("51"))
	content += titleStyle.Render(titleText) + "\n"
	content += underline(len(titleText)) + "\n\n"

	// Basic info
	content += m.renderInfoLine("Date:", m.formatDate(m.conversation.CreatedAt))
	content += m.renderInfoLine("Messages:", fmt.Sprintf("%d", m.conversation.MessageCount()))
	content += m.renderInfoLine("Tokens:", m.formatTokens(int(m.conversation.TotalTokens())))
	content += m.renderInfoLine("Model:", m.conversation.Model)

	// Statistics
	statsHeader := "Statistics"
	content += "\n\n" + statsHeader + "\n" + underline(len(statsHeader)) + "\n"
	content += m.renderInfoLine("User Msgs:", fmt.Sprintf("%d", m.conversation.UserMessageCount()))
	content += m.renderInfoLine("Avg Msg:", m.formatAvgTokens())

	// Token breakdown
	if totalTokens := m.conversation.TotalTokens(); totalTokens > 0 {
		inputTokens, outputTokens, _, _ := m.getTokenBreakdown()
		breakdownHeader := "Token Breakdown"
		content += "\n\n" + breakdownHeader + "\n" + underline(len(breakdownHeader)) + "\n"
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

	// Truncate value to fit available width (width - keyWidth - spacing)
	maxValueLen := m.width - 14
	if maxValueLen < 5 {
		maxValueLen = 5
	}
	if len(value) > maxValueLen {
		value = value[:maxValueLen-3] + "..."
	}

	return keyStyle.Render(key) + " " + value + "\n"
}

// formatDate formats a timestamp for display.
func (m *Metadata) formatDate(t time.Time) string {
	if t.IsZero() {
		return "Unknown"
	}
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

// underline creates an underline of the given length.
func underline(length int) string {
	if length < 3 {
		length = 3
	}
	return strings.Repeat("─", length)
}
