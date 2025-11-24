// ABOUTME: Configuration management for claudex application settings and paths
package app

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds application configuration including paths and settings.
type Config struct {
	ClaudeProjectsPath string // Path to ~/.claude/projects/ directory
	CachePath          string // Path to cache directory
	CacheEnabled       bool   // Whether to use caching
	LogLevel           string // Logging level (debug, info, warn, error)
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	projectsPath := expandPath("~/.claude/projects")
	cachePath := expandPath("~/.cache/claudex")

	return &Config{
		ClaudeProjectsPath: projectsPath,
		CachePath:          cachePath,
		CacheEnabled:       true,
		LogLevel:           "info",
	}
}

// LoadConfig loads configuration from defaults and environment variables.
// Environment variables:
//   - CLAUDEX_PROJECTS_PATH: Override projects directory path
//   - CLAUDEX_LOG_LEVEL: Override logging level
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Override with environment variables
	if projectsPath := os.Getenv("CLAUDEX_PROJECTS_PATH"); projectsPath != "" {
		config.ClaudeProjectsPath = expandPath(projectsPath)
	}

	if logLevel := os.Getenv("CLAUDEX_LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	return config, nil
}

// Validate checks that the configuration is valid.
func (c *Config) Validate() error {
	if c.ClaudeProjectsPath == "" {
		return fmt.Errorf("projects path must not be empty")
	}

	// Check if projects directory exists
	if _, err := os.Stat(c.ClaudeProjectsPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("projects directory does not exist: %s", c.ClaudeProjectsPath)
		}
		return fmt.Errorf("failed to validate projects directory: %w", err)
	}

	return nil
}

// expandPath expands ~ to the user's home directory.
func expandPath(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		// If we can't get home directory, return the path as-is
		return path
	}

	return filepath.Join(home, path[1:])
}
