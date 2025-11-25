// ABOUTME: Main application bootstrap coordinating repository, UI, and Bubble Tea
package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/TrevorS/claudex/internal/repository"
	"github.com/TrevorS/claudex/internal/ui"
)

// App is the main application coordinator.
type App struct {
	config     *Config
	repository repository.Repository
	program    *tea.Program
}

// New creates and initializes a new App with the given configuration.
// It validates the config and initializes the repository with caching.
func New(config *Config) (*App, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create base JSONL repository
	repo, err := repository.NewJSONLRepository(config.ClaudeProjectsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Wrap with caching layer (max 50 conversations in LRU cache)
	cachedRepo, err := repository.NewCachedRepository(repo, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}

	app := &App{
		config:     config,
		repository: cachedRepo,
	}

	return app, nil
}

// Run starts the Bubble Tea program and runs the application.
// This is a blocking call that returns when the user quits the application.
func (a *App) Run() error {
	// Create the root UI model
	model := ui.NewModel(a.repository)

	// Create and run the Bubble Tea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	a.program = p

	// Run the program
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("application error: %w", err)
	}

	return nil
}

// Shutdown performs graceful application shutdown.
func (a *App) Shutdown() error {
	if a.program != nil {
		a.program.Quit()
	}
	return nil
}

// Config returns the application configuration.
func (a *App) Config() *Config {
	return a.config
}

// Repository returns the repository instance.
func (a *App) Repository() repository.Repository {
	return a.repository
}
