// ABOUTME: Centralized message routing and state transitions
package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/TrevorS/claudex/internal/ui/components"
	"github.com/TrevorS/claudex/internal/ui/views"
)

// Update handles all messages and returns a new model and command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case views.ConversationsLoadedMsg:
		// Handle conversations loaded from ListView
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}
		m.conversations = msg.Conversations
		m.loading = false // Done loading
		// Update sidebar with conversations
		if m.sidebar != nil {
			m.sidebar = m.sidebar.SetConversations(m.conversations)
		}
		// Forward to ListView
		if m.listView != nil && m.state == ViewList {
			m.listView, cmd = m.listView.Update(msg)
		}
		return m, cmd

	case views.SelectConversationMsg:
		// Handle conversation selection from ListView
		if msg.Conversation != nil {
			m.selected = msg.Conversation.ID
			m.state = ViewDetail
			m.sidebarVisible = false // Hide sidebar in detail view
			// Create and initialize DetailView
			m.detailView = views.NewDetailView(m.repository, msg.Conversation.ID)
			cmd = m.detailView.Init()
			// Update metadata panel
			if m.metadata != nil {
				m.metadata = m.metadata.SetConversation(msg.Conversation)
			}
		}
		return m, cmd

	case views.ConversationLoadedMsg:
		// Handle conversation loaded from DetailView
		if m.detailView != nil && m.state == ViewDetail {
			m.detailView, cmd = m.detailView.Update(msg)
			// Also update metadata with fully loaded conversation
			if msg.Conversation != nil && m.metadata != nil {
				m.metadata = m.metadata.SetConversation(msg.Conversation)
			}
		}
		return m, cmd

	case views.BackToListMsg:
		// Handle back to list from DetailView
		m.state = ViewList
		m.sidebarVisible = true // Show sidebar again
		m.detailView = nil
		return m, nil

	case components.SearchDebounceMsg:
		// Handle search debounce - execute search
		if msg.Query != "" {
			cmd = m.executeSearch(msg.Query)
			m = m.SetCurrentQuery(msg.Query)
		}
		return m, cmd

	case components.SearchResultsMsg:
		// Handle search results
		m = m.SetSearchResults(msg.Results).SetSearchActive(true)
		// Forward to search bar
		if m.searchBar != nil {
			m.searchBar, cmd = m.searchBar.Update(msg)
		}
		return m, cmd

	case components.SearchErrorMsg:
		// Handle search error
		m = m.SetError(nil) // Clear any previous errors
		// Forward to search bar
		if m.searchBar != nil {
			m.searchBar, cmd = m.searchBar.Update(msg)
		}
		return m, cmd

	case components.ProjectSelectedMsg:
		// Handle project selection from sidebar
		if msg.ProjectPath == "" {
			// "All" selected - show all conversations
			m = m.SetSearchResults(nil).SetSearchActive(false)
			if m.listView != nil {
				m.listView = m.listView.SetConversations(m.conversations)
			}
		} else {
			// Filter conversations to selected project
			var filtered []*domain.Conversation
			for _, conv := range m.conversations {
				if conv.ProjectPath == msg.ProjectPath {
					filtered = append(filtered, conv)
				}
			}
			m = m.SetSearchResults(filtered).SetSearchActive(true)
			if m.listView != nil {
				m.listView = m.listView.SetConversations(filtered)
			}
		}
		return m, nil

	case components.FilterSelectedMsg:
		// Handle filter selection from sidebar
		// Filter conversations based on selected filter
		filtered := m.filterConversationsByFilter(string(msg.Filter))
		m = m.SetSearchResults(filtered).SetSearchActive(true)
		if m.listView != nil {
			m.listView = m.listView.SetConversations(filtered)
		}
		return m, nil

	case views.CursorMovedMsg:
		// Handle cursor movement from ListView - update metadata
		if m.metadata != nil && msg.Conversation != nil {
			m.metadata = m.metadata.SetConversation(msg.Conversation)
		}
		return m, nil

	case tea.WindowSizeMsg:
		// Handle terminal resize
		m = m.handleWindowSize(msg)
		// Forward to active view
		if m.listView != nil && m.state == ViewList {
			m.listView, cmd = m.listView.Update(msg)
		}
		if m.detailView != nil && m.state == ViewDetail {
			m.detailView, cmd = m.detailView.Update(msg)
		}
		return m, cmd

	case tea.KeyMsg:
		// Handle keyboard input
		return m.handleKeyPress(msg)

	default:
		// Forward unknown messages to active view
		if m.listView != nil && m.state == ViewList {
			m.listView, cmd = m.listView.Update(msg)
			return m, cmd
		}
		if m.detailView != nil && m.state == ViewDetail {
			m.detailView, cmd = m.detailView.Update(msg)
			return m, cmd
		}
		return m, nil
	}
}

