Output testdata/vhs/sidebar.gif

# VHS Test: Sidebar and Focus Management
# Tests: Tab focus switching, sidebar navigation, project selection, date filters

Set FontSize 14
Set Width 1200
Set Height 800
Set Shell "bash"

# Launch the application
Type "./bin/claudex"
Enter
Sleep 2s

# Screenshot: Initial state with sidebar focused (cyan border)
Screenshot testdata/vhs/sidebar-focused.png

# Test: Navigate through projects in sidebar with j/k
Type "j"
Sleep 200ms
Type "j"
Sleep 300ms
Screenshot testdata/vhs/sidebar-project-nav.png

# Test: Switch focus to center pane with Tab
Tab
Sleep 300ms
Screenshot testdata/vhs/focus-center.png

# Test: Switch focus back to sidebar with Tab
Tab
Sleep 300ms
Screenshot testdata/vhs/focus-sidebar-again.png

# Test: Date filter cycling with 'f' key
Type "f"
Sleep 500ms
Screenshot testdata/vhs/filter-today.png

Type "f"
Sleep 500ms
Screenshot testdata/vhs/filter-this-week.png

Type "f"
Sleep 500ms
Screenshot testdata/vhs/filter-this-month.png

Type "f"
Sleep 500ms
Screenshot testdata/vhs/filter-all.png

# Test: Select project with Enter
Home
Sleep 200ms
Enter
Sleep 500ms
Screenshot testdata/vhs/project-selected.png

# Test: Clear project selection with Escape
Escape
Sleep 500ms
Screenshot testdata/vhs/project-cleared.png

# Clean exit
Ctrl+C
