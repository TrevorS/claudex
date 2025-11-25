// ABOUTME: Three-pane layout composition with Lipgloss
package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View returns the complete view string for the current UI state.
func (m Model) View() string {
	if m.err != nil {
		return m.renderErrorView()
	}

	if m.width < 40 {
		return m.renderMinimalView()
	}

	if m.width < 80 {
		return m.renderNarrowView()
	}

	return m.renderWideView()
}

// renderErrorView displays an error message.
func (m Model) renderErrorView() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")). // Red
		Bold(true)

	content := fmt.Sprintf("Error: %v\n\nPress Ctrl+C to exit", m.err)

	return errorStyle.Render(content)
}

// renderMinimalView renders a single-pane view for very narrow terminals.
func (m Model) renderMinimalView() string {
	header := m.renderHeader()
	content := m.renderContent()
	footer := m.renderFooter()

	return fmt.Sprintf("%s\n%s\n%s", header, content, footer)
}

// renderNarrowView renders a simplified two-pane view for narrow terminals.
func (m Model) renderNarrowView() string {
	header := m.renderHeader()
	content := m.renderContent()
	footer := m.renderFooter()

	return fmt.Sprintf("%s\n%s\n%s", header, content, footer)
}

// renderWideView renders the full three-pane layout.
func (m Model) renderWideView() string {
	header := m.renderHeader()
	body := m.renderThreePane()
	footer := m.renderFooter()

	return fmt.Sprintf("%s\n%s\n%s", header, body, footer)
}

// renderHeader renders the header section.
func (m Model) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("63")).  // Purple
		Foreground(lipgloss.Color("255")). // White
		Bold(true).
		Padding(0, 1)

	title := "claudex • conversation browser"

	return headerStyle.Render(title)
}

// renderThreePane renders the three-pane layout body.
func (m Model) renderThreePane() string {
	// Calculate content height (total - header - footer)
	// Header = 1 line, Footer = 1 line
	contentHeight := m.height - 2
	if contentHeight < 5 {
		contentHeight = 5
	}

	// If sidebar is hidden (detail view), use two-pane layout
	if !m.sidebarVisible {
		return m.renderTwoPane(contentHeight)
	}

	// Calculate pane widths as percentages
	sidebarWidth := m.width * 25 / 100 // 25% for sidebar (wider for paths)
	rightWidth := m.width * 22 / 100   // 22% for right panel
	centerWidth := m.width - sidebarWidth - rightWidth

	// Enforce minimum widths
	if sidebarWidth < 25 {
		sidebarWidth = 25
	}
	if rightWidth < 22 {
		rightWidth = 22
	}
	if centerWidth < 30 {
		centerWidth = 30
	}

	// Render panes with explicit dimensions
	sidebar := m.renderSidebarWithSize(sidebarWidth, contentHeight)
	center := m.renderCenterWithSize(centerWidth, contentHeight)
	right := m.renderRightWithSize(rightWidth, contentHeight)

	// Join horizontally (borders provide visual separation)
	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		center,
		right,
	)

	return body
}

// renderTwoPane renders a two-pane layout (center + right) when sidebar is hidden.
func (m Model) renderTwoPane(contentHeight int) string {
	// Center gets more space without sidebar
	rightWidth := m.width * 25 / 100
	centerWidth := m.width - rightWidth

	// Enforce minimum widths
	if rightWidth < 25 {
		rightWidth = 25
	}
	if centerWidth < 40 {
		centerWidth = 40
	}

	center := m.renderCenterWithSize(centerWidth, contentHeight)
	right := m.renderRightWithSize(rightWidth, contentHeight)

	return lipgloss.JoinHorizontal(lipgloss.Top, center, right)
}

// renderSidebarWithSize renders the sidebar with explicit dimensions.
func (m Model) renderSidebarWithSize(width, height int) string {
	borderColor := lipgloss.Color("59") // Gray (unfocused)
	if m.sidebar != nil && m.sidebar.HasFocus() {
		borderColor = lipgloss.Color("6") // Cyan (focused)
	}

	// Inner dimensions account for border (2) and padding (2)
	innerWidth := width - 4
	innerHeight := height - 2
	if innerWidth < 10 {
		innerWidth = 10
	}
	if innerHeight < 3 {
		innerHeight = 3
	}

	var content string
	if m.sidebar != nil {
		content = m.sidebar.View()
	} else {
		content = "Projects\n─────────\n\n(Loading...)"
	}

	style := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	return style.Render(content)
}

// renderSidebar renders the left sidebar with focus-aware border (for narrow views).
func (m Model) renderSidebar() string {
	borderColor := lipgloss.Color("59") // Gray (unfocused)
	if m.sidebar != nil && m.sidebar.HasFocus() {
		borderColor = lipgloss.Color("6") // Cyan (focused)
	}

	sidebarStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	var content string
	if m.sidebar != nil {
		content = m.sidebar.View()
	} else {
		content = "Projects\n─────────\n\n(Loading...)"
	}

	return sidebarStyle.Render(content)
}

