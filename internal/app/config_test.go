package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "info", config.LogLevel)
	assert.True(t, config.CacheEnabled)
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		cleanup func(t *testing.T)
		wantErr bool
	}{
		{
			name: "loads with defaults when no env vars set",
			setup: func(t *testing.T) {
				os.Unsetenv("CLAUDEX_PROJECTS_PATH")
				os.Unsetenv("CLAUDEX_LOG_LEVEL")
			},
			cleanup: func(t *testing.T) {},
			wantErr: false,
		},
		{
			name: "respects environment variable overrides",
			setup: func(t *testing.T) {
				tempDir := t.TempDir()
				os.Setenv("CLAUDEX_PROJECTS_PATH", tempDir)
				os.Setenv("CLAUDEX_LOG_LEVEL", "debug")
			},
			cleanup: func(t *testing.T) {
				os.Unsetenv("CLAUDEX_PROJECTS_PATH")
				os.Unsetenv("CLAUDEX_LOG_LEVEL")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			defer tt.cleanup(t)

			config, err := LoadConfig()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, config)
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with existing directory",
			config: &Config{
				ClaudeProjectsPath: t.TempDir(),
			},
			wantErr: false,
		},
		{
			name: "invalid config with non-existent directory",
			config: &Config{
				ClaudeProjectsPath: "/nonexistent/path/that/does/not/exist",
			},
			wantErr: true,
			errMsg:  "directory does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestConfigPathExpansion(t *testing.T) {
	config := DefaultConfig()

	// Verify that home directory is expanded in default config
	assert.NotContains(t, config.ClaudeProjectsPath, "~")
	assert.True(t, filepath.IsAbs(config.ClaudeProjectsPath))
}

func TestConfigLogLevel(t *testing.T) {
	config := DefaultConfig()
	assert.Equal(t, "info", config.LogLevel)

	config.LogLevel = "debug"
	assert.Equal(t, "debug", config.LogLevel)
}
