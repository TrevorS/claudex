---
name: dev-mode
description: Use Air for live reload during Bubble Tea UI development in Go projects. Activates when developing terminal UIs, working on view logic, iterating on layout changes, or needing instant feedback on UI modifications. Watches Go files and rebuilds on changes.
---

# Dev Mode - Live Reload with Air

## Purpose
Enable rapid UI iteration using Air's live reload functionality for Bubble Tea applications.

## When to Use
- Developing or modifying Bubble Tea UI components
- Working on view rendering logic in `internal/ui/`
- Iterating on layout, styling, or component composition
- Testing UI behavior changes without manual rebuild cycles
- Any work requiring instant visual feedback

## Instructions

### Starting Dev Mode
1. Ensure Air is installed: `which air` (install with `go install github.com/air-verse/air@latest` if missing)
2. Check for `.air.toml` configuration in project root
3. Start Air with Makefile target:
   ```bash
   make dev
   # or with debug logging:
   make dev-debug
   ```
4. Air will:
   - Build the application
   - Start it automatically
   - Watch for file changes
   - Rebuild and restart on save

### During Development
- Make changes to any `.go` files in the project
- Save the file
- Air automatically rebuilds and restarts within 1-2 seconds
- View changes immediately in the terminal

### Configuration
The `.air.toml` file controls:
- **Include patterns**: Which files trigger rebuilds (default: `**/*.go`)
- **Exclude patterns**: Ignored paths (e.g., `testdata/`, `vendor/`)
- **Build command**: How to compile (builds `./cmd/claudex`)
- **Clear screen**: Clears terminal on each rebuild for clean view

### Stopping Dev Mode
- Press `Ctrl+C` to stop Air
- Cleans up temporary build artifacts automatically

## Examples

**Example 1: Starting dev mode for UI work**
```bash
make dev
```

**Example 2: Modifying view logic with debug logging**
```bash
make dev-debug
# Edit internal/ui/view.go in another terminal
# Save file → Air detects change → Auto rebuild → UI refreshes
# Debug logs appear in debug.log (watch with: tail -f debug.log)
```

**Example 3: Working with components**
```bash
# Start dev mode
make dev

# In editor, modify internal/ui/components/sidebar.go
# Save file → automatically rebuilt and running in <2s
```

## Troubleshooting
- **Air not found**: Run `go install github.com/air-verse/air@latest` and ensure `~/go/bin` is in PATH
- **Changes not detected**: Check `.air.toml` include/exclude patterns
- **Build errors**: Air shows compiler errors in real-time; fix and save to retry
- **Port conflicts**: Ensure no other instance of the app is running on the same port
