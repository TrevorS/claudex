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
	sidebar := m.renderSidebar()
	center := m.renderCenter()
	right := m.renderRight()

	// Calculate pane widths with padding
	sidebarWidth := 22
	rightWidth := 38
	centerWidth := m.width - sidebarWidth - rightWidth - 6 // Account for borders/spacing

	if centerWidth < 20 {
		centerWidth = 20
	}

	// Apply width constraints
	sidebar = lipgloss.NewStyle().Width(sidebarWidth).Render(sidebar)
	center = lipgloss.NewStyle().Width(centerWidth).Render(center)
	right = lipgloss.NewStyle().Width(rightWidth).Render(right)

	// Join with borders
	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		renderVerticalSeparator(m.height-2),
		center,
		renderVerticalSeparator(m.height-2),
		right,
	)

	return body
}

// renderSidebar renders the left sidebar.
func (m Model) renderSidebar() string {
	sidebarStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("59")). // Gray
		Padding(1, 1)

	var content string
	if m.sidebar != nil {
		content = m.sidebar.View()
	} else {
		content = "Projects\n─────────\n\n(Loading...)"
	}

	return sidebarStyle.Render(content)
}

// renderCenter renders the center pane.
func (m Model) renderCenter() string {
	centerStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("59")). // Gray
		Padding(1, 1)

	content := m.renderCenterContent()

	return centerStyle.Render(content)
}

// renderCenterContent renders the content based on current state.
func (m Model) renderCenterContent() string {
	switch m.state {
	case ViewList:
		if m.listView != nil {
			return m.listView.View()
		}
		return "Conversations\n──────────────\n\n(Loading...)"
	case ViewDetail:
		if m.detailView != nil {
			return m.detailView.View()
		}
		return "Message Detail\n───────────────\n\n(Loading...)"
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

// renderRight renders the right pane.
func (m Model) renderRight() string {
	rightStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("59")). // Gray
		Padding(1, 1)

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
