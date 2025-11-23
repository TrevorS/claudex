Claudex: Claude Code Conversation History Viewer
Complete UX Specification Document v1.0
Executive Summary
Claudex is a terminal user interface application that enables developers to browse, search, and analyze their Claude Code conversation history with the same sophistication they expect from modern GUI applications. The application transforms the raw JSONL conversation data stored in ~/.claude/ into an intuitive, navigable interface featuring three-pane layouts, real-time search, statistical dashboards, and rich message rendering.

This specification defines every aspect of the user experience from initial application launch through advanced filtering workflows. Engineers implementing this specification will find exact layout dimensions, complete component behaviors, state transition logic, and visual design tokens necessary to build claudex without requiring additional design decisions.

Application Architecture and Navigation Model
Root Application States
The application operates in five primary states that define the user's current context and available interactions:

List View State serves as the default entry point where users see the three-pane layout with sidebar, conversation list, and empty detail pane. This state allows users to browse projects, apply filters, and select conversations for detailed viewing. The focus management system begins with the conversation list pane focused by default.

Detail View State activates when users select a conversation, shifting focus to the message viewer pane while maintaining sidebar and list visibility. The detail pane expands to show full conversation content with scrollable messages. Users can navigate back to list view or select different conversations without losing their current filter context.

Search Active State overlays a search interface on the current view, dimming the background content and presenting a centered search input with live-filtered results. The search scope respects the current sidebar filter selections, allowing users to narrow searches to specific projects or date ranges.

Statistics Dashboard State transitions the entire interface to a full-screen dashboard displaying aggregated metrics, charts, and usage patterns. This state provides a dedicated analysis environment while maintaining quick access back to conversation browsing through keyboard shortcuts.

Command Palette State presents a modal overlay for command-driven navigation and actions. The palette activates via keyboard shortcut and provides fuzzy-matched access to all application functionality, serving as a power-user alternative to mouse-based navigation.

State Transition Rules
Transitions between states follow these patterns to maintain predictable navigation flow:

From List View, users can enter Detail View by selecting a conversation with Enter or mouse click, activate Search with forward slash or Ctrl+F, open Statistics Dashboard with S key, or summon Command Palette with Ctrl+P. The Escape key in List View focuses the sidebar if not already focused, otherwise it does nothing to prevent accidental application exit.

From Detail View, users return to List View by pressing Escape or B for back, activate Search to filter within the current conversation, open Statistics showing metrics for the current conversation, or use Command Palette for quick actions. The back navigation preserves scroll position in the conversation list so users return to their previous browsing context.

From Search Active, pressing Escape dismisses the search overlay and returns to the previous state while preserving search terms for quick re-activation. Executing a search by pressing Enter applies filters and returns to List View with updated results, maintaining the search term in the UI for modification.

From Statistics Dashboard, users press Escape or Q to return to the previous browsing state, press S again to toggle back, or use Command Palette to jump directly to other contexts. The dashboard remembers which conversation or project was active before opening statistics.

From Command Palette, pressing Escape dismisses the palette without action, executing a command with Enter performs the action and dismisses the palette, and typing refines the command list through fuzzy matching. Commands that require parameters transition to specialized input prompts before execution.

Initialization Sequence
Application startup follows a four-stage loading process that provides progressive feedback to users:

Stage One: Discovery scans the ~/.claude/projects/ directory to identify all project folders, decoding the hyphenated path names into readable project paths. During this stage, users see a centered loading indicator with the message "Discovering Claude Code projects..." and a count of projects found. This stage typically completes in under 100 milliseconds for most systems.

Stage Two: Metadata Loading reads the first and last entries from each conversation JSONL file to extract session IDs, timestamps, and message counts without parsing complete files. The loading message updates to "Loading conversation metadata... X of Y projects" with a progress bar indicating completion percentage. This stage prioritizes recently modified conversations to show relevant content quickly.

Stage Three: UI Initialization constructs the Bubble Tea component hierarchy, initializes the conversation list table with discovered conversations sorted by last activity, and prepares the viewport for message rendering. Users see the application layout appear with skeleton loading indicators in the conversation list showing placeholder rows.

Stage Four: Content Population replaces skeleton rows with actual conversation data, calculates aggregate statistics for the sidebar, and marks the application ready for interaction. The loading indicators fade and the interface becomes fully responsive, with keyboard focus automatically set to the conversation list for immediate browsing.

If any stage encounters errors such as permission issues or corrupted data files, the application displays a detailed error modal explaining the issue and offering recovery options including retry, skip corrupted files, or exit. The error modal includes collapsible technical details for troubleshooting.

Three-Pane Layout Design and Responsive Behavior
Layout Composition and Dimensions
The application employs a fixed three-pane horizontal layout optimized for terminals 120 characters wide or greater. The layout adapts gracefully to various terminal sizes while maintaining usability.

Sidebar Pane occupies the leftmost 22 characters (approximately 18 percent of a 120-character terminal). This width accommodates typical project path names with moderate truncation while providing sufficient space for filter labels and metadata. The minimum viable width is 18 characters, below which the sidebar must collapse.

Conversation List Pane occupies 60 characters (50 percent of width) providing room for the timestamp column (12 characters), title column (flexible, minimum 25 characters), message count column (8 characters), token count column (10 characters), and cost column (8 characters) with appropriate spacing. This pane scales proportionally with terminal width increases but maintains minimum widths to prevent cramped displays.

Detail Pane occupies the remaining 38 characters (approximately 32 percent), providing adequate width for message content with 80-character line wrapping for code blocks and prose. This pane expands to fill available space when the conversation list collapses or narrows in wide terminals.

Vertical space divides into header (1 row), main content area (terminal height minus 3 rows), and footer/status bar (2 rows including help hints). Each pane manages its own vertical scrolling independently through viewport components.

Responsive Breakpoint Behaviors
Terminal width determines layout adaptation through three breakpoints:

Wide Layout (120+ characters) displays all three panes at their optimal proportions with full feature visibility including all conversation list columns, complete sidebar sections, and comfortably wrapped message content.

Standard Layout (100-119 characters) maintains three panes but reduces conversation list columns to timestamp, title, and cost only, hides extended filter groups in the sidebar requiring expansion to view, and tightens padding throughout the interface.

Narrow Layout (80-99 characters) collapses the sidebar by default, showing only when explicitly toggled with a keyboard shortcut, displays conversation list with timestamp and truncated title only, and expands detail pane to use the full width when active.

Minimal Layout (below 80 characters) switches to a single-pane mode where users toggle between conversation list and detail view, sidebar becomes a full-screen overlay when activated, and the interface displays a persistent banner suggesting terminal resize for optimal experience.

Responsive behavior updates dynamically when terminal size changes during active use, with pane content reflowing smoothly and focus remaining on the same logical element even as layout shifts.

Pane Separation and Focus Indicators
Border styling clearly delineates pane boundaries while indicating focus state through color variation:

Unfocused panes display borders using Lipgloss RoundedBorder style in the color #374151 (neutral gray), creating subtle separation without visual distraction. Pane headers render in #6B7280 (lighter gray) for metadata and status text.

The focused pane's border changes to #7D56F4 (purple accent) with increased contrast, making the active context immediately apparent. The focused pane's header text renders in matching purple, and the header background highlights in #2D1B4E (deep purple) for additional emphasis.

Focus transitions animate using Harmonica spring physics with angular frequency 6.0 and damping ratio 1.0, creating smooth border color interpolation over approximately 200 milliseconds. This subtle animation helps users track focus movement without creating distracting motion.

Corner treatments use the rounded border characters ╭─╮ for top and ╰─╯ for bottom, providing modern aesthetic consistent with contemporary terminal applications. Vertical separators between panes use the thin line character │ to minimize visual weight.

Pane Headers and Status Displays
Each pane includes a header row providing context and status information:

