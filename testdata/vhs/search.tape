Output testdata/vhs/search.gif

# VHS Test: Search Functionality
# Tests: Ctrl+F to enter search, query input, results display, escape to exit

Set FontSize 14
Set Width 1200
Set Height 800
Set Shell "bash"

Type "./bin/claudex"
Enter
Sleep 2s

# Open search with Ctrl+F (correct keybinding, not /)
Ctrl+F
Sleep 500ms
Screenshot testdata/vhs/search-mode.png

# Type search query
Type "test"
Sleep 1s
Screenshot testdata/vhs/search-results.png

# Execute search with Enter
Enter
Sleep 500ms
Screenshot testdata/vhs/search-executed.png

# Clear and exit search with Escape
Escape
Sleep 500ms
Screenshot testdata/vhs/search-cleared.png

Ctrl+C