// handleWindowSize updates the model with new terminal dimensions.
func (m Model) handleWindowSize(msg tea.WindowSizeMsg) Model {
	m = m.SetSize(msg.Width, msg.Height)

	// Calculate content area height (total - header - footer)
	// Header = 1 line, Footer = 1 line
	contentHeight := m.height - 2
	if contentHeight < 5 {
		contentHeight = 5
	}

	// Calculate pane widths based on sidebar visibility
	var sidebarWidth, centerWidth, rightWidth int

	if m.sidebarVisible {
		// Three-pane layout
		sidebarWidth = m.width * 25 / 100 // 25% for sidebar (wider for paths)
		rightWidth = m.width * 22 / 100   // 22% for right panel
		centerWidth = m.width - sidebarWidth - rightWidth

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
	} else {
		// Two-pane layout (detail view)
		sidebarWidth = 0
		rightWidth = m.width * 25 / 100
		centerWidth = m.width - rightWidth

		if rightWidth < 25 {
			rightWidth = 25
		}
		if centerWidth < 40 {
			centerWidth = 40
		}
	}

	// Account for border and padding in center pane (border=2, padding=2)
	innerCenterWidth := centerWidth - 4
	innerCenterHeight := contentHeight - 4

	if innerCenterWidth < 10 {
		innerCenterWidth = 10
	}
	if innerCenterHeight < 3 {
		innerCenterHeight = 3
	}

	// Forward dimensions to all components
	if m.sidebar != nil && m.sidebarVisible {
		m.sidebar = m.sidebar.SetSize(sidebarWidth-4, contentHeight-4)
	}
	if m.listView != nil {
		m.listView = m.listView.SetSize(innerCenterWidth, innerCenterHeight)
	}
	if m.detailView != nil {
		m.detailView = m.detailView.SetSize(innerCenterWidth, innerCenterHeight)
	}
	if m.metadata != nil {
		m.metadata = m.metadata.SetSize(rightWidth-4, contentHeight-4)
	}

	return m
}

// handleKeyPress handles global keyboard shortcuts.
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {

	case tea.KeyCtrlC:
		// Quit application
		return m, tea.Quit

	case tea.KeyEsc:
		// Return to list view or close modal
		if m.state == ViewDetail {
			m.state = ViewList
			m.sidebarVisible = true // Show sidebar again
			m.detailView = nil
			return m, nil
		}
		if m.state == ViewSearch || m.state == ViewStats || m.state == ViewPalette {
			return m.SetState(ViewList), nil
		}
		return m, nil

	case tea.KeyTab:
		// Tab key for focus switching between panes
		return m.handleTabKey(), nil

	case tea.KeyCtrlP:
		// Command palette (Phase 2+)
		return m, nil

	case tea.KeyCtrlF:
		// Toggle search mode
		if m.state == ViewSearch {
			// Exit search mode
			m = m.SetState(ViewList).SetSearchActive(false).SetSearchResults(nil)
			if m.searchBar != nil {
				m.searchBar.SetValue("")
			}
		} else {
			// Enter search mode
			m = m.SetState(ViewSearch)
			if m.searchBar != nil {
				cmd = m.searchBar.Focus()
			}
		}
		return m, cmd

	case tea.KeyRunes:
		// Handle global keys before focus-based routing
		if len(msg.Runes) > 0 && m.state == ViewList {
			switch msg.Runes[0] {
			case 'f', 'F':
				// Global filter toggle - works regardless of focus
				if m.sidebar != nil {
					m.sidebar = m.sidebar.CycleFilter()
				}
				return m, nil
			}
		}
		// Fall through to focus-based routing below
		fallthrough

	default:
		// Forward keyboard input to active view based on focus
		if m.state == ViewSearch && m.searchBar != nil {
			m.searchBar, cmd = m.searchBar.Update(msg)
			return m, cmd
		}
		if m.state == ViewList {
			// Route to focused pane only
			switch m.focusPane {
			case FocusSidebar:
				if m.sidebar != nil {
					m.sidebar, cmd = m.sidebar.Update(msg)
				}
			case FocusCenter:
				if m.listView != nil {
					m.listView, cmd = m.listView.Update(msg)
				}
			}
			return m, cmd
		}
		if m.state == ViewDetail && m.detailView != nil {
			m.detailView, cmd = m.detailView.Update(msg)
			return m, cmd
		}
		return m, nil
	}
}

// executeSearch executes a search query and returns a command.
func (m Model) executeSearch(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := m.searchEngine.Search(context.Background(), query)
		if err != nil {
			return components.SearchErrorMsg{
				Error: err.Error(),
				Query: query,
			}
		}
		return components.SearchResultsMsg{
			Results: results,
			Query:   query,
		}
	}
}

// filterConversationsByFilter filters conversations based on the selected filter.
func (m Model) filterConversationsByFilter(filter string) []*domain.Conversation {
	now := time.Now()
	var filtered []*domain.Conversation

	for _, conv := range m.conversations {
		switch filter {
		case "All":
			filtered = append(filtered, conv)
		case "Today":
			// Conversations from today
			if conv.CreatedAt.Year() == now.Year() &&
				conv.CreatedAt.Month() == now.Month() &&
				conv.CreatedAt.Day() == now.Day() {
				filtered = append(filtered, conv)
			}
		case "This Week":
			// Conversations from last 7 days
			if time.Since(conv.CreatedAt) <= 7*24*time.Hour {
				filtered = append(filtered, conv)
			}
		case "This Month":
			// Conversations from last 30 days
			if time.Since(conv.CreatedAt) <= 30*24*time.Hour {
				filtered = append(filtered, conv)
			}
		}
	}

	if len(filtered) == 0 {
		// If no results, return all conversations
		return m.conversations
	}
	return filtered
}

// handleTabKey cycles focus between panes: Sidebar -> Center -> Right -> Sidebar.
func (m Model) handleTabKey() Model {
	// Only allow Tab navigation in list view (not search, detail, etc.)
	if m.state != ViewList {
		return m
	}

	// Toggle between Sidebar and Center only (right pane is display-only)
	if m.focusPane == FocusSidebar {
		m = m.SetFocusPane(FocusCenter)
	} else {
		m = m.SetFocusPane(FocusSidebar)
	}

	// Update component focus states
	if m.sidebar != nil {
		m.sidebar = m.sidebar.SetFocus(m.focusPane == FocusSidebar)
	}

	return m
}
