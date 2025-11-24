// ABOUTME: Sidebar component displays the list of projects with metadata (conversation count,
// last updated) and quick filters (All, Today, This Week, This Month). Supports keyboard
// navigation for project selection and filter switching.
package components

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/TrevorS/claudex/internal/domain"
)

// FilterType represents the type of date filter applied to conversations.
type FilterType string

const (
	FilterAll       FilterType = "All"
	FilterToday     FilterType = "Today"
	FilterThisWeek  FilterType = "This Week"
	FilterThisMonth FilterType = "This Month"
)

// ProjectMeta holds metadata about a project's conversations.
type ProjectMeta struct {
	Count       int       // Number of conversations in this project
	LastUpdated time.Time // Most recent conversation update time
}

// Sidebar displays projects and filters in the left pane.
type Sidebar struct {
	conversations        []*domain.Conversation
	projects             []string               // Unique sorted project paths
	projectMetadata      map[string]ProjectMeta // Project path -> metadata
	selectedProjectIndex int                    // Currently selected project index
	selectedFilter       FilterType             // Currently selected filter
	hasFocus             bool                   // Whether sidebar has keyboard focus
	viewWidth            int
	viewHeight           int
}

// Message types for communication with parent model

// ProjectSelectedMsg is sent when a user selects a project.
type ProjectSelectedMsg struct {
	ProjectPath string
}

// FilterSelectedMsg is sent when a user selects a filter.
type FilterSelectedMsg struct {
	Filter FilterType
}

// FocusNextMsg is sent when the user presses Tab to move focus.
type FocusNextMsg struct{}

// NewSidebar creates a new Sidebar component.
func NewSidebar() *Sidebar {
	return &Sidebar{
		conversations:        make([]*domain.Conversation, 0),
		projects:             make([]string, 0),
		projectMetadata:      make(map[string]ProjectMeta),
		selectedProjectIndex: -1, // -1 means no selection (all projects)
		selectedFilter:       FilterAll,
		hasFocus:             false,
		viewWidth:            0,
		viewHeight:           0,
	}
}

// Init initializes the Sidebar component.
func (s *Sidebar) Init() tea.Cmd {
	return nil
}

// Update handles Bubble Tea messages and returns updated Sidebar.
func (s *Sidebar) Update(msg tea.Msg) (*Sidebar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return s.SetSize(msg.Width, msg.Height), nil

	case tea.KeyMsg:
		if !s.hasFocus {
			// Ignore key events when not focused
			return s, nil
		}
		return s.handleKeyPress(msg)

	default:
		return s, nil
	}
}

// handleKeyPress processes keyboard input when sidebar has focus.
func (s *Sidebar) handleKeyPress(msg tea.KeyMsg) (*Sidebar, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		return s.moveSelection(-1), nil

	case tea.KeyDown:
		return s.moveSelection(1), nil

	case tea.KeyEnter:
		// Send project selected message
		if len(s.projects) > 0 && s.selectedProjectIndex < len(s.projects) {
			project := s.projects[s.selectedProjectIndex]
			return s, func() tea.Msg {
				return ProjectSelectedMsg{ProjectPath: project}
			}
		}
		return s, nil

	case tea.KeyTab:
		// Send focus next message
		return s, func() tea.Msg {
			return FocusNextMsg{}
		}

	case tea.KeyRunes:
		if len(msg.Runes) > 0 {
			switch msg.Runes[0] {
			case 'f', 'F':
				// Cycle through filters
				return s.cycleFilter(), nil

			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// Jump to project by number
				index := int(msg.Runes[0] - '1') // '1' -> index 0
				return s.jumpToProject(index), nil
			}
		}
		return s, nil

	default:
		return s, nil
	}
}

// moveSelection moves the project selection up or down.
func (s *Sidebar) moveSelection(delta int) *Sidebar {
	newS := *s
	newIndex := s.selectedProjectIndex + delta

	// Clamp to valid range
	if newIndex < 0 {
		newIndex = 0
	}
	if newIndex >= len(s.projects) {
		newIndex = len(s.projects) - 1
	}

	newS.selectedProjectIndex = newIndex
	return &newS
}

// jumpToProject jumps to a specific project index (0-based).
func (s *Sidebar) jumpToProject(index int) *Sidebar {
	newS := *s

	// Clamp to valid range
	if index < 0 {
		index = 0
	}
	if index >= len(s.projects) {
		index = len(s.projects) - 1
	}

	newS.selectedProjectIndex = index
	return &newS
}

