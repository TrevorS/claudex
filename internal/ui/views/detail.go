// ABOUTME: DetailView component displays a full conversation with scrollable messages.
// Shows conversation header with metadata, all messages with role labels and token counts,
// and supports keyboard navigation including scrolling and returning to list view.
package views

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/TrevorS/claudex/internal/render"
	"github.com/TrevorS/claudex/internal/repository"
)

// DetailView displays a full conversation with scrollable messages.
type DetailView struct {
	repository     repository.Repository
	conversation   *domain.Conversation
	err            *domain.AppError
	viewport       viewport.Model
	content        string // pre-rendered message content
	conversationID string
	viewHeight     int
	viewWidth      int
	ready          bool // viewport initialized
	mdRenderer     *render.Renderer
}

// ConversationLoadedMsg is sent when a conversation is loaded from the repository.
type ConversationLoadedMsg struct {
	Conversation *domain.Conversation
	Err          error
}

// BackToListMsg is sent when the user wants to return to the list view.
type BackToListMsg struct{}

// Styles for message rendering
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")). // Bright blue
			Margin(0, 0, 1, 0)

	userLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14")). // Bright cyan
			Margin(0, 1, 0, 0)

	assistantLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("10")). // Bright green
				Margin(0, 1, 0, 0)

	toolLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("11")). // Bright yellow
			Margin(0, 1, 0, 0)

	tokenStyle = lipgloss.NewStyle().
			Faint(true).
			Foreground(lipgloss.Color("8")). // Gray
			Margin(0, 0, 0, 1)

	messageContentStyle = lipgloss.NewStyle().
				Margin(0, 0, 1, 2) // Indent message content

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")). // Gray
			Margin(1, 0, 1, 0)

	metadataStyle = lipgloss.NewStyle().
			Faint(true).
			Foreground(lipgloss.Color("8"))

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")). // Bright red
			Margin(1, 0, 1, 0)
)

// NewDetailView creates a new DetailView for the given conversation ID.
func NewDetailView(repo repository.Repository, conversationID string) *DetailView {
	return &DetailView{
		repository:     repo,
		conversationID: conversationID,
		viewWidth:      80,
		viewHeight:     24,
		ready:          false,
		mdRenderer:     render.NewRenderer(80),
	}
}

// Init initializes the DetailView and loads the conversation.
func (v *DetailView) Init() tea.Cmd {
	return v.loadConversationCmd()
}

// loadConversationCmd returns a command that loads the conversation from the repository.
func (v *DetailView) loadConversationCmd() tea.Cmd {
	return func() tea.Msg {
		conversation, err := v.repository.GetByID(context.Background(), v.conversationID)
		return ConversationLoadedMsg{
			Conversation: conversation,
			Err:          err,
		}
	}
}

// Update handles messages and updates the DetailView state.
func (v *DetailView) Update(msg tea.Msg) (*DetailView, tea.Cmd) {
	switch msg := msg.(type) {

	case ConversationLoadedMsg:
		if msg.Err != nil {
			newView := *v
			if appErr, ok := msg.Err.(*domain.AppError); ok {
				newView.err = appErr
			} else {
				newView.err = domain.NewAppError(domain.ErrRepository, "failed to load conversation", msg.Err)
			}
			return &newView, nil
		}

		newView := *v
		newView.conversation = msg.Conversation
		newView.content = newView.renderConversation()
		newView = *newView.initViewport()
		return &newView, nil

	case tea.WindowSizeMsg:
		return v.SetSize(msg.Width, msg.Height), nil

	case tea.KeyMsg:
		return v.handleKeyPress(msg)

	default:
		// Pass other messages to viewport
		if v.ready {
			newView := *v
			var cmd tea.Cmd
			newView.viewport, cmd = v.viewport.Update(msg)
			return &newView, cmd
		}
		return v, nil
	}
}

// handleKeyPress handles keyboard input.
func (v *DetailView) handleKeyPress(msg tea.KeyMsg) (*DetailView, tea.Cmd) {
	switch msg.Type {

	case tea.KeyEsc:
		return v, func() tea.Msg {
			return BackToListMsg{}
		}

	case tea.KeyRunes:
		if len(msg.Runes) > 0 {
			switch msg.Runes[0] {
			case 'b':
				return v, func() tea.Msg {
					return BackToListMsg{}
				}
			}
		}
	}

	// Pass to viewport for scrolling
	if v.ready {
		newView := *v
		var cmd tea.Cmd
		newView.viewport, cmd = v.viewport.Update(msg)
		return &newView, cmd
	}

	return v, nil
}

// initViewport initializes the viewport with content.
func (v *DetailView) initViewport() *DetailView {
	newView := *v

	// Calculate content height (leave room for detail header/footer only)
	// Note: viewHeight already has border/padding subtracted by parent
	contentHeight := v.viewHeight - 2
	if contentHeight < 10 {
		contentHeight = 10
	}

	vp := viewport.New(v.viewWidth, contentHeight)
	vp.SetContent(v.content)
	vp.YPosition = 0

	newView.viewport = vp
	newView.ready = true

	return &newView
}

