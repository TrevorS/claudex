---
name: debug-ui
description: Use Bubble Tea debug logging to inspect message flow, state transitions, and update cycles in terminal UI applications. Activates when debugging UI behavior, investigating message handling, tracking state changes, or diagnosing rendering issues.
---

# Debug UI - Bubble Tea Message Flow Inspection

## Purpose
Enable comprehensive debug logging for Bubble Tea applications to trace message flow, state transitions, and update cycles.

## When to Use
- Debugging unexpected UI behavior or state issues
- Investigating why keyboard input isn't working as expected
- Tracking message propagation through nested models
- Understanding update cycle timing and order
- Diagnosing rendering or layout problems
- Verifying that messages are being sent/received correctly

## Instructions

### Enabling Debug Logging
1. Start dev mode with debug logging enabled:
   ```bash
   make dev-debug
   ```
2. Debug logs automatically write to `debug.log`

### Reading Debug Logs
Open the log file in a separate terminal:
```bash
# Follow logs in real-time
tail -f debug.log

# Search for specific messages
grep "tea.KeyMsg" debug.log

# View recent entries
tail -100 debug.log
```

### What Gets Logged
Bubble Tea automatically logs:
- All `tea.Msg` types received by the root model
- Message type names (e.g., `tea.KeyMsg`, `tea.WindowSizeMsg`)
- Message payloads for keyboard, mouse, and system events
- Timing information for update cycles

### Debug Workflow
1. Run `make dev-debug` in one terminal
2. Open `debug.log` in another terminal with `tail -f debug.log`
3. Interact with the UI
4. Observe message flow in real-time
5. Identify unexpected behavior or missing messages
6. Fix issue and verify with continued logging

## Examples

**Example 1: Debug a keyboard handling issue**
```bash
# Terminal 1: Start dev mode with debug
make dev-debug

# Terminal 2: Watch logs
tail -f debug.log | grep -i key
# Then interact with UI in terminal 1
# See which KeyMsg events are being received
```

**Example 2: Track state transitions**
```bash
# Terminal 1: Start dev mode with debug
make dev-debug

# Terminal 2: Watch all events
tail -f debug.log
# Interact with UI (navigate, select, search, etc.)
# See complete message flow and update order
```

**Example 3: Investigate message timing**
```bash
# Terminal 1: Start dev mode with debug
make dev-debug

# Terminal 2: Filter window resize events
tail -f debug.log | grep WindowSizeMsg
# Resize terminal window in terminal 1
# See how UI responds to size changes
```

## Best Practices
- Always run `make dev-debug` when debugging UI issues
- Use `grep` to filter logs and focus on specific message types
- Combine with Air's auto-rebuild for fast iteration + debugging
- Close log file and stop debugging when issue is resolved
- Debug logs are only active when `DEBUG=1` is set (Makefile handles this)

## Troubleshooting
- **No log file created**: Ensure working directory is project root
- **Log file empty**: Verify `make dev-debug` is actually running (not just `make dev`)
- **Too much output**: Use `grep` to filter specific message types (e.g., `grep KeyMsg`)
- **Logs not updating**: Check that the app is actually running and receiving input