// cycleFilter cycles through the available filters.
func (s *Sidebar) cycleFilter() *Sidebar {
	newS := *s

	filters := []FilterType{FilterAll, FilterToday, FilterThisWeek, FilterThisMonth}
	currentIndex := 0
	for i, f := range filters {
		if f == s.selectedFilter {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(filters)
	newS.selectedFilter = filters[nextIndex]

	return &newS
}

// View renders the Sidebar as a string.
func (s *Sidebar) View() string {
	if s.viewWidth < 18 {
		// Too narrow to render properly
		return s.renderCollapsed()
	}

	var b strings.Builder

	// Render header
	b.WriteString(s.renderHeader())
	b.WriteString("\n\n")

	// Render projects
	b.WriteString(s.renderProjects())
	b.WriteString("\n\n")

	// Render filters
	b.WriteString(s.renderFilters())
	b.WriteString("\n\n")

	// Render footer hints
	b.WriteString(s.renderFooter())

	// Apply border styling
	style := lipgloss.NewStyle().
		Width(s.viewWidth).
		Height(s.viewHeight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.borderColor()).
		Padding(1)

	return style.Render(b.String())
}

// renderCollapsed renders a minimal view when width is too narrow.
func (s *Sidebar) renderCollapsed() string {
	return lipgloss.NewStyle().
		Width(s.viewWidth).
		Height(s.viewHeight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1).
		Render("...")
}

// renderHeader renders the "Projects" header.
func (s *Sidebar) renderHeader() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6"))

	return style.Render("PROJECTS")
}

// renderProjects renders the list of projects with metadata.
func (s *Sidebar) renderProjects() string {
	if len(s.projects) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("No projects")
	}

	var b strings.Builder

	for i, project := range s.projects {
		isSelected := i == s.selectedProjectIndex && s.hasFocus

		// Format project path
		displayPath := s.formatProjectPath(project)

		// Get metadata
		meta := s.projectMetadata[project]

		// Render project line
		projectLine := s.renderProjectLine(displayPath, meta, isSelected)
		b.WriteString(projectLine)
		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String())
}

// renderProjectLine renders a single project line with metadata.
func (s *Sidebar) renderProjectLine(path string, meta ProjectMeta, isSelected bool) string {
	// Build the display text
	var b strings.Builder

	// Icon
	icon := "📁"
	if path == "" {
		icon = "📄"
		path = "(no project)"
	}

	// Truncate path if too long
	maxPathLen := s.viewWidth - 10
	if maxPathLen < 5 {
		maxPathLen = 5
	}
	if len(path) > maxPathLen {
		path = path[:maxPathLen-3] + "..."
	}

	b.WriteString(fmt.Sprintf("%s %s\n", icon, path))

	// Metadata line (indented)
	metaLine := fmt.Sprintf("   %d conv · %s", meta.Count, s.formatLastUpdated(meta.LastUpdated))
	b.WriteString(metaLine)

	// Apply styling
	style := lipgloss.NewStyle()
	if isSelected {
		style = style.
			Background(lipgloss.Color("6")).
			Foreground(lipgloss.Color("0")).
			Bold(true)
	} else {
		style = style.Foreground(lipgloss.Color("7"))
	}

	return style.Render(b.String())
}

// renderFilters renders the filter selection buttons.
func (s *Sidebar) renderFilters() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6"))
	b.WriteString(headerStyle.Render("FILTERS"))
	b.WriteString("\n")

	// Filter buttons
	filters := []FilterType{FilterAll, FilterToday, FilterThisWeek, FilterThisMonth}
	for _, filter := range filters {
		isSelected := filter == s.selectedFilter
		filterLine := s.renderFilterButton(filter, isSelected)
		b.WriteString(filterLine)
		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String())
}

// renderFilterButton renders a single filter button.
func (s *Sidebar) renderFilterButton(filter FilterType, isSelected bool) string {
	style := lipgloss.NewStyle()

	if isSelected {
		style = style.
			Foreground(lipgloss.Color("2")).
			Bold(true)
	} else {
		style = style.Foreground(lipgloss.Color("240"))
	}

	prefix := "  "
	if isSelected {
		prefix = "▶ "
	}

	return style.Render(fmt.Sprintf("%s%s", prefix, string(filter)))
}

// renderFooter renders keyboard hints at the bottom.
func (s *Sidebar) renderFooter() string {
	if !s.hasFocus {
		return ""
	}

	hints := "↑↓ navigate · ⏎ select · f filter"

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	return style.Render(hints)
}

// borderColor returns the border color based on focus state.
func (s *Sidebar) borderColor() lipgloss.Color {
	if s.hasFocus {
		return lipgloss.Color("6") // Cyan for focused
	}
	return lipgloss.Color("240") // Gray for unfocused
}

// formatProjectPath formats a project path for display.
func (s *Sidebar) formatProjectPath(path string) string {
	if path == "" {
		return "(no project)"
	}

	// Simple display: just show the path
	// Could implement smart truncation here (e.g., /Users/.../project)
	return path
}