// renderCenterWithSize renders the center pane with explicit dimensions.
func (m Model) renderCenterWithSize(width, height int) string {
	borderColor := lipgloss.Color("59") // Gray (unfocused)
	if m.focusPane == FocusCenter {
		borderColor = lipgloss.Color("6") // Cyan (focused)
	}

	// Inner dimensions account for border (2) and padding (2)
	innerWidth := width - 4
	innerHeight := height - 2
	if innerWidth < 10 {
		innerWidth = 10
	}
	if innerHeight < 3 {
		innerHeight = 3
	}

	content := m.renderCenterContent()

	style := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	return style.Render(content)
}

// renderCenter renders the center pane with focus-aware border (for narrow views).
func (m Model) renderCenter() string {
	borderColor := lipgloss.Color("59") // Gray (unfocused)
	if m.focusPane == FocusCenter {
		borderColor = lipgloss.Color("6") // Cyan (focused)
	}

	centerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	content := m.renderCenterContent()

	return centerStyle.Render(content)
}

// renderCenterContent renders the content based on current state.
func (m Model) renderCenterContent() string {
	// Show loading indicator during initial load
	if m.loading && m.state == ViewList {
		return "Loading conversations..."
	}

	switch m.state {
	case ViewList:
		if m.listView != nil {
			return m.listView.View()
		}
		return "No conversations found."
	case ViewDetail:
		if m.detailView != nil {
			return m.detailView.View()
		}
		return "Loading conversation..."
	case ViewSearch:
		return m.renderSearchView()
	case ViewStats:
		return "Statistics\n───────────────\n\n(Phase 2)"
	default:
		return "Unknown View"
	}
}

// renderSearchView renders the search view with search bar and results.
func (m Model) renderSearchView() string {
	title := lipgloss.NewStyle().Bold(true).Render("Search Conversations")
	separator := "─────────────────────────"

	var content string
	if m.searchBar != nil {
		content = fmt.Sprintf("%s\n%s\n\n%s\n\n", title, separator, m.searchBar.View())
	} else {
		content = fmt.Sprintf("%s\n%s\n\n", title, separator)
	}

	// Show results or hints
	if m.isSearchActive && len(m.searchResults) > 0 {
		content += fmt.Sprintf("Results (%d):\n", len(m.searchResults))
		for i, conv := range m.searchResults {
			if i >= 10 { // Show first 10
				content += fmt.Sprintf("... and %d more\n", len(m.searchResults)-10)
				break
			}
			content += fmt.Sprintf("• %s\n", conv.Title)
		}
	} else if m.isSearchActive && len(m.searchResults) == 0 {
		content += "No results found.\n"
	} else {
		content += "Press Ctrl+F to search • Enter to execute • Esc to cancel\n\n"
		content += "Examples:\n"
		content += "  • python - simple text search\n"
		content += "  • title:\"my chat\" - search in title\n"
		content += "  • tokens:>5000 - filter by token count\n"
		content += "  • python AND rust - boolean AND\n"
		content += "  • (ai OR ml) AND NOT archived - complex query\n"
	}

	return content
}

// renderRightWithSize renders the right pane with explicit dimensions.
func (m Model) renderRightWithSize(width, height int) string {
	borderColor := lipgloss.Color("59") // Gray (unfocused)
	if m.focusPane == FocusRight {
		borderColor = lipgloss.Color("6") // Cyan (focused)
	}

	// Inner dimensions account for border (2) and padding (2)
	innerWidth := width - 4
	innerHeight := height - 2
	if innerWidth < 10 {
		innerWidth = 10
	}
	if innerHeight < 3 {
		innerHeight = 3
	}

	var content string
	if m.metadata != nil {
		content = m.metadata.View()
	} else {
		content = "Metadata\n─────────\n\n(Loading...)"
	}

	style := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	return style.Render(content)
}

// renderRight renders the right pane with focus-aware border (for narrow views).
func (m Model) renderRight() string {
	borderColor := lipgloss.Color("59") // Gray (unfocused)
	if m.focusPane == FocusRight {
		borderColor = lipgloss.Color("6") // Cyan (focused)
	}

	rightStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	var content string
	if m.metadata != nil {
		content = m.metadata.View()
	} else {
		content = "Metadata\n─────────\n\n(Loading...)"
	}

	return rightStyle.Render(content)
}

// renderContent renders the main content area (fallback).
func (m Model) renderContent() string {
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("59")).
		Padding(1, 1)

	text := "Claudex Conversation Browser\n\nPress Tab to navigate, Escape to go back, Ctrl+C to quit"

	return contentStyle.Render(text)
}

// renderFooter renders the footer section.
func (m Model) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("59")). // Gray
		Padding(0, 1)

	hints := "^C:quit  ^F:search  Tab:focus  Esc:back  ^P:palette"

	return footerStyle.Render(hints)
}

// renderVerticalSeparator creates a vertical line.
func renderVerticalSeparator(height int) string {
	if height <= 0 {
		return "│"
	}

	separator := "│"
	lines := make([]string, height)
	for i := range lines {
		lines[i] = separator
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// ErrAppError is a helper to create an error for display.
func ErrAppError(msg string) error {
	return fmt.Errorf("%s", msg)
}
