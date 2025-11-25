Output testdata/vhs/startup.gif

Set FontSize 14
Set Width 1200
Set Height 800
Set Shell "bash"

Type "./bin/claudex"
Enter
Sleep 2s

Screenshot testdata/vhs/startup.png

# Test navigation
Type@100ms "jjj"
Sleep 500ms
Screenshot testdata/vhs/list-navigation.png

# Open detail view
Enter
Sleep 1s
Screenshot testdata/vhs/detail-view.png

# Return to list
Escape
Sleep 500ms

Ctrl+C
