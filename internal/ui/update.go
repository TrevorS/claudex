// ABOUTME: Centralized message routing and state transitions
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all messages and returns a new model and command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// Handle terminal resize
		return m.handleWindowSize(msg), nil

	case tea.KeyMsg:
		// Handle keyboard input
		return m.handleKeyPress(msg)

	default:
		return m, nil
	}
}

// handleWindowSize updates the model with new terminal dimensions.
func (m Model) handleWindowSize(msg tea.WindowSizeMsg) Model {
	return m.SetSize(msg.Width, msg.Height)
}

// handleKeyPress handles global keyboard shortcuts.
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {

	case tea.KeyCtrlC:
		// Quit application
		return m, tea.Quit

	case tea.KeyEsc:
		// Return to list view or close modal
		if m.state == ViewDetail || m.state == ViewSearch || m.state == ViewStats {
			return m.SetState(ViewList), nil
		}
		if m.state == ViewPalette {
			return m.SetState(ViewList), nil
		}
		return m, nil

	case tea.KeyTab:
		// Tab key for focus switching (no-op in Phase 1)
		return m, nil

	case tea.KeyCtrlP:
		// Command palette (Phase 2+)
		return m, nil

	case tea.KeyEnter:
		// Handle enter based on current state (no-op in Phase 1)
		return m, nil

	default:
		return m, nil
	}
}