// View renders the DetailView as a string.
func (v *DetailView) View() string {
	if v.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error loading conversation: %s", v.err.Error()))
	}

	if v.conversation == nil {
		return errorStyle.Render("Conversation not found.")
	}

	if !v.ready {
		return "Loading..."
	}

	return v.viewport.View()
}

// renderConversation renders the entire conversation including header and messages.
func (v *DetailView) renderConversation() string {
	if v.conversation == nil {
		return ""
	}

	var parts []string

	// Render header
	parts = append(parts, v.renderHeader())

	// Render messages
	if len(v.conversation.Messages) == 0 {
		parts = append(parts, "\nNo messages in this conversation.\n")
	} else {
		for i, msg := range v.conversation.Messages {
			parts = append(parts, v.renderMessage(msg, i))
		}
	}

	return strings.Join(parts, "\n")
}

// renderHeader renders the conversation header with metadata.
func (v *DetailView) renderHeader() string {
	conv := v.conversation

	var parts []string

	// Title
	title := headerStyle.Render(conv.Title)
	parts = append(parts, title)

	// Metadata line
	metadata := fmt.Sprintf("Model: %s  |  Messages: %d  |  Tokens: %s  |  Updated: %s",
		v.formatModel(conv.Model),
		len(conv.Messages),
		v.formatTokenCount(conv.TotalTokens()),
		v.formatTimestamp(conv.UpdatedAt),
	)
	parts = append(parts, metadataStyle.Render(metadata))

	// Separator
	separator := strings.Repeat("─", v.viewWidth-4)
	parts = append(parts, separatorStyle.Render(separator))

	return strings.Join(parts, "\n")
}

// renderMessage renders a single message with role label, content, and token count.
func (v *DetailView) renderMessage(msg *domain.Message, index int) string {
	var parts []string

	// Role label with token count
	roleLabel := v.formatRole(msg.Role)
	tokens := v.formatTokenCount(msg.TotalTokens())
	headerLine := fmt.Sprintf("%s  %s", roleLabel, tokenStyle.Render(tokens))
	parts = append(parts, headerLine)

	// Message content - render markdown to terminal format
	renderedContent := v.mdRenderer.Render(msg.Content)
	parts = append(parts, messageContentStyle.Render(renderedContent))

	return strings.Join(parts, "\n")
}

// formatRole formats the message role as a labeled string.
func (v *DetailView) formatRole(role domain.Role) string {
	switch role {
	case domain.RoleUser:
		return userLabelStyle.Render("● User")
	case domain.RoleAssistant:
		return assistantLabelStyle.Render("◆ Assistant")
	case domain.RoleTool:
		return toolLabelStyle.Render("▶ Tool")
	default:
		return role.String()
	}
}

// formatTokenCount formats a token count for display.
func (v *DetailView) formatTokenCount(tokens int64) string {
	if tokens == 0 {
		return ""
	}

	if tokens < 1000 {
		return fmt.Sprintf("⊛ %d tokens", tokens)
	}

	if tokens < 1000000 {
		return fmt.Sprintf("⊛ %.1fK tokens", float64(tokens)/1000)
	}

	return fmt.Sprintf("⊛ %.1fM tokens", float64(tokens)/1000000)
}

// formatModel simplifies model names for display.
func (v *DetailView) formatModel(model string) string {
	// Simplify "claude-sonnet-4" -> "sonnet-4"
	parts := strings.Split(model, "-")
	if len(parts) > 1 && parts[0] == "claude" {
		return strings.Join(parts[1:], "-")
	}
	return model
}

// formatTimestamp formats a timestamp for display.
func (v *DetailView) formatTimestamp(ts time.Time) string {
	now := time.Now()
	diff := now.Sub(ts)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	default:
		return ts.Format("2006-01-02 15:04")
	}
}

// wrapText wraps text to the specified width, preserving code blocks and paragraphs.
func (v *DetailView) wrapText(text string, width int) string {
	if width < 20 {
		width = 20 // Minimum readable width
	}

	lines := strings.Split(text, "\n")
	var wrapped []string

	for _, line := range lines {
		// Preserve empty lines
		if len(line) == 0 {
			wrapped = append(wrapped, "")
			continue
		}

		// Preserve code blocks (lines starting with spaces or tabs)
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			wrapped = append(wrapped, line)
			continue
		}

		// Wrap long lines
		if len(line) <= width {
			wrapped = append(wrapped, line)
		} else {
			wrapped = append(wrapped, v.wrapLine(line, width)...)
		}
	}

	return strings.Join(wrapped, "\n")
}

// wrapLine wraps a single long line to the specified width.
func (v *DetailView) wrapLine(line string, width int) []string {
	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{""}
	}

	var wrapped []string
	currentLine := ""

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			wrapped = append(wrapped, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		wrapped = append(wrapped, currentLine)
	}

	return wrapped
}

// SetSize sets the view dimensions and reinitializes the viewport.
func (v *DetailView) SetSize(width, height int) *DetailView {
	newView := *v
	newView.viewWidth = width
	newView.viewHeight = height

	// Update markdown renderer width
	newView.mdRenderer = render.NewRenderer(width - 4) // Account for indentation

	// Reinitialize viewport if already set up
	if newView.ready && newView.conversation != nil {
		newView.content = newView.renderConversation()
		newView = *newView.initViewport()
	}

	return &newView
}