// formatLastUpdated formats a timestamp as a relative time string.
func (s *Sidebar) formatLastUpdated(t time.Time) string {
	if t.IsZero() {
		return "never"
	}

	diff := time.Since(t)

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
	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / 24 / 7)
		return fmt.Sprintf("%dw ago", weeks)
	default:
		return t.Format("Jan 2")
	}
}

// SetConversations updates the conversation list and recalculates projects and metadata.
func (s *Sidebar) SetConversations(conversations []*domain.Conversation) *Sidebar {
	newS := *s
	newS.conversations = conversations

	// Extract unique projects
	newS.projects = extractUniqueProjects(conversations)

	// Calculate metadata for each project
	newS.projectMetadata = calculateProjectMetadata(conversations, newS.projects)

	// Clamp selection to valid range (but keep -1 as "all projects")
	if newS.selectedProjectIndex >= len(newS.projects) {
		if len(newS.projects) > 0 {
			newS.selectedProjectIndex = len(newS.projects) - 1
		} else {
			newS.selectedProjectIndex = -1
		}
	}
	// Keep -1 as valid ("all projects") - don't auto-select first project

	return &newS
}

// extractUniqueProjects extracts and sorts unique project paths from conversations.
func extractUniqueProjects(conversations []*domain.Conversation) []string {
	projectSet := make(map[string]bool)

	for _, conv := range conversations {
		projectSet[conv.ProjectPath] = true
	}

	projects := make([]string, 0, len(projectSet))
	for project := range projectSet {
		projects = append(projects, project)
	}

	// Sort projects alphabetically, with empty string last
	sort.Slice(projects, func(i, j int) bool {
		if projects[i] == "" {
			return false
		}
		if projects[j] == "" {
			return true
		}
		return projects[i] < projects[j]
	})

	return projects
}

// calculateProjectMetadata calculates conversation counts and last updated times for projects.
func calculateProjectMetadata(conversations []*domain.Conversation, projects []string) map[string]ProjectMeta {
	metadata := make(map[string]ProjectMeta)

	for _, project := range projects {
		var count int
		var lastUpdated time.Time

		for _, conv := range conversations {
			if conv.ProjectPath == project {
				count++
				if conv.UpdatedAt.After(lastUpdated) {
					lastUpdated = conv.UpdatedAt
				}
			}
		}

		metadata[project] = ProjectMeta{
			Count:       count,
			LastUpdated: lastUpdated,
		}
	}

	return metadata
}

// FilteredConversations returns conversations filtered by the current filter and selected project.
func (s *Sidebar) FilteredConversations() []*domain.Conversation {
	filtered := s.conversations

	// Apply project filter if a specific project is selected (not -1 which means "all")
	if s.selectedProjectIndex >= 0 && s.selectedProjectIndex < len(s.projects) {
		selectedProject := s.projects[s.selectedProjectIndex]
		projectFiltered := make([]*domain.Conversation, 0)
		for _, conv := range filtered {
			if conv.ProjectPath == selectedProject {
				projectFiltered = append(projectFiltered, conv)
			}
		}
		filtered = projectFiltered
	}

	// Apply date filter
	now := time.Now()
	dateFiltered := make([]*domain.Conversation, 0)

	for _, conv := range filtered {
		include := false

		switch s.selectedFilter {
		case FilterAll:
			include = true

		case FilterToday:
			diff := now.Sub(conv.UpdatedAt)
			include = diff < 24*time.Hour

		case FilterThisWeek:
			diff := now.Sub(conv.UpdatedAt)
			include = diff < 7*24*time.Hour

		case FilterThisMonth:
			diff := now.Sub(conv.UpdatedAt)
			include = diff < 30*24*time.Hour
		}

		if include {
			dateFiltered = append(dateFiltered, conv)
		}
	}

	return dateFiltered
}

// SelectedProject returns the currently selected project path.
func (s *Sidebar) SelectedProject() string {
	if len(s.projects) == 0 || s.selectedProjectIndex >= len(s.projects) {
		return ""
	}
	return s.projects[s.selectedProjectIndex]
}

// SelectedFilter returns the currently selected filter.
func (s *Sidebar) SelectedFilter() FilterType {
	return s.selectedFilter
}

// SetFilter sets the current filter.
func (s *Sidebar) SetFilter(filter FilterType) *Sidebar {
	newS := *s
	newS.selectedFilter = filter
	return &newS
}

// Filters returns all available filters.
func (s *Sidebar) Filters() []FilterType {
	return []FilterType{FilterAll, FilterToday, FilterThisWeek, FilterThisMonth}
}

// SetFocus sets whether the sidebar has keyboard focus.
func (s *Sidebar) SetFocus(focus bool) *Sidebar {
	newS := *s
	newS.hasFocus = focus
	return &newS
}

// SetSize sets the sidebar dimensions.
func (s *Sidebar) SetSize(width, height int) *Sidebar {
	newS := *s
	newS.viewWidth = width
	newS.viewHeight = height
	return &newS
}