The Sidebar Header displays "Projects & Filters" in bold white text on purple background (#7D56F4), right-aligned with a small icon indicating expand/collapse state. When sidebar is focused, the header additionally shows the count of visible projects and active filters.

The Conversation List Header shows "Conversations" on the left with active filter tags displayed beneath in small gray text with dismissible indicators. The right side displays sort indicator (ascending/descending icon plus column name) and result count such as "247 conversations" or "3 of 247" when filters active. When searching, the header transitions to show "Search Results" with match count.

The Detail Pane Header displays the selected conversation's title in bold with timestamp beneath in smaller gray text. The right side shows token summary "45.2k tokens · $0.67" and a conversation menu icon. When no conversation is selected, the header shows "Select a conversation" in gray italic text.

Footer bars span the full application width beneath all panes, divided into three sections aligned with the panes above. The left section shows sidebar shortcuts, center shows conversation list shortcuts, right shows detail pane shortcuts. The footer adapts to show context-sensitive hints based on which pane has focus and current application state.

Sidebar Component Specification
Overall Structure and Organization
The sidebar organizes content into three primary sections stacked vertically with clear visual separation:

The Global Header occupies the top 3 rows including the "Projects & Filters" title row (styled with purple background as specified earlier), a thin separator line, and the application status indicator showing "Ready" in green, "Loading..." in yellow with spinner, or "Error" in red when appropriate.

The Projects Section follows immediately below, consuming approximately 60 percent of available vertical space. The section begins with a "PROJECTS" label in small gray uppercase text, then lists project items with 1-row spacing between entries. Projects display in a scrollable region when count exceeds available space.

The Filters Section occupies the remaining 40 percent of vertical space below projects. It begins with a "FILTERS" label matching projects section styling, then displays filter groups in a collapsible hierarchy. Common filters like "Today," "This Week," and "High Token Usage" appear as top-level items, while advanced filters nest under expandable groups.

Scroll indicators appear as small arrows (↑↓) in the section headers when content exceeds visible space, with smooth scrolling animation when users navigate beyond viewport boundaries.

Project Item Template and Display
Each project item renders across 3 rows providing comprehensive information at a glance:

Row 1: Project Path displays with the folder icon emoji 📁 followed by the decoded project path truncated intelligently. The truncation algorithm removes intermediate path segments while preserving the first and last segments, showing ellipsis for omitted parts. For example /Users/developer/code/my-awesome-project becomes /Users/.../my-awesome-project when space is limited. The full path appears in a tooltip on hover or in the footer when keyboard-focused.

Row 2: Metadata shows conversation counts and recent activity in small gray text indented by 3 spaces to align beneath the icon. The format reads "24 conversations · 12 this week" with the counts automatically updating as conversations load. The "this week" count calculates from conversations with timestamps within the last 7 days.

Row 3: Spacing provides visual separation from the next project with an empty row.

Selection state modifies the project item's appearance dramatically. The selected project displays with purple background (#2D1B4E) across all 3 rows, purple text (#7D56F4) for the path name, and white text for metadata. A vertical bar indicator (│) appears in the left margin at the selected project's first row in matching purple.

When sidebar has keyboard focus, the currently focused project (not necessarily selected) displays with a subtle gray background (#1F2937) to distinguish from the selected project's more prominent purple highlighting. This two-level visual hierarchy allows users to navigate through projects without changing the active conversation list.

Project Discovery and Sorting
Projects populate the sidebar through automatic discovery of the ~/.claude/projects/ directory, with each encoded directory name decoded back to its original path by replacing hyphens with forward slashes.

The default sort order presents projects by last activity timestamp descending, ensuring recently active projects appear at the top. Users can change sort order through the command palette or a dedicated sidebar menu, with options including:

Alphabetical by project name ascending or descending, total conversation count descending or ascending, last activity timestamp (default), total token usage descending, and total cost descending. The current sort indicator displays subtly in the Projects section header.

Projects with zero conversations do not appear in the list by default but can be shown through a "Show empty projects" toggle in the sidebar settings menu. This prevents clutter while allowing users to verify complete project discovery when needed.

Filter Section Organization and Behavior
The Filters section employs a hierarchical structure supporting both quick access filters and advanced filter composition:

Quick Filters appear as top-level items with simple labels and single-click activation. These include temporal filters (Today, Yesterday, This Week, This Month), usage filters (High Cost >$1, High Tokens >50k, Many Messages >20), and status filters (Has Errors, Uses Tool X). Each quick filter displays with a circle icon (○) when inactive and filled circle (●) when active.

Advanced Filter Groups appear below quick filters with expandable triangles (▶ collapsed, ▼ expanded). Groups organize related filters such as "Date Range" for custom start/end date selection, "Token Range" for minimum/maximum token thresholds, "Tool Usage" for filtering by specific tools used, "Cost Range" for budget-conscious filtering, and "Project Scope" for limiting to specific projects or excluding others.

Filter activation follows additive logic within groups (OR) and intersective logic between groups (AND). For example, selecting "Uses Bash" OR "Uses Read" within Tool Usage, AND "This Week" within Date Range shows conversations from this week that used either Bash or Read tools.

Active filter count displays in the Filters section header as "FILTERS (3 active)" in purple text when any filters apply. A "Clear All" action appears at the bottom of the filter list when filters are active, styled as a small button with the text in red to indicate destructive action.

Keyboard Navigation Within Sidebar
Sidebar keyboard controls enable rapid navigation without mouse:

Up/Down Arrow Keys (or Vim j/k) move focus between items in the current section, scrolling content when focus moves beyond visible area. Navigation wraps at list boundaries, so Down from the last project moves to the first filter and vice versa.

Right Arrow (or l) expands collapsed filter groups when focused on a group header. When focused on an already expanded group, Right moves focus to the first filter within that group.

Left Arrow (or h) collapses expanded filter groups when focused on a group header. When focused on a filter within a group, Left moves focus back to the group header.

Enter/Space activates the focused item. For projects, this updates the conversation list to show only that project's conversations. For filters, this toggles the filter on or off. For filter groups, this expands or collapses the group.

Tab moves focus from the sidebar to the conversation list pane. Shift+Tab moves focus back from conversation list to sidebar, returning to the previously focused sidebar item.

Numbers 1-9 serve as quick shortcuts jumping directly to projects 1 through 9 in the current list order. This enables instant project switching for users with regular projects.

Forward Slash (/) opens search mode even when sidebar is focused, maintaining the current project/filter context for the search scope.

The focused item within the sidebar displays with the subtle gray background mentioned earlier, with a > cursor character appearing in the left margin to reinforce focus position for keyboard-only users.

Conversation List Component Specification
Table Structure and Column Definitions
The conversation list implements a sophisticated table using Evertras/bubble-table with five columns providing comprehensive conversation metadata:

Timestamp Column (12 characters fixed width) displays relative timestamps for recent conversations ("2 hours ago," "Yesterday," "3 days ago") and absolute dates for older conversations ("Jan 15, 2025"). The column header reads "Last Active" and includes a sort indicator showing current sort direction. Timestamps update dynamically during active use to maintain accuracy.

Title Column (flexible width, minimum 25 characters) shows conversation summaries extracted from the JSONL summary messages when available, or derives intelligent titles from the first user message when summaries are absent. The derivation algorithm prioritizes messages containing questions, removes common prefixes like "Claude," and truncates intelligently at word boundaries. Full titles appear in tooltips on hover or in the detail pane when selected.

Messages Column (8 characters fixed width) displays message count as a simple integer with thousands separator for high counts ("1,234"). The column header "Msgs" keeps the table compact while remaining clear.

Tokens Column (10 characters fixed width) shows total tokens formatted with K/M suffixes for readability ("45.2k," "1.2M"). The column includes a small icon indicator for conversations with high cache token usage (showing a ⚡ lightning bolt when cache reads exceed 50k tokens). Header reads "Tokens" with sort indicator.

Cost Column (8 characters fixed width) displays estimated cost in USD with two decimal precision ("$0.67," "$12.45"). Costs above $10 render in yellow text to highlight expensive conversations. Header reads "Cost" with sort indicator.

Column widths adapt in narrow layouts by hiding Messages and Tokens columns first, then Cost column, leaving only Timestamp and Title as the minimum viable display. When columns hide, their information becomes available in the detail pane or through tooltip overlays.

Conversation Item Rendering and Visual Hierarchy
Each conversation row renders with subtle visual cues indicating status and selection:

Unselected rows display with dark background (#111827), light gray text (#D1D5DB) for content, and medium gray (#6B7280) for metadata columns. Row borders use the thin line character ─ in dark gray (#374151) to separate items without creating heavy visual weight.

Selected row highlights with purple-tinted background (#2D1B4E), purple text (#7D56F4) for the title, and white text for metadata columns. A purple vertical bar (│) appears in the left margin to reinforce selection. The selected row's border intensifies to match the purple accent color.

Hovered rows (when mouse is active) display a slightly lighter background than unselected (#1F2937) without changing text colors, providing subtle feedback for mouse-based navigation.

Focus indicator for keyboard navigation shows a more prominent version of the hover background with an additional cursor marker > in the left margin. This allows users to distinguish between the selected conversation (which determines detail pane content) and the focused conversation (which receives keyboard input).

Row height remains constant at 1 line to maximize list density, with all information fitting into the defined column structure. Long titles truncate with ellipsis, and users access full content by selecting the conversation to view in the detail pane.

Selection Model and Multi-Select Behavior
The conversation list supports both single-selection and multi-selection modes:

Single-Select Mode (default) maintains one selected conversation at a time, updating the detail pane immediately when selection changes. Clicking or pressing Enter on a conversation selects it and moves keyboard focus. Pressing Up/Down arrows moves focus and selection together, providing fluid browsing.

Multi-Select Mode activates when users press Space on a focused conversation without deselecting previous selections, or when holding Shift while clicking. In this mode, Space toggles selection of the focused item without changing focus, Shift+Up/Down extends selection to include all conversations between the anchor point and current focus, and Ctrl+A selects all visible conversations matching current filters.

When multiple conversations are selected, the detail pane displays a summary panel showing aggregate statistics (total messages, total tokens, total cost, date range) rather than individual conversation content. A header message "3 conversations selected" with individual deselect icons provides clarity.

Multi-select enables batch operations accessed through the command palette or context menu including bulk export, batch tagging (future feature), and deletion. The interface displays confirmation dialogs before executing destructive operations on multiple conversations.

Pagination Strategy and Performance
The conversation list implements virtual scrolling for optimal performance with large conversation counts:

Visible Window renders approximately 30 conversation rows at once (adjusting based on terminal height), with buffer rows above and below the viewport to enable smooth scrolling. This windowing approach keeps memory usage constant regardless of total conversation count.

Scroll Behavior responds to Page Up/Down keys jumping by full page height, Home/End keys jumping to list start/end, and arrow keys moving by single items. Mouse wheel support scrolls by 3 items per wheel click. Scroll position indicators appear in the column headers showing "1-30 of 247" or similar ranges.

Lazy Loading defers full conversation parsing until needed, loading only metadata during list population. When users select a conversation, the application parses the complete JSONL file asynchronously, displaying a loading spinner in the detail pane while parsing occurs. Parsed conversations cache in memory for instant recall when users navigate back.

Progressive Rendering populates the list in batches during initial load, showing skeleton rows that gradually fill with real data. The first batch of 50 conversations (sorted by most recent) renders within 200 milliseconds, providing immediate responsiveness even on slower systems.

Search Integration and Filter Display
Search and filter states integrate directly into the conversation list interface:

Search Bar appears above the conversation list when users activate search mode with / or Ctrl+F. The bar displays a text input field spanning the full pane width with placeholder text "Search conversations..." and a subtle purple underline indicating focus. Search executes on each keystroke with 300ms debounce, updating the list with matching conversations.

Active Filter Tags render below the search bar (or directly below the column headers when search is inactive) as small pill-shaped badges with gray backgrounds. Each tag shows the filter name (e.g., "This Week," "Uses Bash," "$0.50-$2.00") with a small × close button for quick removal. When multiple filters are active, tags wrap to multiple rows if needed.

Result Count Indicator displays prominently in the list header as "Showing 23 of 247 conversations" when filters or search are active, or "247 conversations" when showing all. The indicator updates in real-time as users type search terms or toggle filters.

Clear Filters Action appears as a small "Clear all filters" link in purple text at the right edge of the filter tags area when any filters are active. Clicking or pressing Enter on this link removes all filters and search terms, restoring the full conversation list.

Empty state for filtered results shows a centered message "No conversations match your filters" with a list of active filters beneath and a large "Clear filters" button. This prevents confusion when legitimate filter combinations yield no results.

Sorting Capabilities and Indicators
Users control conversation list sorting through column header interactions:

Sort Direction toggles between ascending and descending when users click a column header or press Enter while focused on the header. The sort indicator updates to show current direction using ↑ for ascending and ↓ for descending.

Multi-Column Sort supports sort hierarchies where users can specify primary, secondary, and tertiary sort columns. Shift+Click on headers adds columns to the sort hierarchy, with small numbered indicators (1, 2, 3) showing sort order. For example, sorting by "Cost ↓ (1)" then "Tokens ↓ (2)" groups expensive conversations with highest token counts appearing first.

Default Sort orders conversations by Timestamp descending (most recent first) with a secondary sort by Cost descending to break ties. This ensures the most relevant and most expensive recent conversations appear at the top.

Sort Persistence saves the current sort configuration per project in application settings, restoring the preferred sort order when users switch between projects. A "Reset to default sort" command in the command palette provides quick access to standard sorting.

Column headers indicate sortable status with a subtle hover effect and arrow icons. The currently sorted column's header displays in bold purple text with the sort direction indicator prominent.

Detail Pane Message Viewer Specification
Conversation Header Design
The detail pane begins with a comprehensive header providing conversation context:

Title Row displays the conversation title or summary in large bold white text (#FFFFFF) spanning the full pane width. When title exceeds available width, it wraps to a second line rather than truncating, ensuring full context visibility.

Metadata Row shows timestamp on the left in medium gray (#9CA3AF) using full ISO format ("January 15, 2025 at 2:34 PM") for precision. The right side displays aggregate token usage and cost in the format "45.2k tokens · $0.67" with the same gray styling. When cache tokens are present, an additional indicator shows "12k cached" in lighter gray with the lightning bolt emoji.

Summary Statistics Bar provides quick metrics in a third row including total message count, conversation duration (calculated from first to last message timestamp), and participant count (typically just "user + assistant" but can include system messages). Icons prefix each metric using 💬 for messages, ⏱️ for duration, and 👥 for participants.

Action Icons Row aligns to the right with small icon buttons for conversation operations including export (📤), copy link (🔗), refresh (🔄), and more options menu (⋮). These icons respond to hover with purple tinting and show tooltips explaining their function.

A thin horizontal divider (─ repeated across pane width in dark gray) separates the header from the message list below, providing clear visual distinction between metadata and content.

Message List Layout and Grouping
Messages render chronologically in a vertically scrollable region using bubbles/viewport for smooth performance:

Message Grouping clusters messages by time proximity, adding timestamp dividers between groups when gaps exceed 1 hour. Dividers render as centered gray text showing relative or absolute time ("3 hours later," "Next day: January 16") with thin horizontal lines extending left and right.

Message Spacing provides 2 blank rows between individual messages within a group and 3 blank rows between groups (including the divider row). This spacing prevents visual crowding while maintaining readable density.

Thread Indicators show parent-child relationships for tool use and tool result sequences through subtle indentation and connecting lines. When an assistant message includes tool use requests, the subsequent tool result messages indent by 2 characters with a vertical line (│) in gray connecting them to their parent message.

The viewport scrolls smoothly using keyboard controls (Up/Down arrows, Page Up/Down, Home/End, j/k) and mouse wheel. Scroll position indicators appear as small percentage markers in the top-right corner of the viewport showing position like "↑ 34% ↓" when not at list boundaries.

User Message Component Template
User messages render with a distinct visual treatment emphasizing the conversational nature:

Message Container uses a rounded border box (╭─╮ top, ╰─╯ bottom) with blue-tinted background (#1E3A5F) and blue accent border (#60A5FA). The container spans the full pane width minus 2 characters padding on each side.

Header Row displays the user icon emoji (👤) followed by "user" in bold blue text (#60A5FA), followed by a dot separator and timestamp in gray. The timestamp shows relative time for recent messages ("2 minutes ago") and absolute time for older messages ("Jan 15 at 2:34 PM").

Content Section renders the user's message text in light gray (#D1D5DB) with markdown formatting support including bold, italic, inline code, and code blocks. Line breaks preserve original formatting, and URLs render as clickable links (underlined in blue) when terminal supports hyperlinks.

Content Width respects a maximum of 80 characters per line within the message box, wrapping longer lines at word boundaries. Code blocks receive special treatment with monospace formatting, syntax highlighting when language is detectable, and no line wrapping to preserve code structure.

Padding within the message container provides 1 character space on all sides, creating comfortable reading space without excessive whitespace.

Assistant Message Component Template
Assistant messages employ a distinct green-tinted design emphasizing AI responses:

Message Container uses the same rounded border style as user messages but with green-tinted background (#0D3B2F) and green accent border (#04B575).

Header Row displays the robot icon emoji (🤖) followed by "assistant" in bold green text (#04B575), timestamp, and token usage metrics. The token display shows "4.2k in · 1.8k out" in small gray text, providing immediate visibility into response cost. When thinking content is present, a small "💭" icon appears indicating collapsible reasoning.

Content Sections separate the assistant's response into distinct blocks:

The Thinking Section (when present) renders in a collapsible panel with lighter background (#1F2937) and gray text (#9CA3AF). The section header shows "💭 Thinking" with a collapse/expand indicator (▶/▼). When expanded, the thinking content displays in italic gray text providing insight into Claude's reasoning process. This section collapses by default to reduce visual noise.

The Text Content Section follows thinking with the main response text in light gray (#D1D5DB) with full markdown rendering support including headers, lists, quotes, tables, and code blocks. Syntax highlighting applies to code blocks using appropriate colors for the detected language.

The Tool Use Section appears when the assistant requests tool operations, rendering each tool use as a sub-panel with purple accent (#7D56F4). The tool use panel header shows tool name, input parameters, and execution status (pending, running, success, error) with appropriate status colors.

Tool Use and Tool Result Display
Tool operations receive specialized rendering emphasizing their importance in the conversation flow:

Tool Use Request Panel displays within assistant messages using purple-tinted styling to distinguish from regular text:

The panel header shows "🔧 Tool Use" followed by the tool name in bold purple ("Read", "Bash", "Edit"). Parameters render beneath in a compact format showing key-value pairs in gray text. For example, a Read tool might show path: /src/utils.py indented beneath the header.

When the tool request is large or complex, parameters collapse by default with an expansion control. Expanded view shows full JSON-formatted parameters with syntax highlighting for readability.

Tool Result Panel appears as a separate user message (since tool results come from the user role in the JSONL format) but with specialized styling:

The panel uses purple-tinted background (#2D1B4E) matching the tool use accent color, creating visual connection. The header shows "📄 Tool Result" with the matching tool use ID and execution duration ("completed in 125ms").

Status Indicators prominently display success or error states:

Success results show a green checkmark (✓) in the header with green accent text for the word "Success". Error results show a red X (✗) with red accent text for "Error" and error message details in red below.

Content Rendering adapts to result type:

File content renders in code blocks with syntax highlighting based on file extension. Command output displays in monospace with preserved formatting and ANSI color codes rendered when supported. JSON responses receive pretty-printing with syntax highlighting. Large results (over 1000 lines) truncate with "Show more" expansion controls.

The tool result panel includes action icons for copying content, opening in external editor, or saving to file, appearing on hover in the header area.

Token Usage Visualization
Token accounting displays throughout the message viewer providing cost transparency:

Per-Message Tokens appear in each assistant message header as described above, showing input/output breakdown. Cache token information displays with the lightning bolt emoji when present: "12.3k cached ⚡".

Cumulative Token Display appears as a sticky footer at the bottom of the detail pane showing running totals as users scroll through the conversation. The display format shows "Total: 127.5k tokens (45.2k in, 67.8k out, 14.5k cached) · Est. cost: $1.89" in medium gray text.

Visual Token Bar (optional, enabled in settings) renders a small horizontal bar in the conversation header showing relative proportions of input tokens (blue), output tokens (green), and cached tokens (yellow). The bar segments scale proportionally with hover tooltips showing exact counts.

When users hover over or focus token displays, tooltips appear explaining the token types and their cost implications: "Input tokens: Text sent to Claude. Output tokens: Claude's response text. Cached tokens: Retrieved from cache, not billed."

Code Syntax Highlighting
Code blocks throughout message content receive rich syntax highlighting using Chroma:

Language Detection analyzes code fences for language specifiers (python, javascript) or attempts automatic detection for unfenced code blocks based on syntax patterns.

Color Scheme adapts to terminal capabilities and theme settings, using TrueColor syntax highlighting in capable terminals and degrading to ANSI256 colors otherwise. The color scheme provides distinct colors for keywords (purple), strings (green), comments (gray), numbers (orange), and function names (blue).

Line Numbers appear in code blocks exceeding 10 lines, rendering in dark gray in the left margin. Line numbers aid in discussing code with Claude and help users reference specific lines when providing context.

Horizontal Scrolling handles long lines in code blocks rather than wrapping, preserving code structure. Scroll indicators (◀ ▶) appear when content exceeds visible width, and users scroll horizontally with arrow keys or mouse wheel when focused on the code block.

Copy Code Action appears as a small icon in the top-right corner of code blocks, allowing users to copy the code content to clipboard with a single click or keystroke. A brief confirmation message "Copied!" appears on successful copy.

Search and Filter Interface Specification
Global Search Activation and Layout
Search mode activates through multiple entry points providing flexibility for user workflows:

Pressing forward slash (/) from any view immediately activates search with scope defaulting to current project if a project is selected in the sidebar, or all projects if no specific project is focused. Pressing Ctrl+F activates the same search interface. The command palette includes a "Search Conversations" command allowing activation through that workflow as well.

The search interface renders as a centered modal overlay with the following characteristics:

Background Dimming applies a semi-transparent dark overlay (50% opacity black) over the entire application, creating focus on the search interface while maintaining context visibility beneath.

Search Modal Box appears centered both horizontally and vertically, sized at 80 characters wide and 30 rows tall (or smaller if terminal height is limited). The box uses rounded borders in purple (#7D56F4) with a slightly lighter background (#1F2937) than the main interface to create depth.

Modal Header displays "🔍 Search Conversations" in bold white text on purple background spanning the full modal width. A small "ESC to close" hint appears right-aligned in the header.

Search Input and Query Syntax
The search input field appears prominently at the top of the modal:

Input Field Styling uses a purple underline indicating focus, placeholder text in medium gray reading "Search messages, tools, or metadata...", and live character count in the top-right corner showing "0/500" to indicate the search query length limit.

Query Parsing supports sophisticated search syntax enabling precise filtering:

Plain Text matches conversation content, participant names, and summaries using case-insensitive substring matching. Multiple words search as phrase unless separated by operators.

Quoted Phrases use double quotes to search exact phrases: "optimize this function" matches only conversations containing that exact phrase. Single quotes work equivalently.

Boolean Operators combine search terms: python AND performance requires both terms, bug OR error matches either term, and code NOT typescript excludes conversations mentioning typescript.

Field-Specific Searches target particular metadata fields:

title:optimization searches only conversation titles
tool:bash finds conversations using the Bash tool
date:2025-01-15 filters to specific date
cost:>1.00 finds conversations over $1
tokens:>50000 filters high-token conversations
branch:feature/new-api searches by git branch
Date Range Syntax supports flexible date expressions:

date:today or date:yesterday
date:this-week or date:last-week
date:2025-01 for month
date:2025-01-15..2025-01-20 for range
Combination Queries mix field searches with text: tool:read AND python AND date:this-week finds this week's conversations using the Read tool that mention Python.

The search input provides autocomplete suggestions appearing in a dropdown below the field, showing recently used queries, common field names, and tool names as users type. Arrow keys navigate suggestions, Enter accepts the highlighted suggestion.

Search Scope Selection
A scope selector appears directly below the search input:

Scope Options render as radio buttons or pill toggles:

Current Project (default when project selected) - scope is limited to conversations within the focused sidebar project
All Projects - searches across all discovered projects
Selected Project(s) - appears when multiple projects are selected, searching only those projects
Current Conversation - appears when in detail view, searching within the displayed conversation only
The selected scope displays with purple accent, while unselected scopes show gray. A conversation count indicator shows "Searching 247 conversations" updating as scope changes.

Search Results Presentation
Below the scope selector, search results populate in real-time:

Results List renders similarly to the main conversation list but in a more compact format showing timestamp, title, and match excerpt. The excerpt displays up to 100 characters of context around the matched term with the term itself highlighted in yellow background.

Match Indicators show match count per conversation: "3 matches" in small gray text beneath each result item. Clicking a result opens that conversation in the detail pane and highlights the matched sections.

Result Ranking orders results by relevance score combining factors: exact phrase matches score higher than individual word matches, title matches score higher than content matches, more recent conversations score slightly higher, and match concentration (multiple matches near each other) increases score.

Live Updating refreshes results as users type with 300ms debounce, showing a small loading indicator during search execution. Result count displays in the header as "234 results" or "No results found".

Keyboard Navigation allows arrow keys to move through results, Enter to open the selected result, and Escape to dismiss search without opening anything.

Search History and Saved Searches
The search modal includes access to search history:

Recent Searches appear below the input field when empty, showing the last 10 unique queries in chronological order (most recent first). Each history item renders in gray with a small clock icon, becoming purple on hover or focus. Clicking a history item populates the search field with that query.

Clear History button appears at the bottom of history list, removing all saved searches when confirmed.

Saved Searches (optional feature) allows users to save frequently used complex queries with custom names. The save button (📌) appears next to the search input when a query is entered. Saved searches display in a separate tab or section accessible via keyboard shortcut, allowing instant re-execution of complex filter combinations.

Filter Construction Interface
Advanced filtering opens from a "Advanced Filters" button below the search input:

The advanced filter interface expands inline, showing filter builder controls:

Filter Groups organize by category (Time, Cost, Tokens, Tools, Branches, Projects) with each category expandable to show relevant controls.

Time Filters provide:

Date range picker with calendar interface for selecting start/end dates
Relative time presets (last hour, today, yesterday, this week, last week, this month, last month, this year)
Custom day-of-week filters (show only weekday/weekend conversations)
Numeric Range Filters for cost and tokens show:

Minimum and maximum input fields with validation
Quick preset buttons (e.g., "Under $1", "$1-$5", "Over $10" for cost)
Slider controls for visual range selection in capable terminals
Tool Filters display:

Checkbox list of all tools discovered in conversations
Search box to filter long tool lists
"Used" versus "Not Used" toggle for inclusion/exclusion logic
Tool usage count showing how many conversations match
Multi-Select Filters for projects and branches:

Searchable checkbox lists with select all/none controls
Hierarchical display for nested organization
Selected item count indicators
Filter Logic Toggles allow switching between AND/OR combination modes for filters within each group, displayed as radio buttons or toggle switches.

Apply Filters button at the bottom commits the filter configuration and updates the main conversation list. Reset button clears all advanced filters. Cancel dismisses the advanced interface without applying changes.

Command Palette Specification
Activation and Modal Presentation
The command palette provides keyboard-driven access to all application functionality:

Activation occurs via Ctrl+P (or Cmd+P on macOS) from any application state. The activation shortcut works globally regardless of current pane focus or active modal, making it a reliable escape hatch for navigation.

Modal Overlay renders similarly to the search interface with semi-transparent dark background (50% opacity black) dimming the underlying application. The modal box itself measures 70 characters wide and up to 25 rows tall, centered on screen.

Modal Header displays "⌘ Command Palette" in bold white text on purple background with "ESC to cancel" hint right-aligned.

Input Field and Fuzzy Matching
The command input field appears prominently below the header:

Input Field spans the full modal width with purple underline focus indicator and placeholder text "Type a command or search...". Character count shows in the corner.

Fuzzy Matching filters commands as users type, supporting:

Substring matching ignoring case ("conv" matches "Conversation")
Abbreviation matching ("gcr" matches "Go to Conversation Recent")
Out-of-order matching ("statproj" matches "Project Statistics")
Typo tolerance for common mistakes
Match Highlighting shows matched characters in bold purple within command names in the results list, helping users understand why certain commands appear.

Command List Organization
Commands organize into logical groups for discoverability:

Navigation Commands control view and focus movement:

"Go to Sidebar" - Focuses sidebar pane
"Go to Conversation List" - Focuses conversation list
"Go to Detail Pane" - Focuses detail/message viewer
"Switch Project: [name]" - Changes to specific project (shows recent projects)
"Back to List View" - Returns from detail to list (when applicable)
View Commands change application state:

"Toggle Statistics Dashboard" - Opens/closes stats view
"Toggle Help" - Shows/hides help modal
"Refresh Conversations" - Reloads conversation data from disk
"Compact View" - Reduces spacing for denser display (toggle)
Search and Filter Commands provide quick access:

"Search Conversations" - Opens search modal
"Filter by Date Range" - Opens date range picker
"Filter by Tool: [tool name]" - Applies specific tool filter
"Clear All Filters" - Removes active filters
"Show High Cost Conversations" - Applies cost filter
Data Actions perform operations on conversations:

"Export Current Conversation" - Saves conversation to file
"Export Selected Conversations" - Batch export
"Copy Conversation Link" - Copies deep link to clipboard
"Refresh Project Index" - Rescans .claude directory
Application Actions control global settings:

"Toggle Theme" - Switches between light/dark (if supported)
"Open Settings" - Shows settings modal
"Show Keyboard Shortcuts" - Displays help
"Quit Application" - Exits with confirmation
Each command displays with its keyboard shortcut (if available) right-aligned in gray text. Commands that require additional input show "..." after their name.

Command Execution Patterns
Command execution follows different patterns based on command type:

Immediate Commands execute instantly when selected with Enter key, dismissing the command palette and performing the action. Examples include navigation commands and view toggles.

Parameter Commands transition to a parameter input interface when selected. For example, "Filter by Date Range" replaces the command list with a date picker interface showing start date and end date inputs with calendar widgets. Pressing Enter applies the filter, Escape cancels back to command list.

Confirmation Commands show a confirmation dialog before executing, particularly for destructive actions. "Clear Cache" prompts "Are you sure? This will remove all cached conversation data." with Yes/No options.

Progressive Commands show progress feedback for long-running operations. "Export All Conversations" displays a progress bar in the command palette modal showing "Exporting... 45/234 conversations" with a cancel button.

Command Context and Availability
Command availability adapts to current application state:

Context-Sensitive Commands appear only when relevant. "Export Current Conversation" only shows when a conversation is selected in detail view. "Clear Filters" only appears when filters are active.

Disabled Commands render in gray with explanatory text. "Go to Detail Pane" disables in list view showing "(select a conversation first)" beneath the command name.

Recent Commands appear at the top of the list when command palette opens with empty input, showing the 5 most recently executed commands with a small clock icon prefix. This enables rapid repeated actions.

Favorite Commands (optional feature) allows users to pin frequently used commands to the top of the list with a star icon, persisting across sessions.

Command Customization
Advanced users can extend the command palette:

Custom Commands support user-defined shortcuts executing shell commands or scripts. Configuration occurs through a commands.yaml file allowing definitions like:

commands:
  - name: "Export to HTML"
    shortcut: "Ctrl+E"
    script: "./scripts/export-html.sh {{conversation_id}}"
    description: "Export current conversation as HTML"
Custom commands appear in a separate "Custom" group in the command palette with distinct icon prefix.

Statistics Dashboard Specification
Dashboard Layout and Access
The statistics dashboard provides comprehensive usage analytics:

Access Methods include pressing S key from any view, selecting "View Statistics" from command palette, or clicking a stats icon in the application header (if mouse UI enabled).

Layout Mode offers two options:

Full-Screen Mode (default) replaces the entire interface with the statistics dashboard, providing maximum space for charts and data. A header bar shows "📊 Statistics Dashboard" with "ESC to close" hint. The close action returns users to their previous view state.

Overlay Mode (accessible via command palette or settings) displays statistics as a large centered modal over the current view with semi-transparent background, allowing users to reference conversation list while viewing stats.

Dashboard Content Organization
The dashboard organizes into six primary sections arranged in a grid layout:

Overview Cards (top row, 4 cards horizontally) display key metrics in large-number format:

Card 1: Total Projects shows project count with icon 📁 and subtext "across [N] directories"

Card 2: Total Conversations shows conversation count with icon 💬 and subtext showing count from current month

Card 3: Total Messages shows message count with icon 📝 and subtext showing average per conversation

Card 4: Total Cost shows cumulative cost with icon 💰 and subtext showing average per conversation

Each card uses a lightly colored background with the metric number prominently displayed in large text. Cards respond to hover/focus by highlighting and showing a "click for details" hint.

Temporal Analysis Charts (middle section, 2 charts horizontally):

Left chart: Conversations Over Time displays a line chart showing conversation count per day over the selected time range (defaults to last 30 days). The chart uses NimbleMarkets/ntcharts line chart component with purple line color, grid lines for readability, and axis labels. Hover shows exact counts for each day.

Right chart: Token Usage Trends shows a stacked area chart with three layers: input tokens (blue), output tokens (green), and cache tokens (yellow). This visualizes token consumption patterns and cache efficiency over time.

Tool Usage Analysis (lower left section, 40% width):

Displays a horizontal bar chart ranking tools by usage count. Each bar shows the tool name, usage count, and success rate percentage. Colors code based on success rate: green for >90%, yellow for 70-90%, red for <70%. Clicking a tool bar filters the conversation list to show only conversations using that tool.

Cost Distribution (lower middle section, 30% width):

Shows a pie chart breaking down cost by project, with slices colored distinctly and labeled with project name and percentage. A legend lists all projects with exact cost values. An option toggles between "by project" and "by time period" views.

Activity Heat Map (lower right section, 30% width):

Displays a heat map showing conversation activity by hour of day and day of week. Darker colors indicate higher activity. This helps users understand their usage patterns. Hover shows exact conversation counts for each cell.

Date Range Controls appear at the top of the dashboard allowing users to adjust the analysis time window:

Preset buttons for "Last 7 Days", "Last 30 Days", "Last 3 Months", "Last Year", and "All Time". A custom range option opens date pickers. The selected range updates all dashboard visualizations simultaneously.

Interactive Features
Dashboard elements support rich interactions:

Chart Interactions include:

Hover tooltips showing exact values and additional context
Click on data points to filter conversation list to that segment
Drag to zoom into chart regions (for time series)
Double-click to reset zoom
Card Drill-Down opens detailed views when cards are clicked:

Clicking "Total Conversations" opens a breakdown showing conversations per project, conversations per day, and distribution by message count. Similar drill-downs exist for other metric cards.

Export Dashboard button (top-right corner) exports current statistics as:

PNG image (renders dashboard in capable terminals)
CSV data files (one per chart/section)
Markdown report (formatted text with ASCII charts)
JSON data (raw statistics for external analysis)
Real-Time Updates
The dashboard supports optional real-time updating:

Auto-Refresh Toggle in the dashboard header enables automatic statistics recalculation every 30 seconds, useful when monitoring active usage or when conversations are being added/modified.

Manual Refresh button forces immediate recalculation, showing a loading overlay during processing. Large datasets may take several seconds to analyze.

Cached Indicators show small icons on charts indicating whether data comes from cache (fast but potentially stale) or live calculation (slower but current). Users can force live recalculation via settings.

Help System and Keyboard Shortcuts
Help Modal Structure
The help system provides comprehensive documentation accessible at any time:

Activation occurs via ? key from any view, "Show Help" in command palette, or F1 key. The help modal appears as a centered overlay sized at 90 characters wide and 35 rows tall with rounded purple border.

Modal Header displays "❓ Claudex Help" with search box for filtering help content and "ESC to close" hint.

Content Sections organize using a tabbed interface at the top:

Keyboard Shortcuts tab (default) lists all keyboard commands
Getting Started tab provides quick-start guide
Features tab explains major functionality
Troubleshooting tab addresses common issues
About tab shows version, credits, and links
Tabs render as clickable buttons with purple highlight for active tab. Number keys 1-5 provide quick tab switching.

Keyboard Shortcuts Reference
The shortcuts tab organizes commands by context:

Global Shortcuts (available everywhere):

? or F1          Show this help
Ctrl+P           Open command palette
Ctrl+Q or q      Quit application
Ctrl+C           Force quit (in most contexts)
Navigation Shortcuts:

Tab              Move focus to next pane
Shift+Tab        Move focus to previous pane
1, 2, 3          Jump to sidebar, list, detail panes
Escape           Return to previous view/close modal
b                Back to conversation list (from detail)
Sidebar Shortcuts (when focused):

↑↓ or j/k        Navigate projects/filters
→ or l           Expand group
← or h           Collapse group
Enter/Space      Select/activate item
1-9              Quick jump to project 1-9
Conversation List Shortcuts (when focused):

↑↓ or j/k        Navigate conversations
Page Up/Down     Scroll by page
Home/End or gg/G Jump to start/end
Enter            Open conversation details
Space            Toggle selection (multi-select)
Shift+↑↓         Extend selection
Ctrl+A           Select all (filtered)
/ or Ctrl+F      Search conversations
s                Sort by column (when header focused)
Detail Pane Shortcuts (when focused):

↑↓ or j/k        Scroll message viewer
Page Up/Down     Scroll by page
Home/End or gg/G Jump to conversation start/end
Space            Expand/collapse section
c                Copy focused message
e                Expand all collapsed sections
Search Modal Shortcuts (when active):

Escape           Close search
Enter            Apply search/select result
↑↓               Navigate search results
Ctrl+R           Show recent searches
Ctrl+K           Clear search query
Command Palette Shortcuts (when active):

Escape           Close palette
Enter            Execute command
↑↓               Navigate commands
Ctrl+R           Show recent commands
Each shortcut displays with the keys in a monospace pill-shaped box (like Ctrl+P) followed by the action description. Alternative shortcuts for the same action show with "or" between them.

Contextual Help Integration
The help system provides context-aware assistance:

First-Use Hints appear as small tooltip overlays when users first interact with features. For example, when first focusing the sidebar, a hint appears: "Use arrow keys to navigate, Enter to select. Press ? for full help." Hints dismiss after 3 seconds or on next interaction.

Feature Flags track which features users have discovered, enabling progressive disclosure of advanced functionality. Settings allow disabling hints for experienced users.

Inline Help Icons appear next to complex features (like advanced filters) showing a small "ⓘ" icon. Hovering or focusing these icons displays expanded explanations.

Status Bar Hints show relevant shortcuts for the current context in the footer bar. When sidebar is focused, hints show sidebar navigation keys. When a search is active, hints show search-specific shortcuts.

Extended Documentation
The Features tab provides detailed explanations:

Feature Walkthroughs include:

"Browsing Conversation History" - explains three-pane navigation
"Searching and Filtering" - covers search syntax and filter options
"Understanding Token Usage" - explains input/output/cache tokens
"Analyzing Statistics" - guides through dashboard features
"Exporting Data" - describes export formats and uses
Each walkthrough combines text explanation with ASCII art screenshots showing the interface elements being described.

Common Workflows section provides step-by-step guides:

"Finding conversations about a specific topic"
"Tracking costs across projects"
"Identifying high-token conversations"
"Reviewing tool usage patterns"
Troubleshooting Guide
The Troubleshooting tab addresses common issues:

Common Issues section includes:

Problem: "No conversations appear in the list" Solutions: Check ~/.claude/projects/ exists, verify project paths are correct, check file permissions, try "Refresh Index" in command palette.

Problem: "Application is slow with large conversation history" Solutions: Enable conversation caching in settings, use filters to narrow results, consider archiving old conversations.

Problem: "Token counts seem incorrect" Solutions: Refresh conversation data, verify JSONL files are not corrupted, check that all assistant messages include usage metadata.

Problem: "Search returns unexpected results" Solutions: Review search syntax, check active filters, try quoted phrases for exact matches, clear search history if stale.

Each issue includes expandable "Technical Details" with log file locations, debug mode instructions, and reporting guidelines for bugs.

Visual Design System and Theming
Color Palette Definition
The application employs a comprehensive color system optimized for terminal readability:

Primary Colors establish brand identity and interactive elements:

Purple #7D56F4 - primary accent for borders, highlights, active states
Deep Purple #2D1B4E - tinted backgrounds for selected items
Light Purple #A78BFA - hover states and secondary accents
Semantic Colors convey meaning through color association:

Success Green #04B575 - successful operations, assistant role
Error Red #EF4444 - errors, destructive actions
Warning Yellow #F59E0B - cautions, high-cost indicators
Info Blue #60A5FA - informational elements, user role
Role-Based Colors distinguish conversation participants:

User Blue background #1E3A5F, border #60A5FA
Assistant Green background #0D3B2F, border #04B575
Tool Result Purple background #2D1B4E, border #7D56F4
System Gray background #1F2937, border #6B7280
Neutral Grays provide structure and hierarchy:

Background Dark #111827 - main application background
Background Medium #1F2937 - pane backgrounds
Border Gray #374151 - default borders and dividers
Text Light #D1D5DB - primary text content
Text Medium #9CA3AF - secondary text, metadata
Text Dark #6B7280 - tertiary text, hints
Adaptive Color Strategy
The color system adapts to terminal capabilities and themes:

Terminal Detection occurs at startup, identifying color support level (4-bit/16 colors, 8-bit/256 colors, 24-bit/TrueColor) and background lightness (light versus dark).

Adaptive Color Mapping using Lipgloss AdaptiveColor provides appropriate colors for light and dark terminals:

For light terminals (white background):

Primary Purple maps to #6D28D9 (darker purple)
Text Light maps to #1F2937 (dark gray)
Backgrounds use lighter tones
For dark terminals (black background):

Primary Purple uses #7D56F4 (standard purple)
Text Light uses #D1D5DB (light gray)
Backgrounds use darker tones
Color Degradation ensures graceful fallback in limited terminals:

In 256-color mode, TrueColor hex values map to nearest ANSI256 equivalents while maintaining relative contrast and semantic meaning.

In 16-color mode, the palette simplifies to basic ANSI colors: purple becomes bright magenta, greens remain green, blues remain blue, with careful contrast testing ensuring readability.

Typography System
Text styling creates visual hierarchy through weight, style, and size:

Font Styles apply semantic meaning:

Bold for headers, titles, focused items, command names
Italic for secondary information, hints, thinking content
Regular for primary content, conversation text
Monospace for code blocks, file paths, technical data
Text Size Simulation in terminals lacking true font sizing uses techniques like:

Large text: bold + uppercase for headers
Normal text: regular weight and case
Small text: regular weight + dimmed color for metadata
Text Alignment follows consistent rules:

Left-aligned for content, conversation text, lists
Right-aligned for metadata, timestamps, token counts
Center-aligned for headers, empty states, loading messages
Spacing and Layout System
Consistent spacing creates visual rhythm:

Padding Values use a 4-point base unit:

Tight: 1 space (0.25 units) - within small components
Normal: 2 spaces (0.5 units) - standard padding
Comfortable: 4 spaces (1 unit) - section padding
Spacious: 8 spaces (2 units) - major section separation
Margin Rules define component relationships:

Related items: 1 row spacing
Component sections: 2 rows spacing
Major sections: 3-4 rows spacing
Line Height ensures readability:

Single-spacing for compact lists and metadata
1.5-spacing for prose content in messages
Double-spacing for distinct conversation segments
Icon and Emoji Usage
Icons provide visual anchors and semantic meaning:

Role Indicators:

👤 User messages
🤖 Assistant messages
📄 Tool results
⚙️ System messages
Status Indicators:

✓ Success (green)
✗ Error (red)
⏳ Loading (yellow)
⚡ Cached (yellow)
💭 Thinking (gray)
Action Indicators:

📤 Export
🔗 Copy link
🔄 Refresh
🔍 Search
📊 Statistics
⋮ More options
Navigation Indicators:

▶ Collapsed/Expandable (right-pointing)
▼ Expanded (down-pointing)
→ Next/Forward
← Back
↑ Up/Top
↓ Down/Bottom
Icons always precede their associated text with one space separation. Emoji render natively in capable terminals and degrade to ASCII alternatives in limited environments.

Border and Divider Styling
Borders create structure and focus indication:

Border Styles use Lipgloss borders:

RoundedBorder (╭─╮╰─╯│) for panes and modals - modern, friendly
NormalBorder (┌─┐└─┘│) for tables and structured data - clean, precise
Hidden borders for seamless layouts when structure is implicit
Border Weight indicates hierarchy:

Thin borders (standard line weight) for internal divisions
Emphasized borders (color contrast) for focus indication
No borders for tightly integrated content
Divider Styles separate content:

Horizontal: repeated ─ in gray for section separation
Vertical: repeated │ in gray for column division
Heavy: ━ for major divisions (headers, footers)
Animation and Transition Specifications
Smooth animations enhance perceived responsiveness:

Scroll Animation uses Harmonica spring physics:

Angular frequency: 6.0 (moderately fast)
Damping ratio: 1.0 (critically damped, no bounce)
Frame rate: 60 FPS target
Duration: 200-300ms typical
Fade Transitions between views:

Opacity interpolation from 0 to 1 over 150ms
Easing: ease-out for appearing, ease-in for disappearing
No fade on fast navigation (keyboard shortcuts)
Progress Animation for loading states:

Animated spinner using bubbles/spinner Dot style
Rotation speed: 10 frames per full rotation
Color: purple accent matching theme
Loading Skeleton Animation:

Pulsing opacity between 30% and 60%
Pulse period: 1.5 seconds
No animation in slow terminals (< 30 FPS)
Performance Constraints:

All animations degrade gracefully in slow terminals
Animations disable automatically if frame rate drops below 30 FPS
User setting allows disabling animations entirely
Critical operations never block on animations
Error Handling and Edge Cases
File System Error States
The application gracefully handles various file system issues:

Permission Denied Errors display modal with:

Title: "⚠️ Permission Error"
Message: "Cannot access conversation history at ~/.claude/projects/"
Explanation: "Claudex needs read access to browse conversation files."
Actions: "Open Terminal Guide" (shows chmod commands), "Change Directory" (allows selecting different path), "Quit"
File Not Found Errors occur when:

~/.claude/ directory doesn't exist → Show "No Claude Code History Found" with setup guide
Specific conversation file missing → Show warning in conversation list with "File missing" indicator
Project directory empty → Show "No conversations in this project" empty state
Corrupted JSONL Errors handle gracefully:

Line-level errors skip the malformed line, log error, continue parsing
File-level errors mark conversation as "corrupted" in list with warning icon
Users can view technical error details through context menu
"Attempt Recovery" action tries alternative parsing strategies
Data Parsing Edge Cases
The application handles malformed or unexpected data:

Missing Required Fields like sessionId or timestamp:

Generate fallback values (UUID for sessionId, file modification time for timestamp)
Mark conversation with warning indicator
Display "(incomplete metadata)" in conversation list
Unexpected Field Types such as string instead of number:

Attempt type coercion with validation
Use default values for failed coercion
Log warning for debugging
Empty Conversation Files with zero messages:

Show in list as "Empty conversation" with gray styling
Detail pane displays "This conversation contains no messages"
Offer "Delete" and "Refresh" options
Incomplete Tool Use Sequences where tool result doesn't match tool use:

Display both independently with warning indicators
Show "Orphaned tool result" or "Unfulfilled tool use" labels
Allow manual matching through context menu
Resource Constraint Handling
Performance degrades gracefully under constraints:

Memory Limits when conversation history is very large:

Implement aggressive caching with LRU eviction
Show warning when memory usage exceeds 80% of limit
Offer "Compact Mode" reducing memory footprint
Suggest filtering to reduce dataset size
Disk Space Issues when caching requires storage:

Check available space before writing cache
Show warning when space below 100MB
Offer "Clear Cache" to free space
Gracefully disable caching if space unavailable
Slow Disk Access on network drives or slow media:

Show loading indicators during slow operations
Implement longer timeouts for file operations
Cache more aggressively to reduce disk hits
Show "Performance degraded - slow disk" status
Terminal Size Constraints in very small terminals:

Show persistent banner "Terminal too small - resize for better experience" when below 80x24
Switch to minimal single-pane layout
Disable complex features like statistics dashboard
Provide text-only fallbacks for visual elements
Error Message Templates
All errors follow consistent message structure:

Error Title (bold red): Brief description Error Message (normal white): Clear explanation of what went wrong
Technical Details (collapsed gray): File paths, error codes, stack traces Suggested Actions (purple buttons): "Retry", "Skip", "Open Settings", "Report Bug", "Quit" Help Link (blue underlined): "Learn more about this error"

Example error modal:

╭─────────────────────────────────────────────╮
│  ⚠️  Failed to Load Conversations           │
├─────────────────────────────────────────────┤
│                                              │
│  Claudex encountered errors while reading   │
│  conversation files from disk.              │
│                                              │
│  3 of 247 conversations could not be loaded │
│  due to file corruption or permission       │
│  issues.                                     │
│                                              │
│  ▶ Show technical details                   │
│                                              │
│  [ Continue with 244 conversations ]        │
│  [ Retry loading all files ]                │
│  [ Open troubleshooting guide ]             │
│                                              │
╰─────────────────────────────────────────────╯
Validation and Bounds Checking
Input validation prevents errors before they occur:

Search Query Validation:

Maximum length 500 characters
Sanitize special regex characters if regex mode enabled
Warn about query complexity (too many OR clauses)
Numeric Input Validation:

Token ranges must be positive integers
Cost ranges must be non-negative decimals with max 2 decimal places
Date ranges must have start before end
File Path Validation:

Check existence before operations
Verify write permissions for export
Validate path format and characters
Configuration Validation:

Ensure window dimensions are positive and reasonable
Validate color hex codes
Check command palette shortcuts don't conflict
Recovery Mechanisms
The application provides multiple recovery options:

Automatic Retry with exponential backoff:

First retry immediately
Second retry after 1 second
Third retry after 3 seconds
Fourth retry after 10 seconds
Give up after 4 retries, show error
Fallback to Cached Data:

When fresh data unavailable, use last successful cache
Display "Showing cached data" indicator with timestamp
Offer manual refresh action
Graceful Degradation:

Missing features disable cleanly with explanatory messages
Complex features simplify when resources limited
Read-only mode when write permissions unavailable
Manual Recovery Options:

"Rebuild Index" scans all files fresh
"Clear Cache" removes potentially corrupted cache
"Reset Settings" returns to defaults
"Recover Conversation" attempts deep parsing of corrupted files
Accessibility and Usability Considerations
Keyboard-Only Navigation
Complete functionality remains accessible without mouse:

Tab Order Logic follows natural reading flow:

Main panes: Sidebar → Conversation List → Detail Pane
Within lists: Top to bottom, wrapping at ends
Modal dialogs: Input fields → Action buttons → Close button
Settings: Groups → Fields within groups → Apply/Cancel
Skip-to-Content Shortcuts enable quick navigation:

Ctrl+1: Jump to sidebar
Ctrl+2: Jump to conversation list
Ctrl+3: Jump to detail pane
Ctrl+H: Jump to help
Ctrl+S: Jump to search
Dropdown Menu Access:

Arrow keys navigate menu items
Enter selects highlighted item
Escape closes without selection
First letter navigation jumps to items starting with typed letter
Focus Indicators clearly show keyboard position:

Focused elements display enhanced border/background
Focus ring never disappears during keyboard navigation
Focus moves logically and predictably
No focus traps where keyboard users get stuck
Screen Reader Considerations
While TUI screen reader support is limited, the application provides:

Descriptive Labels on all interactive elements:

Buttons: "Export conversation button"
Input fields: "Search conversations input field"
List items: "Project: my-project, 24 conversations, last active 2 hours ago"
Status Announcements for dynamic updates:

"Loading conversations... 45 of 234 loaded"
"Search complete. 12 results found"
"Conversation selected. Rendering messages"
Semantic Markup where possible:

Headers use bold and spacing for structure
Lists use consistent indentation
Tables use clear column alignment
Alternative Text for visual elements:

Emoji descriptions: 👤 becomes "user icon"
Charts: Text summary appears below visual
Progress bars: Percentage shown as text
Internationalization Support
The application supports global usage patterns:

Timestamp Formatting respects locale:

Date format: Uses ISO 8601 or locale preference (MM/DD/YYYY vs DD/MM/YYYY)
Time format: 12-hour vs 24-hour based on locale
Relative time: Localized strings ("2 hours ago" → "hace 2 horas")
Number Formatting follows locale conventions:

Thousands separator: Comma vs period (1,234 vs 1.234)
Decimal separator: Period vs comma (1.5 vs 1,5)
Currency: Symbol position and spacing ($1.50 vs 1,50€)
Text Encoding handles multi-byte characters:

UTF-8 throughout for full Unicode support
Proper character width calculation for CJK and emoji
Right-to-left text support (future enhancement)
Language Support (future):

UI strings externalized for translation
Command palette supports localized command names
Help text available in multiple languages
Performance Requirements
Specific performance targets ensure responsive experience:

Load Time Targets:

Application startup: < 500ms to first render
Project list: < 200ms to display
Conversation list: < 200ms to populate first page
Conversation detail: < 300ms to render
Search results: < 500ms for typical query
Responsiveness Targets:

Keyboard input latency: < 50ms perceived lag
Scroll smoothness: 60 FPS sustained
Focus change: < 100ms visual feedback
Filter application: < 300ms list update
Resource Usage Targets:

Memory: < 100MB for 1000 conversations
Disk cache: < 50MB typical
CPU: < 5% idle, < 30% active scrolling
Scalability Targets:

Support 10,000+ conversations without degradation
Handle individual conversations with 1000+ messages
Process search across 100,000+ messages in < 2 seconds
Progressive Disclosure Strategy
Features reveal complexity gradually:

First Launch Experience:

Welcome screen with quick-start guide
Highlight sidebar showing projects
Prompt to select first project
Show conversation list with tooltip hints
Encourage opening first conversation
Brief feature tour highlighting key shortcuts
Feature Introduction Timing:

Basic navigation: Immediate (first launch)
Search: After viewing 3+ conversations
Filters: After performing first search
Statistics: After 10+ conversation views
Command palette: After using 5+ different features
Advanced features: Show in "Pro Tips" hints
Hint System Progression:

Level 1 (New User): Show all hints, detailed explanations
Level 2 (Active User): Show hints for unused features only
Level 3 (Power User): Minimal hints, shortcuts emphasized
Level 4 (Expert): No hints, assume familiarity
Tutorial Mode (optional):

Interactive walkthrough covering all major features
Skippable at any time
Resumable from settings
Includes practice exercises with sample data
User Preference Management
Extensive customization maintains user preferences:

Settings Categories:

Appearance Settings:

Theme: Auto-detect, Light, Dark, Custom
Color accents: Purple (default), Blue, Green, Red, Custom
Animation speed: None, Slow, Normal (default), Fast
Compact mode: On/Off
Behavior Settings:

Auto-refresh: On/Off, interval selection
Default view: List, Last viewed, Statistics
Scroll speed: Slow, Normal (default), Fast
Confirmation dialogs: All, Destructive only, None
Data Settings:

Cache size limit: 10MB, 50MB (default), 100MB, Unlimited
Cache location: ~/.cache/claudex (default), Custom
Auto-export: On/Off, format selection, schedule
Data retention: Keep all, Archive old, Delete after period
Advanced Settings:

Debug mode: On/Off
Performance mode: Auto, Optimize for speed, Optimize for memory
Custom commands: Manage custom command palette entries
Experimental features: Enable/disable beta features
Settings Persistence:

Saved in ~/.config/claudex/settings.yaml
Auto-save on change with debounce
Manual save button for batch changes
"Reset to defaults" for each category or all settings
Settings Import/Export:

Export settings to YAML file for backup
Import settings from file for sync across machines
Share settings with team members
Validate imported settings before applying
Implementation Priorities and Phasing
Minimal Viable Product (MVP) Scope
The initial release focuses on core browsing and viewing:

Phase 1 Features:

Three-pane layout with basic responsiveness
Project discovery and sidebar list
Conversation list with essential columns (timestamp, title, cost)
Basic message viewer with user/assistant messages
Simple tool use/result display
Keyboard navigation (arrow keys, Enter, Escape)
Basic search (text-only, no advanced syntax)
Quit functionality
Technical Foundation:

JSONL parsing with error handling
Bubble Tea architecture with component composition
Lipgloss styling system
Bubbles list and viewport components
Configuration file support
Success Criteria for MVP:

Users can browse all projects
Users can view any conversation
Users can search for text
Application handles 1000+ conversations performantly
No crashes or data corruption
Clear error messages for common issues
Enhanced Features (Phase 2)
Second release adds polish and power features:

Phase 2 Features:

Advanced search syntax (fields, operators, dates)
Filter system with quick and advanced filters
Command palette with fuzzy matching
Statistics dashboard with basic charts
Token usage displays and cost calculations
Multi-select and batch operations
Export functionality (markdown, JSON, CSV)
Help system with keyboard shortcuts reference
Settings and preferences
Improved visual design and animations
Technical Enhancements:

Evertras/bubble-table for enhanced list
erikgeiser/promptkit for command palette
Harmonic animations throughout
Caching layer for performance
Full error recovery system
Success Criteria for Phase 2:

Users can find specific conversations quickly
Users can analyze their usage patterns
Users can export data for external use
Power users discover command palette
Application feels polished and professional
Advanced Capabilities (Phase 3)
Future releases expand functionality:

Phase 3 Features:

Advanced statistics with NimbleMarkets/ntcharts
Custom tags and organization
Conversation annotations and notes
Saved searches and filter presets
Custom command palette commands
Theme customization and color schemes
Conversation comparison and diffing
Real-time conversation monitoring
Integration with external tools
Team usage analytics (if applicable)
Technical Sophistication:

Full-text search indexing
Advanced charting and visualization
Plugin/extension system
API for external integrations
Performance optimization for huge datasets
Open Technical Questions
Several areas require collaboration with engineering:

Performance at Scale:

What is the actual performance profile with 10,000+ conversations?
Should we implement conversation pagination or fully load everything?
Is full-text indexing (e.g., Bleve) worth the complexity?
Terminal Compatibility:

What is the minimum terminal size we should support?
How do we handle exotic terminals with limited capabilities?
Should we provide an ASCII-only fallback mode?
Data Validation:

How strict should JSONL parsing be?
What recovery strategies work best for corrupted data?
Should we validate conversation integrity on load?
Feature Scope:

Should Phase 1 include basic statistics or defer to Phase 2?
Is export critical for MVP or can it wait?
Does command palette belong in Phase 1 for power users?
Testing Strategy:

How do we test TUI applications effectively?
What automated testing is possible with Bubble Tea?
How do we test across different terminal types?
Conclusion
This specification provides a comprehensive blueprint for implementing claudex as a sophisticated, professional-grade terminal user interface for browsing Claude Code conversation history. The design balances developer ergonomics with visual polish, keyboard efficiency with mouse accessibility, and immediate usability with power-user features.

The three-phase implementation approach ensures a solid foundation with the MVP while providing a clear roadmap for enhancement. The specification intentionally provides more detail than strictly necessary for engineering implementation, allowing designers and engineers to collaborate on trade-offs and optimizations during development.

Key design principles throughout include consistency in keyboard shortcuts and visual patterns, progressive disclosure of complexity, graceful degradation under constraints, clear error messages with recovery options, and performance optimization for large datasets.

The Bubble Tea ecosystem provides mature, battle-tested components that make this ambitious TUI achievable without building everything from scratch. By following proven patterns from applications like gh-dash and Glow, claudex can deliver a first-class terminal experience that Claude Code users will find indispensable for understanding and leveraging their conversation history.
