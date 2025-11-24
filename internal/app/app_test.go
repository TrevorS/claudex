package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TrevorS/claudex/internal/repository"
)

func TestNewApp(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{
		ClaudeProjectsPath: tempDir,
		CacheEnabled:       false,
		LogLevel:           "debug",
	}

	app, err := New(config)
	require.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, config, app.config)
	assert.NotNil(t, app.repository)
}

func TestNewAppInvalidConfig(t *testing.T) {
	config := &Config{
		ClaudeProjectsPath: "/nonexistent/path",
	}

	app, err := New(config)
	assert.Error(t, err)
	assert.Nil(t, app)
}

func TestAppConfig(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{
		ClaudeProjectsPath: tempDir,
		LogLevel:           "debug",
	}

	app, err := New(config)
	require.NoError(t, err)
	assert.Equal(t, config.LogLevel, app.config.LogLevel)
	assert.Equal(t, config.ClaudeProjectsPath, app.config.ClaudeProjectsPath)
}

func TestAppRepository(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{
		ClaudeProjectsPath: tempDir,
	}

	app, err := New(config)
	require.NoError(t, err)

	// Verify repository is initialized and implements Repository interface
	var _ repository.Repository = app.repository
}
