Output testdata/vhs/startup.gif

# VHS Test: Application Startup
# Tests: Initial launch, three-pane layout, header/footer rendering

Set FontSize 14
Set Width 1200
Set Height 800
Set Shell "bash"

# Launch the application
Type "./bin/claudex"
Enter
Sleep 2s

# Screenshot 1: Initial startup with three-pane layout
# Expected: Header "claudex • conversation browser", sidebar with projects,
#           conversation list in center, metadata panel on right
Screenshot testdata/vhs/startup.png

# Test: Verify initial focus is on sidebar (default behavior)
# The sidebar should have a cyan border indicating focus
Sleep 500ms
Screenshot testdata/vhs/startup-initial-focus.png

# Clean exit
Ctrl+C
