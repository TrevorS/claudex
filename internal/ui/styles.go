// ABOUTME: Reusable Lipgloss style definitions for UI components
package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// BorderStyle is the base border style for panels.
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("59"))

	// FocusedBorderStyle is the border style for focused panels.
	FocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")) // Purple

	// HeaderStyle is the style for header text.
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63")). // Purple
			Padding(0, 1)

	// ErrorStyle is the style for error messages.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")). // Red
			Bold(true)

	// SuccessStyle is the style for success messages.
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")). // Green
			Bold(true)

	// InfoStyle is the style for info messages.
	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("4")). // Blue
			Italic(true)

	// SubtleStyle is the style for subtle text.
	SubtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")) // Bright black/gray
)
