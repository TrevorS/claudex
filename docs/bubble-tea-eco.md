# Bubble Tea Ecosystem - Comprehensive Technical Research



## Table of Contents

1. [Core Bubble Tea Framework](#1-core-bubble-tea-framework)

2. [Lipgloss Styling](#2-lipgloss-styling)

3. [Bubbles Components](#3-bubbles-components)

4. [Third-Party Libraries](#4-third-party-libraries)

5. [Performance Considerations](#5-performance-considerations)

6. [Testing](#6-testing)

7. [Code Patterns and Best Practices](#7-code-patterns-and-best-practices)

8. [Notable Example Applications](#8-notable-example-applications)



---



## 1. Core Bubble Tea Framework



### Latest Versions and Stability



**Current Versions (2025):**

- **v1.x**: Published September 17, 2025, imported by 9,310+ projects - **Production-ready and stable**

  - Install: `go get github.com/charmbracelet/bubbletea`

- **v2.x (Beta)**: Published October 30, 2025, imported by 234 projects - **Beta, active development**

  - Install: `go get github.com/charmbracelet/bubbletea/v2@v2.0.0-beta.3`



**Stability**: Bubble Tea is production-ready and used by 8,000+ applications including tools from **AWS**, **NVIDIA**, and other major organizations. It powers all of Charm Bracelet's applications including Glow, Soft Serve, and hundreds of open-source projects.



### Core Concepts: The Elm Architecture



Bubble Tea is based on **The Elm Architecture**, a functional, predictable pattern with three core methods:



#### The Three Core Methods



**1. Init() - Initial Command**

```go

func (m model) Init() tea.Cmd {

    // Return initial commands to run on startup

    return tea.Batch(someCommand, anotherCommand)

}

```



**2. Update(msg tea.Msg) - Event Handler**

```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        switch msg.String() {

        case "ctrl+c", "q":

            return m, tea.Quit

        case "enter":

            return m, doSomething()

        }

    case customMsg:

        // Handle custom messages

        m.data = msg.data

        return m, nil

    }

    return m, nil

}

```



**3. View() - Render UI**

```go

func (m model) View() string {

    // Return a string representing your entire UI

    // No need to worry about redrawing logic - Bubble Tea handles it

    return lipgloss.JoinVertical(

        lipgloss.Left,

        m.header.View(),

        m.content.View(),

        m.footer.View(),

    )

}

```



#### Architecture Principles



- **Unidirectional Data Flow**: Messages flow down, updates flow up

- **Single Source of Truth**: Model contains all application state

- **Predictability**: Same state + same message = same result

- **Testability**: Pure functions make testing straightforward



### Event Handling and Message Passing



#### Messages (tea.Msg)



A message represents something that happened. Built-in messages include:

- `tea.KeyMsg` - Keyboard input

- `tea.MouseMsg` - Mouse events (clicks, movement, scrolling)

- `tea.WindowSizeMsg` - Terminal resize events

- `tea.PasteMsg` (v2) - Clipboard paste events



Custom messages are just Go types:

```go

type recordUpdatedMsg struct {

    id   int

    data string

}



type errorMsg struct {

    err error

}

```



#### Commands (tea.Cmd)



**Definition**: `type Cmd func() tea.Msg`



Commands perform I/O and return messages. All I/O (network, disk, timers) should use commands:



```go

// Example: HTTP request command

func fetchData(url string) tea.Cmd {

    return func() tea.Msg {

        resp, err := http.Get(url)

        if err != nil {

            return errorMsg{err}

        }

        // Process response...

        return dataReceivedMsg{data}

    }

}

```



**Key Pattern**: Commands run asynchronously in goroutines. The message they return is sent to Update().



#### Batch - Concurrent Execution



Run multiple commands concurrently with no ordering guarantees:



```go

func (m model) Init() tea.Cmd {

    return tea.Batch(

        fetchUserData(),

        startTimer(),

        checkForUpdates(),

    )

}

```



#### Sequence - Sequential Execution



Run commands in order, one at a time:



```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg.(type) {

    case tea.KeyMsg:

        return m, tea.Sequence(

            saveData(),

            showNotification(),

            refreshUI(),

        )

    }

    return m, nil

}

```



#### Combining Batch and Sequence



```go

return tea.Batch(

    tea.Sequence(doThis, thenThat, thenThis),

    tea.Sequence(otherCmds...),

    anotherCmd,

)

```



### Best Practices for Component Composition



#### Component Structure



Any non-trivial application should use nested models:



```go

type AppModel struct {

    header   HeaderModel

    sidebar  SidebarModel

    content  ContentModel

    footer   FooterModel



    // Shared state

    windowSize tea.WindowSizeMsg

}

```



#### The Parent as Message Router



The top-level model becomes a **message router** and **screen compositor**:



```go

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    var cmds []tea.Cmd



    // Handle global messages

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:

        m.windowSize = msg

        // Propagate to all children

    case tea.KeyMsg:

        if msg.String() == "ctrl+c" {

            return m, tea.Quit

        }

    }



    // Route to children

    var cmd tea.Cmd

    m.header, cmd = m.header.Update(msg)

    cmds = append(cmds, cmd)



    m.content, cmd = m.content.Update(msg)

    cmds = append(cmds, cmd)



    return m, tea.Batch(cmds...)

}



func (m AppModel) View() string {

    return lipgloss.JoinVertical(

        lipgloss.Left,

        m.header.View(),

        lipgloss.JoinHorizontal(lipgloss.Top,

            m.sidebar.View(),

            m.content.View(),

        ),

        m.footer.View(),

    )

}

```



#### Component Interface Pattern



While Bubbles components don't satisfy `tea.Model` interface exactly (they return their own type, not `tea.Model`), they follow the same pattern:



```go

type MyComponent struct {

    state string

}



func (c MyComponent) Init() tea.Cmd {

    return nil

}



func (c MyComponent) Update(msg tea.Msg) (MyComponent, tea.Cmd) {

    // Note: returns MyComponent, not tea.Model

    return c, nil

}



func (c MyComponent) View() string {

    return c.state

}

```



### State Management Patterns



#### Single Source of Truth



**Problem**: Duplicated state between components leads to sync issues.



**Solution**: Keep state in one place, pass it down:



```go

type AppModel struct {

    items []Item  // Single source of truth

    list  list.Model

    table table.Model

}



func (m AppModel) View() string {

    // Both components render from same data

    return lipgloss.JoinHorizontal(

        lipgloss.Top,

        m.list.View(),

        m.table.View(),

    )

}

```



#### Window Size Management



Propagate window size to all components:



```go

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:

        m.windowSize = msg

        // Update child dimensions

        m.sidebar.SetSize(msg.Height, 30)

        m.content.SetSize(msg.Height, msg.Width-30)

    }

    return m, nil

}

```



#### State Machine Pattern



For complex UI flows, use explicit state machines:



```go

type ViewState int



const (

    ViewLoading ViewState = iota

    ViewList

    ViewDetail

    ViewEdit

)



type Model struct {

    state ViewState

    // ... other fields

}



func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch m.state {

    case ViewLoading:

        // Handle loading state

    case ViewList:

        // Handle list state

    case ViewDetail:

        // Handle detail state

    }

    return m, nil

}



func (m Model) View() string {

    switch m.state {

    case ViewLoading:

        return m.spinner.View()

    case ViewList:

        return m.list.View()

    case ViewDetail:

        return m.detail.View()

    }

    return ""

}

```



---



## 2. Lipgloss Styling



### Current Versions



- **v1**: Stable, widely used

- **v2**: Alpha/Beta (v2.0.0-alpha.2+) - Published March 26, 2025



**Install**:

- v1: `go get github.com/charmbracelet/lipgloss`

- v2: `go get github.com/charmbracelet/lipgloss/v2@v2.0.0-alpha.2`



### Color System



#### Automatic Color Profile Detection



Lipgloss automatically detects the terminal's color capabilities and adapts:



**Profiles**:

- **TrueColor**: 24-bit RGB (16.7M colors)

- **ANSI256**: 256 colors

- **ANSI**: 16 colors

- **Ascii**: No colors

- **NoTTY**: Non-terminal output



**Automatic Degradation**: Colors outside the terminal's gamut are automatically coerced to the closest available value.



#### Color Definition Methods



**1. Simple Colors**:

```go

style := lipgloss.NewStyle().

    Foreground(lipgloss.Color("205")).      // ANSI256 color

    Background(lipgloss.Color("#FF5733"))   // Hex color

```



**2. AdaptiveColor - Light/Dark Backgrounds**:

```go

style := lipgloss.NewStyle().

    Foreground(lipgloss.AdaptiveColor{

        Light: "#000000",  // Dark text on light background

        Dark:  "#FFFFFF",  // Light text on dark background

    })

```



The terminal's background color is automatically detected and the appropriate color chosen.



**3. CompleteColor - Explicit Per-Profile**:

```go

style := lipgloss.NewStyle().

    Foreground(lipgloss.CompleteColor{

        TrueColor: "#FF5733",

        ANSI256:   "205",

        ANSI:      "5",

    })

```



No automatic degradation - you specify exact values for each profile.



**4. CompleteAdaptiveColor - Full Control**:

```go

style := lipgloss.NewStyle().

    Foreground(lipgloss.CompleteAdaptiveColor{

        Light: lipgloss.CompleteColor{

            TrueColor: "#000000",

            ANSI256:   "0",

            ANSI:      "0",

        },

        Dark: lipgloss.CompleteColor{

            TrueColor: "#FFFFFF",

            ANSI256:   "255",

            ANSI:      "15",

        },

    })

```



### Border Styles



Lipgloss provides multiple built-in border styles:



```go

// Normal border - standard weight, 90° corners

lipgloss.NormalBorder()



// Rounded corners

lipgloss.RoundedBorder()



// Thicker than normal

lipgloss.ThickBorder()



// Double-line strokes

lipgloss.DoubleBorder()



// Hidden border - maintains spacing without visible border

lipgloss.HiddenBorder()

```



**Usage**:

```go

style := lipgloss.NewStyle().

    Border(lipgloss.RoundedBorder()).

    BorderForeground(lipgloss.Color("62")).

    Padding(1, 2)

```



**Partial Borders**:

```go

style := lipgloss.NewStyle().

    Border(lipgloss.RoundedBorder()).

    BorderTop(true).

    BorderBottom(true).

    BorderLeft(false).

    BorderRight(false)

```



### Layout Capabilities



#### Join Functions



**JoinVertical** - Stack content vertically:

```go

content := lipgloss.JoinVertical(

    lipgloss.Left,    // 0.0 = left, 0.5 = center, 1.0 = right

    header,

    body,

    footer,

)

```



**JoinHorizontal** - Arrange content horizontally:

```go

content := lipgloss.JoinHorizontal(

    lipgloss.Top,     // 0.0 = top, 0.5 = middle, 1.0 = bottom

    sidebar,

    mainContent,

    rightPanel,

)

```



#### Place Function



Position content within a defined space:



```go

positioned := lipgloss.Place(

    80,                   // width

    24,                   // height

    lipgloss.Right,       // horizontal position

    lipgloss.Bottom,      // vertical position

    styledContent,

)

```



**Position Values**:

- `0.0` = start (left/top)

- `0.5` = center

- `1.0` = end (right/bottom)



#### Table Layout (Lipgloss v2+)



```go

import "github.com/charmbracelet/lipgloss/v2/table"



t := table.New().

    Border(lipgloss.ASCIIBorder()).

    Headers("NAME", "AGE", "CITY").

    Rows(

        []string{"Alice", "25", "NYC"},

        []string{"Bob", "30", "SF"},

    )



fmt.Println(t)

```



### Text Styling



```go

style := lipgloss.NewStyle().

    Bold(true).

    Italic(true).

    Underline(true).

    Strikethrough(true).

    Blink(true).

    Faint(true).

    Foreground(lipgloss.Color("205")).

    Background(lipgloss.Color("235")).

    Width(50).

    Height(10).

    Align(lipgloss.Center).          // Horizontal alignment

    AlignVertical(lipgloss.Middle).  // Vertical alignment

    Padding(1, 2, 1, 2).            // Top, Right, Bottom, Left

    Margin(1, 2).                    // Vertical, Horizontal

    MaxWidth(80).

    MaxHeight(24)



rendered := style.Render("Hello, World!")

```



### Lipgloss v2 Improvements



**Major Changes**:

1. **Deterministic Styles**: Styles are now deterministic and predictable

2. **Precise I/O Control**: You control input/output sources explicitly

3. **Better Integration**: Works in lockstep with Bubble Tea v2 (no more lock-ups)

4. **Improved Color API**: More explicit color handling



**Breaking Changes**:

- Can no longer use hexadecimal and integer formats interchangeably for colors

- Default I/O behavior changed - no longer assumes stdin/stdout by default



---



## 3. Bubbles Components



**Repository**: https://github.com/charmbracelet/bubbles

**Docs**: https://pkg.go.dev/github.com/charmbracelet/bubbles



Bubbles provides production-ready TUI components used in Glow and many other applications.



### Available Built-in Components



#### 1. Viewport - Scrollable Content



**Purpose**: Display and scroll large content (logs, markdown, help text)



**Features**:

- Vertical scrolling

- Mouse wheel support

- Standard pager keybindings

- High-performance mode for alternate screen buffer

- Percentage and line number display



**Usage**:

```go

import "github.com/charmbracelet/bubbles/viewport"



type model struct {

    viewport viewport.Model

    content  string

}



func (m model) Init() tea.Cmd {

    return nil

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    var cmd tea.Cmd



    switch msg := msg.(type) {

    case tea.WindowSizeMsg:

        m.viewport = viewport.New(msg.Width, msg.Height-2)

        m.viewport.SetContent(m.content)

        return m, nil

    }



    m.viewport, cmd = m.viewport.Update(msg)

    return m, cmd

}



func (m model) View() string {

    return m.viewport.View()

}

```



**Performance**: High-performance rendering available but deprecated. Use for alternate screen buffer apps.



#### 2. List - Interactive Selection



**Purpose**: Browse, filter, and select from lists of items



**Features**:

- Filtering

- Pagination

- Help display

- Status messages

- Activity spinner

- Custom item rendering

- Delegate pattern for item behavior



**Usage**:

```go

import "github.com/charmbracelet/bubbles/list"



type item struct {

    title, desc string

}



func (i item) Title() string       { return i.title }

func (i item) Description() string { return i.desc }

func (i item) FilterValue() string { return i.title }



type model struct {

    list list.Model

}



func (m model) Init() tea.Cmd {

    return nil

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        if msg.String() == "enter" {

            selected := m.list.SelectedItem().(item)

            // Handle selection

        }

    case tea.WindowSizeMsg:

        m.list.SetSize(msg.Width, msg.Height)

    }



    var cmd tea.Cmd

    m.list, cmd = m.list.Update(msg)

    return m, cmd

}



func (m model) View() string {

    return m.list.View()

}

```



#### 3. Table - Tabular Data



**Purpose**: Display and navigate rows and columns



**Features**:

- Vertical scrolling

- Row selection

- Custom column widths

- Style customization

- Header, rows, footer



**Usage**:

```go

import "github.com/charmbracelet/bubbles/table"



columns := []table.Column{

    {Title: "Name", Width: 20},

    {Title: "Age", Width: 10},

    {Title: "City", Width: 20},

}



rows := []table.Row{

    {"Alice", "25", "NYC"},

    {"Bob", "30", "SF"},

}



t := table.New(

    table.WithColumns(columns),

    table.WithRows(rows),

    table.WithFocused(true),

    table.WithHeight(10),

)



s := table.DefaultStyles()

s.Header = s.Header.

    BorderStyle(lipgloss.NormalBorder()).

    BorderForeground(lipgloss.Color("240"))

t.SetStyles(s)

```



#### 4. Spinner - Loading Indicator



**Purpose**: Show activity/loading state



**Features**:

- Multiple built-in spinner styles

- Custom spinner frames

- Customizable colors and speed



**Usage**:

```go

import "github.com/charmbracelet/bubbles/spinner"



type model struct {

    spinner spinner.Model

    loading bool

}



func newModel() model {

    s := spinner.New()

    s.Spinner = spinner.Dot

    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

    return model{spinner: s, loading: true}

}



func (m model) Init() tea.Cmd {

    return m.spinner.Tick

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    if !m.loading {

        return m, nil

    }



    var cmd tea.Cmd

    m.spinner, cmd = m.spinner.Update(msg)

    return m, cmd

}



func (m model) View() string {

    if m.loading {

        return m.spinner.View() + " Loading..."

    }

    return "Done!"

}

```



**Built-in Spinners**: Line, Dot, MiniDot, Jump, Pulse, Points, Globe, Moon, Monkey



#### 5. Progress - Progress Bar



**Purpose**: Show progress of operations



**Features**:

- Solid and gradient color modes

- Percentage display

- Harmonica spring animations

- Custom colors and characters



**Usage**:

```go

import "github.com/charmbracelet/bubbles/progress"



type model struct {

    progress progress.Model

    percent  float64

}



func (m model) Init() tea.Cmd {

    return nil

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case progressMsg:

        m.percent = msg.percent

        if m.percent >= 1.0 {

            return m, tea.Quit

        }

        return m, tickCmd()

    }

    return m, nil

}



func (m model) View() string {

    return m.progress.ViewAs(m.percent)

}

```



#### 6. TextInput - Single-line Input



**Purpose**: Text input field



**Features**:

- Unicode support

- Paste support

- Horizontal scrolling for long values

- Placeholder text

- Character limits

- Input validation

- Password masking

- Custom styling



**Usage**:

```go

import "github.com/charmbracelet/bubbles/textinput"



type model struct {

    textInput textinput.Model

}



func initialModel() model {

    ti := textinput.New()

    ti.Placeholder = "Enter your name"

    ti.Focus()

    ti.CharLimit = 50

    ti.Width = 30



    return model{textInput: ti}

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    var cmd tea.Cmd



    switch msg := msg.(type) {

    case tea.KeyMsg:

        switch msg.String() {

        case "enter":

            value := m.textInput.Value()

            // Process input

            return m, tea.Quit

        }

    }



    m.textInput, cmd = m.textInput.Update(msg)

    return m, cmd

}

```



#### 7. Textarea - Multi-line Input



**Purpose**: Multi-line text editing



**Features**:

- Unicode support

- Paste support

- Vertical scrolling

- Line numbers

- Custom dimensions

- Character/line limits



**Usage**:

```go

import "github.com/charmbracelet/bubbles/textarea"



type model struct {

    textarea textarea.Model

}



func initialModel() model {

    ta := textarea.New()

    ta.Placeholder = "Enter your notes..."

    ta.Focus()

    ta.SetWidth(60)

    ta.SetHeight(10)



    return model{textarea: ta}

}

```



#### 8. Key - Keybinding Helper



**Purpose**: Define and document keyboard shortcuts



**Usage**:

```go

import "github.com/charmbracelet/bubbles/key"



type keyMap struct {

    Up    key.Binding

    Down  key.Binding

    Help  key.Binding

    Quit  key.Binding

}



var keys = keyMap{

    Up: key.NewBinding(

        key.WithKeys("up", "k"),

        key.WithHelp("↑/k", "move up"),

    ),

    Down: key.NewBinding(

        key.WithKeys("down", "j"),

        key.WithHelp("↓/j", "move down"),

    ),

    Help: key.NewBinding(

        key.WithKeys("?"),

        key.WithHelp("?", "toggle help"),

    ),

    Quit: key.NewBinding(

        key.WithKeys("q", "ctrl+c"),

        key.WithHelp("q", "quit"),

    ),

}



// In Update:

switch msg := msg.(type) {

case tea.KeyMsg:

    switch {

    case key.Matches(msg, keys.Quit):

        return m, tea.Quit

    case key.Matches(msg, keys.Up):

        // Handle up

    }

}

```



#### 9. Help - Help Display



**Purpose**: Show keyboard shortcuts to users



**Usage**:

```go

import "github.com/charmbracelet/bubbles/help"



type model struct {

    help help.Model

    keys keyMap

}



func (k keyMap) ShortHelp() []key.Binding {

    return []key.Binding{k.Help, k.Quit}

}



func (k keyMap) FullHelp() [][]key.Binding {

    return [][]key.Binding{

        {k.Up, k.Down},

        {k.Help, k.Quit},

    }

}



func (m model) View() string {

    return m.content + "\n" + m.help.View(m.keys)

}

```



### Custom Component Creation



To create a custom component, implement the three core methods:



```go

type CustomWidget struct {

    width  int

    height int

    style  lipgloss.Style

    data   []string

}



func NewCustomWidget() CustomWidget {

    return CustomWidget{

        style: lipgloss.NewStyle().

            Border(lipgloss.RoundedBorder()).

            Padding(1),

        data: []string{},

    }

}



func (w CustomWidget) Init() tea.Cmd {

    return nil  // Or return initial commands

}



func (w CustomWidget) Update(msg tea.Msg) (CustomWidget, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:

        w.width = msg.Width

        w.height = msg.Height

    case customDataMsg:

        w.data = msg.data

    }

    return w, nil

}



func (w CustomWidget) View() string {

    content := strings.Join(w.data, "\n")

    return w.style.

        Width(w.width).

        Height(w.height).

        Render(content)

}



// Helper methods

func (w *CustomWidget) SetData(data []string) {

    w.data = data

}

```



---



## 4. Third-Party Libraries



### 1. Evertras/bubble-table



**GitHub**: https://github.com/Evertras/bubble-table

**Docs**: https://pkg.go.dev/github.com/evertras/bubble-table/table



**Purpose**: Sophisticated, customizable table component



**Features**:

- **Sorting**: Multi-column sorting, numeric and string comparison

- **Filtering**: Built-in filtering capabilities

- **Pagination**: Automatic page management

- **Column Types**: Fixed-width and flexible columns

- **Horizontal Scrolling**: With frozen left columns

- **Styling**: Global, column, row, and cell-level styles

- **Borders**: Customizable shapes and colors

- **Footer**: Auto page info, custom text, or hidden



**Usage**:

```go

import "github.com/evertras/bubble-table/table"



columns := []table.Column{

    table.NewColumn("id", "ID", 5),

    table.NewFlexColumn("name", "Name", 1),

    table.NewColumn("age", "Age", 8),

}



rows := []table.Row{

    table.NewRow(table.RowData{

        "id":   "1",

        "name": "Alice",

        "age":  "25",

    }),

}



tbl := table.New(columns).

    WithRows(rows).

    Focused(true).

    Border(table.BorderRounded).

    WithPageSize(10).

    WithMaxTotalWidth(80).

    SortByAsc("name")



// In Update:

tbl, cmd = tbl.Update(msg)



// In View:

return tbl.View()

```



**Sorting**:

- Single or multi-column

- Ascending/descending

- Numeric value sorting (ints, floats)

- String comparison fallback



**When to Use**: Need advanced table features beyond basic Bubbles table (sorting, filtering, frozen columns, complex styling)



### 2. erikgeiser/promptkit



**GitHub**: https://github.com/erikgeiser/promptkit

**Docs**: https://pkg.go.dev/github.com/erikgeiser/promptkit



**Purpose**: Collection of command-line prompts (like a command palette)



**Status**: ⚠️ API not yet stable - expect breaking changes in minor versions



**Features**:

- Selection prompts with filtering and pagination

- Text input with validation

- Confirmation prompts (yes/no)

- Integrates as Bubble Tea widget

- Sensible defaults

- Re-mappable keybindings

- Heavy customization options



**Components**:

- `selection` - Filterable, paginated selection lists

- `textinput` - Text input with validation and defaults

- `confirmation` - Yes/no prompts



**Usage**:

```go

import (

    "github.com/erikgeiser/promptkit/selection"

    "github.com/erikgeiser/promptkit/textinput"

    "github.com/erikgeiser/promptkit/confirmation"

)



// Selection

sp := selection.New("Choose an option:", options)

sp.Filter = func(filter string, choice string) bool {

    return strings.Contains(strings.ToLower(choice), strings.ToLower(filter))

}

choice, err := sp.RunPrompt()



// Text input

input := textinput.New("Enter your name:")

input.Placeholder = "John Doe"

input.Validate = func(s string) error {

    if len(s) < 2 {

        return errors.New("name too short")

    }

    return nil

}

name, err := input.RunPrompt()



// Confirmation

confirm := confirmation.New("Are you sure?", confirmation.No)

ok, err := confirm.RunPrompt()

```



**Known Issues**: Windows not explicitly supported due to Bubble Tea input event bug



**When to Use**: Building interactive CLI prompts, wizards, or command palette UIs



### 3. charmbracelet/harmonica



**GitHub**: https://github.com/charmbracelet/harmonica

**Example**: https://github.com/charmbracelet/harmonica/blob/master/examples/spring/tui/main.go



**Purpose**: Spring-based animation library for smooth, natural motion



**Features**:

- Physics-based spring animations

- Frequency (speed) control

- Damping (bounciness) control

- Integrates with Bubble Tea

- Used by Bubbles progress component



**Usage**:

```go

import "github.com/charmbracelet/harmonica"



type model struct {

    spring harmonica.Spring

    x      float64

    target float64

}



func (m model) Init() tea.Cmd {

    return tick()

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tickMsg:

        m.x, m.spring = harmonica.Animate(m.x, m.target, m.spring)

        if !m.spring.Done() {

            return m, tick()

        }

    case targetChangedMsg:

        m.target = msg.newTarget

        return m, tick()

    }

    return m, nil

}

```



**Integration with Progress Bar**:

```go

p := progress.New(

    progress.WithSpringOptions(120, 1.0),  // frequency, damping

)

```



**When to Use**: Smooth animations for progress bars, transitions, sprite movement



### 4. NimbleMarkets/ntcharts



**GitHub**: https://github.com/NimbleMarkets/ntcharts

**Docs**: https://pkg.go.dev/github.com/NimbleMarkets/ntcharts



**Purpose**: Terminal charts for Bubble Tea framework



**Chart Types**:

- **Canvas**: 2D grid for plotting arbitrary runes (foundation for other charts)

- **Bar Charts**: Horizontal rows or vertical columns

- **Heatmaps**: Color-mapped (x,y) values

- **Line Charts**: (X,Y) data points on 2D grid

- **Candlestick Charts**: OHLC (Open, High, Low, Close) data

- **Time Series Line Charts**: Time on X-axis

- **Wave Line Charts**: Wave-pattern connections

- **Sparklines**: Small, simple data visualization



**Integration**: Uses Lipgloss for styling and BubbleZone for mouse support



**Usage Example**:

```go

import (

    "github.com/NimbleMarkets/ntcharts/linechart"

    "github.com/NimbleMarkets/ntcharts/canvas"

)



type model struct {

    chart linechart.Model

}



func (m model) Init() tea.Cmd {

    return nil

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:

        m.chart.SetSize(msg.Width, msg.Height)

    case dataUpdateMsg:

        m.chart.UpdateData(msg.points)

    }

    return m, nil

}



func (m model) View() string {

    return m.chart.View()

}

```



**When to Use**: Data visualization, dashboards, monitoring tools, analytics displays



### 5. alecthomas/chroma



**GitHub**: https://github.com/alecthomas/chroma

**Docs**: https://pkg.go.dev/github.com/alecthomas/chroma/v2



**Purpose**: General-purpose syntax highlighter in pure Go



**Features**:

- Based on Pygments

- Supports 100+ languages

- Multiple output formats (HTML, ANSI, terminal256, terminal16)

- 50+ built-in styles

- Language auto-detection

- Pygments lexer/style translation

- Used by Hugo static site generator



**Output Formats for TUI**:

- `terminal` - 8 color ANSI

- `terminal16` - 16 color ANSI

- `terminal256` - 256 color ANSI

- `terminal16m` - TrueColor (16 million colors)



**Usage for TUI**:

```go

import (

    "github.com/alecthomas/chroma/v2"

    "github.com/alecthomas/chroma/v2/formatters"

    "github.com/alecthomas/chroma/v2/lexers"

    "github.com/alecthomas/chroma/v2/styles"

)



// Quick API

highlighted, err := quick.Highlight(

    writer,

    code,

    "go",           // language

    "terminal256",  // formatter

    "monokai",      // style

)



// Detailed API

lexer := lexers.Get("go")

if lexer == nil {

    lexer = lexers.Fallback

}



style := styles.Get("monokai")

if style == nil {

    style = styles.Fallback

}



formatter := formatters.Get("terminal256")

if formatter == nil {

    formatter = formatters.Fallback

}



iterator, err := lexer.Tokenise(nil, code)

err = formatter.Format(writer, style, iterator)

```



**Language Detection**:

```go

// By filename

lexer := lexers.Match("main.go")



// By content analysis

lexer := lexers.Analyse(sourceCode)



// Explicit

lexer := lexers.Get("go")

```



**Popular Styles**: monokai, solarized-dark, dracula, github, vim



**When to Use**: Code viewers, file browsers, diff displays, documentation, REPLs



### 6. lrstanley/bubblezone



**GitHub**: https://github.com/lrstanley/bubblezone



**Purpose**: Easy mouse event tracking for Bubble Tea components



**Features**:

- Clickable regions/buttons

- Automatic bounds checking

- Eliminates manual coordinate calculation

- Simple zone marking API



**Usage**:

```go

import "github.com/lrstanley/bubblezone"



// Initialize global zone manager

var zone = bubblezone.NewGlobal()



type model struct {

    buttons []button

}



type button struct {

    id    string

    label string

}



func (m model) View() string {

    var buttons []string

    for _, btn := range m.buttons {

        // Mark clickable zone

        buttons = append(buttons, zone.Mark(btn.id, btn.label))

    }



    content := lipgloss.JoinVertical(lipgloss.Left, buttons...)



    // Scan to register zones

    return zone.Scan(content)

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.MouseMsg:

        if msg.Action == tea.MouseActionRelease &&

           msg.Button == tea.MouseButtonLeft {

            // Check which zone was clicked

            for _, btn := range m.buttons {

                if zone.Get(btn.id).InBounds(msg) {

                    // Handle button click

                }

            }

        }

    }

    return m, nil

}

```



**When to Use**: Building mouse-interactive TUIs with clickable elements



---



## 5. Performance Considerations



### Virtual Scrolling and Large Datasets



#### Viewport Component



The built-in `viewport` component is optimized for large content:



```go

viewport := viewport.New(width, height)

viewport.SetContent(largeContent)  // Can handle 1000s of lines



// High-performance mode (deprecated but available)

viewport.HighPerformanceRendering = true

```



**Key Features**:

- Only renders visible lines

- Efficient scrolling

- Mouse wheel delta of 3 lines

- Works with alternate screen buffer



**Gotcha**: High-performance rendering is deprecated - will be removed in future versions.



#### List Component



For large lists, use built-in pagination:



```go

list := list.New(items, delegate, width, height)

list.Paginator.PerPage = 20

list.Paginator.Type = paginator.Dots

```



The list component only renders visible items, making it efficient for 1000s of entries.



#### Custom Virtual Scrolling Pattern



For custom components with large datasets:



```go

type VirtualList struct {

    allItems     []Item      // Full dataset

    visibleStart int         // First visible index

    visibleCount int         // Number of visible items

    height       int         // Component height

}



func (v VirtualList) View() string {

    end := v.visibleStart + v.visibleCount

    if end > len(v.allItems) {

        end = len(v.allItems)

    }



    // Only render visible slice

    visible := v.allItems[v.visibleStart:end]



    var lines []string

    for _, item := range visible {

        lines = append(lines, item.Render())

    }



    return strings.Join(lines, "\n")

}



func (v *VirtualList) ScrollDown() {

    if v.visibleStart+v.visibleCount < len(v.allItems) {

        v.visibleStart++

    }

}



func (v *VirtualList) ScrollUp() {

    if v.visibleStart > 0 {

        v.visibleStart--

    }

}

```



### Efficient Rendering



#### Framerate-Based Renderer



Bubble Tea includes a framerate-based renderer that:

- Batches rapid updates

- Prevents screen flicker

- Optimizes terminal writes



**No configuration needed** - automatic.



#### Avoid Expensive View Calculations



**Bad**:

```go

func (m model) View() string {

    // Recalculates every render

    sorted := sortItems(m.items)

    filtered := filterItems(sorted, m.filter)

    return renderItems(filtered)

}

```



**Good**:

```go

type model struct {

    items         []Item

    displayItems  []Item  // Cached

    filter        string

    dirty         bool

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case filterChangedMsg:

        m.filter = msg.filter

        m.dirty = true  // Mark for recalculation

    }



    if m.dirty {

        m.displayItems = m.recalculateDisplay()

        m.dirty = false

    }



    return m, nil

}



func (m model) View() string {

    return renderItems(m.displayItems)  // Use cached

}

```



#### String Building



Use `strings.Builder` for concatenating many strings:



**Bad**:

```go

result := ""

for _, line := range lines {

    result += line + "\n"  // Many allocations

}

```



**Good**:

```go

var builder strings.Builder

for _, line := range lines {

    builder.WriteString(line)

    builder.WriteRune('\n')

}

result := builder.String()

```



Or use `strings.Join`:

```go

result := strings.Join(lines, "\n")

```



### Memory Management



#### Slice Preallocation



```go

// Bad

var items []Item

for i := 0; i < 1000; i++ {

    items = append(items, Item{})  // Many reallocations

}



// Good

items := make([]Item, 0, 1000)

for i := 0; i < 1000; i++ {

    items = append(items, Item{})  // One allocation

}

```



#### Limit Buffered Content



For log viewers or streams:



```go

type LogModel struct {

    lines     []string

    maxLines  int

}



func (m *LogModel) AddLine(line string) {

    m.lines = append(m.lines, line)



    // Trim old lines

    if len(m.lines) > m.maxLines {

        m.lines = m.lines[len(m.lines)-m.maxLines:]

    }

}

```



#### Avoid Storing Rendered Content



**Bad**:

```go

type model struct {

    items         []Item

    renderedItems []string  // Duplicate memory

}

```



**Good**:

```go

type model struct {

    items []Item

}



func (m model) View() string {

    // Render on demand

    var rendered []string

    for _, item := range m.items {

        rendered = append(rendered, item.View())

    }

    return strings.Join(rendered, "\n")

}

```



### Handling Large Datasets



#### Lazy Loading Pattern



```go

type DataModel struct {

    loaded   []Item

    loader   func(offset, limit int) ([]Item, error)

    offset   int

    pageSize int

    hasMore  bool

}



func (m *DataModel) LoadMore() tea.Cmd {

    return func() tea.Msg {

        items, err := m.loader(m.offset, m.pageSize)

        if err != nil {

            return errorMsg{err}

        }

        return itemsLoadedMsg{items}

    }

}



func (m DataModel) Update(msg tea.Msg) (DataModel, tea.Cmd) {

    switch msg := msg.(type) {

    case itemsLoadedMsg:

        m.loaded = append(m.loaded, msg.items...)

        m.offset += len(msg.items)

        m.hasMore = len(msg.items) == m.pageSize

    case scrolledToBottomMsg:

        if m.hasMore {

            return m, m.LoadMore()

        }

    }

    return m, nil

}

```



#### Pagination Pattern



```go

type PagedModel struct {

    allData     []Item

    currentPage int

    pageSize    int

}



func (m PagedModel) CurrentPageItems() []Item {

    start := m.currentPage * m.pageSize

    end := start + m.pageSize



    if start >= len(m.allData) {

        return []Item{}

    }

    if end > len(m.allData) {

        end = len(m.allData)

    }



    return m.allData[start:end]

}



func (m PagedModel) View() string {

    items := m.CurrentPageItems()

    // Render only current page

}

```



### ANSI Escape Code Optimization



Lipgloss handles ANSI codes efficiently, but be aware:



- **ANSI codes count toward string length** when calculating widths

- Use Lipgloss width functions: `lipgloss.Width(str)`, not `len(str)`

- Viewport has optimizations for ANSI-heavy content



```go

// Bad - incorrect width

width := len(styledString)



// Good - ANSI-aware width

width := lipgloss.Width(styledString)

```



---



## 6. Testing



### 1. teatest (Official Experimental Library)



**Package**: `github.com/charmbracelet/x/exp/teatest`

**Status**: Experimental, work in progress



**Features**:

- Assert entire program output

- Assert partial output

- Assert internal model state

- Golden file testing (snapshot testing)

- Test with simulated terminal

- Send input events to program



**Usage**:



```go

import (

    "testing"

    tea "github.com/charmbracelet/bubbletea"

    "github.com/charmbracelet/x/exp/teatest"

)



func TestApp(t *testing.T) {

    m := NewModel()



    // Create test with simulated terminal

    tm := teatest.NewTestModel(

        t, m,

        teatest.WithInitialTermSize(80, 24),

    )



    // Send key presses

    tm.Send(tea.KeyMsg{

        Type:  tea.KeyRunes,

        Runes: []rune("test input"),

    })

    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})



    // Wait for processing

    teatest.WaitFor(

        t, tm.Output(),

        func(bts []byte) bool {

            return bytes.Contains(bts, []byte("Expected text"))

        },

        teatest.WithCheckInterval(time.Millisecond*100),

        teatest.WithDuration(time.Second*3),

    )



    // Assert final output

    out := tm.FinalOutput(t)

    if !bytes.Contains(out, []byte("Success")) {

        t.Fatal("expected success message")

    }



    // Assert model state

    finalModel := tm.FinalModel(t)

    appModel := finalModel.(AppModel)

    if appModel.completed != true {

        t.Fatal("expected completed state")

    }

}

```



**Golden Files**:

```go

func TestAppGolden(t *testing.T) {

    m := NewModel()

    tm := teatest.NewTestModel(t, m)



    // ... send inputs ...



    // First run captures output

    // Subsequent runs compare against captured

    teatest.RequireEqualOutput(t, tm.Output())

}

```



### 2. catwalk (Third-Party Library)



**GitHub**: https://github.com/knz/catwalk



**Purpose**: Unit test library for Bubble Tea models



**Features**:

- Test Update() with messages

- Compare resulting state

- Compare resulting View output

- No need to run full program



**Usage**:

```go

import "github.com/knz/catwalk"



func TestModelUpdate(t *testing.T) {

    m := NewModel()



    // Apply message

    newModel, cmd := m.Update(tea.KeyMsg{

        Type:  tea.KeyRunes,

        Runes: []rune("x"),

    })



    // Assert state

    if newModel.(Model).input != "x" {

        t.Errorf("expected input 'x', got '%s'", newModel.(Model).input)

    }



    // Assert view

    view := newModel.View()

    if !strings.Contains(view, "x") {

        t.Error("view should contain input")

    }



    // Assert command

    if cmd != nil {

        t.Error("expected no command")

    }

}

```



### 3. VHS - Visual Testing



**GitHub**: https://github.com/charmbracelet/vhs



**Purpose**: Record terminal sessions as GIFs/PNGs from declarative scripts



**Features**:

- Scripted terminal interactions

- Automated screenshot generation

- GIF recording

- Documentation and demo generation



**Example Script** (demo.tape):

```

Output demo.gif



Set FontSize 20

Set Width 1200

Set Height 600



Type "go run main.go"

Enter

Sleep 1s



Type "down"

Sleep 500ms



Type "down"

Sleep 500ms



Type "enter"

Sleep 2s



Screenshot final.png

```



Run: `vhs demo.tape`



**Use Cases**:

- Regression testing (compare screenshots)

- Documentation

- Demos for README

- Bug reproduction



### 4. Testing Patterns



#### Testing Update Logic



```go

func TestUpdateHandlesInput(t *testing.T) {

    m := model{input: ""}



    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}

    newModel, cmd := m.Update(msg)



    m2 := newModel.(model)

    if m2.input != "a" {

        t.Errorf("expected 'a', got '%s'", m2.input)

    }

    if cmd != nil {

        t.Error("expected nil command")

    }

}

```



#### Testing Commands



```go

func TestCommand(t *testing.T) {

    cmd := fetchData("https://api.example.com")



    // Commands return messages

    msg := cmd()



    switch msg := msg.(type) {

    case dataMsg:

        if msg.data == "" {

            t.Error("expected data")

        }

    case errorMsg:

        t.Errorf("unexpected error: %v", msg.err)

    default:

        t.Error("unexpected message type")

    }

}

```



#### Testing View Output



```go

func TestView(t *testing.T) {

    m := model{

        title:   "Test",

        content: "Hello",

    }



    view := m.View()



    if !strings.Contains(view, "Test") {

        t.Error("view should contain title")

    }

    if !strings.Contains(view, "Hello") {

        t.Error("view should contain content")

    }

}

```



#### Testing State Machines



```go

func TestStateTransitions(t *testing.T) {

    tests := []struct {

        name       string

        state      ViewState

        msg        tea.Msg

        wantState  ViewState

    }{

        {

            name:      "loading to list",

            state:     ViewLoading,

            msg:       dataLoadedMsg{},

            wantState: ViewList,

        },

        {

            name:      "list to detail",

            state:     ViewList,

            msg:       itemSelectedMsg{},

            wantState: ViewDetail,

        },

    }



    for _, tt := range tests {

        t.Run(tt.name, func(t *testing.T) {

            m := model{state: tt.state}

            newModel, _ := m.Update(tt.msg)

            m2 := newModel.(model)



            if m2.state != tt.wantState {

                t.Errorf("expected state %v, got %v", tt.wantState, m2.state)

            }

        })

    }

}

```



---



## 7. Code Patterns and Best Practices



### Keyboard Handling



#### Enhanced Keyboard Support (v2)



Bubble Tea v2 adds support for:

- Key releases (for games)

- `shift+enter`, `super+space`, etc.

- Key disambiguation (ctrl+i vs tab)



**v2 Features**:

```go

switch msg := msg.(type) {

case tea.KeyMsg:

    // Distinguish ctrl+i from tab

    if msg.String() == "ctrl+i" {

        // Actual ctrl+i

    }



case tea.KeyReleaseMsg:  // v2 only

    // Handle key up

}



case tea.PasteMsg:  // v2 - dedicated paste message

    text := msg.Text

```



#### Keybinding Patterns



```go

import "github.com/charmbracelet/bubbles/key"



type keyMap struct {

    Up    key.Binding

    Down  key.Binding

    Left  key.Binding

    Right key.Binding

    Enter key.Binding

    Back  key.Binding

    Quit  key.Binding

}



var keys = keyMap{

    Up: key.NewBinding(

        key.WithKeys("up", "k", "w"),

        key.WithHelp("↑/k/w", "move up"),

    ),

    Down: key.NewBinding(

        key.WithKeys("down", "j", "s"),

        key.WithHelp("↓/j/s", "move down"),

    ),

    Enter: key.NewBinding(

        key.WithKeys("enter"),

        key.WithHelp("enter", "select"),

    ),

    Quit: key.NewBinding(

        key.WithKeys("q", "ctrl+c", "esc"),

        key.WithHelp("q", "quit"),

    ),

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        switch {

        case key.Matches(msg, keys.Quit):

            return m, tea.Quit

        case key.Matches(msg, keys.Up):

            m.cursor--

        case key.Matches(msg, keys.Down):

            m.cursor++

        case key.Matches(msg, keys.Enter):

            return m, m.selectItem()

        }

    }

    return m, nil

}

```



### Mouse Handling (v2 Improvements)



**v2 splits MouseMsg into specific types**:



```go

switch msg := msg.(type) {

case tea.MouseClickMsg:

    // Left, right, middle clicks

    x, y := msg.X, msg.Y

    button := msg.Button



case tea.MouseReleaseMsg:

    // Button released



case tea.MouseWheelMsg:

    // Scroll events

    if msg.Up {

        // Scroll up

    } else {

        // Scroll down

    }



case tea.MouseMotionMsg:

    // Mouse movement

    x, y := msg.X, msg.Y

}

```



### Alternate Screen Buffer



**Fullscreen Mode**:

```go

func main() {

    p := tea.NewProgram(

        initialModel(),

        tea.WithAltScreen(),        // Use alternate screen

        tea.WithMouseAllMotion(),   // Enable mouse

    )



    if _, err := p.Run(); err != nil {

        log.Fatal(err)

    }

}

```



**Dynamic Switching** (Enter/Exit fullscreen):

```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        switch msg.String() {

        case "f":

            if m.fullscreen {

                return m, tea.ExitAltScreen

            }

            return m, tea.EnterAltScreen

        }

    }

    return m, nil

}

```



**Inline vs. Fullscreen**:

- **Inline**: Output stays in terminal history

- **Fullscreen**: Uses alternate buffer, clears on exit

- **Mix**: Switch dynamically as needed



### Command Patterns



#### Long-Running Operations



```go

func (m model) startTask() tea.Msg {

    // Simulate long operation

    time.Sleep(5 * time.Second)

    return taskCompleteMsg{result: "done"}

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        if msg.String() == "enter" {

            m.loading = true

            return m, m.startTask

        }

    case taskCompleteMsg:

        m.loading = false

        m.result = msg.result

    }

    return m, nil

}

```



#### HTTP Requests



```go

type httpResponseMsg struct {

    status int

    body   []byte

}



type httpErrorMsg struct {

    err error

}



func fetchURL(url string) tea.Cmd {

    return func() tea.Msg {

        resp, err := http.Get(url)

        if err != nil {

            return httpErrorMsg{err}

        }

        defer resp.Body.Close()



        body, err := io.ReadAll(resp.Body)

        if err != nil {

            return httpErrorMsg{err}

        }



        return httpResponseMsg{

            status: resp.StatusCode,

            body:   body,

        }

    }

}

```



#### Tickers and Intervals



```go

type tickMsg time.Time



func tick() tea.Cmd {

    return tea.Tick(time.Second, func(t time.Time) tea.Msg {

        return tickMsg(t)

    })

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg.(type) {

    case tickMsg:

        m.counter++

        return m, tick()  // Continue ticking

    }

    return m, nil

}



func (m model) Init() tea.Cmd {

    return tick()  // Start ticker

}

```



#### Channel-Based Streams



```go

func listenForEvents(eventChan chan Event) tea.Cmd {

    return func() tea.Msg {

        return eventMsg{event: <-eventChan}

    }

}



func (m model) Init() tea.Cmd {

    return listenForEvents(m.eventChan)

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case eventMsg:

        m.handleEvent(msg.event)

        return m, listenForEvents(m.eventChan)  // Continue listening

    }

    return m, nil

}

```



### Form Building Pattern



```go

type formModel struct {

    inputs   []textinput.Model

    focused  int

    err      error

}



func newForm() formModel {

    inputs := make([]textinput.Model, 3)



    inputs[0] = textinput.New()

    inputs[0].Placeholder = "Name"

    inputs[0].Focus()



    inputs[1] = textinput.New()

    inputs[1].Placeholder = "Email"



    inputs[2] = textinput.New()

    inputs[2].Placeholder = "Phone"



    return formModel{inputs: inputs}

}



func (m formModel) Update(msg tea.Msg) (formModel, tea.Cmd) {

    var cmds []tea.Cmd



    switch msg := msg.(type) {

    case tea.KeyMsg:

        switch msg.String() {

        case "tab", "down":

            m.focused = (m.focused + 1) % len(m.inputs)

            return m, m.updateFocus()



        case "shift+tab", "up":

            m.focused--

            if m.focused < 0 {

                m.focused = len(m.inputs) - 1

            }

            return m, m.updateFocus()



        case "enter":

            return m, m.submitForm()

        }

    }



    // Update focused input

    var cmd tea.Cmd

    m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)

    cmds = append(cmds, cmd)



    return m, tea.Batch(cmds...)

}



func (m formModel) updateFocus() tea.Cmd {

    cmds := make([]tea.Cmd, len(m.inputs))



    for i := range m.inputs {

        if i == m.focused {

            cmds[i] = m.inputs[i].Focus()

        } else {

            m.inputs[i].Blur()

        }

    }



    return tea.Batch(cmds...)

}



func (m formModel) View() string {

    var views []string

    for i, input := range m.inputs {

        cursor := " "

        if i == m.focused {

            cursor = ">"

        }

        views = append(views, cursor+" "+input.View())

    }



    return lipgloss.JoinVertical(lipgloss.Left, views...)

}

```



### Gotchas and Common Mistakes



#### 1. Forgetting to Return Commands



**Bad**:

```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        if msg.String() == "r" {

            fetchData()  // ❌ Command not returned!

        }

    }

    return m, nil

}

```



**Good**:

```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        if msg.String() == "r" {

            return m, fetchData()  // ✅ Return the command

        }

    }

    return m, nil

}

```



#### 2. Blocking Operations in Update



**Bad**:

```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        if msg.String() == "r" {

            time.Sleep(5 * time.Second)  // ❌ Blocks UI!

            m.data = fetchData()

        }

    }

    return m, nil

}

```



**Good**:

```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        if msg.String() == "r" {

            return m, func() tea.Msg {  // ✅ Async command

                time.Sleep(5 * time.Second)

                return dataMsg{fetchData()}

            }

        }

    case dataMsg:

        m.data = msg.data

    }

    return m, nil

}

```



#### 3. Not Propagating WindowSizeMsg



**Bad**:

```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    // ❌ Child components don't know terminal size

    m.list, _ = m.list.Update(msg)

    return m, nil

}

```



**Good**:

```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:

        m.width = msg.Width

        m.height = msg.Height

        // ✅ Resize children

        m.list.SetSize(msg.Width, msg.Height-2)

    }



    m.list, _ = m.list.Update(msg)

    return m, nil

}

```



#### 4. Using len() Instead of lipgloss.Width()



**Bad**:

```go

width := len(styledString)  // ❌ Counts ANSI escape codes

padding := totalWidth - width

```



**Good**:

```go

width := lipgloss.Width(styledString)  // ✅ ANSI-aware

padding := totalWidth - width

```



#### 5. Not Handling All Message Types



**Bad**:

```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        // Only handles keys, ignores WindowSize, Mouse, etc.

    }

    return m, nil  // ❌ Child components never get messages

}

```



**Good**:

```go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    var cmds []tea.Cmd



    // Handle specific messages

    switch msg := msg.(type) {

    case tea.KeyMsg:

        // Handle keys

    case tea.WindowSizeMsg:

        // Handle resize

    }



    // ✅ Always propagate to children

    var cmd tea.Cmd

    m.child, cmd = m.child.Update(msg)

    cmds = append(cmds, cmd)



    return m, tea.Batch(cmds...)

}

```



---



## 8. Notable Example Applications



### Production Applications Using Bubble Tea



**Official Charm Bracelet Apps**:

- **Glow** - Markdown reader and terminal-based markdown viewer

  - GitHub: https://github.com/charmbracelet/glow

- **Soft Serve** - Self-hosted Git server with TUI

  - GitHub: https://github.com/charmbracelet/soft-serve

- **VHS** - Terminal recorder/screenshot tool

  - GitHub: https://github.com/charmbracelet/vhs



**Major Third-Party Apps** (8,000+ total):

- **chezmoi** - Dotfile manager across multiple machines

  - Used by thousands of developers

- **TruffleHog** (Truffle Security Co.) - Find leaked credentials

- **container-canary** (NVIDIA) - Container validator

- **eks-node-viewer** (AWS) - Visualize EKS cluster node usage

- **mc** (MinIO) - Official MinIO client



### Finding More Examples



**Charm & Friends**: https://github.com/charm-and-friends/additional-bubbles

Community-maintained bubbles and components



**Example Directory**: https://github.com/charmbracelet/bubbletea/tree/main/examples



Notable examples in the repo:

- `list-simple` - Interactive list selection

- `textinputs` - Multiple inputs with focus switching

- `credit-card-form` - Multi-step form with validation

- `progress-download` - Download progress tracking

- `fullscreen` - Alternate screen buffer usage

- `realtime` - Channel-based event handling

- `spinner` - Loading indicators

- `table` - Tabular data display



### Learning Resources



**Official**:

- Tutorial: https://github.com/charmbracelet/bubbletea/tree/main/tutorials/basics

- Commands Tutorial: https://github.com/charmbracelet/bubbletea/tree/main/tutorials/commands

- Go Docs: https://pkg.go.dev/github.com/charmbracelet/bubbletea



**Blog Posts**:

- "Commands in Bubble Tea": https://charm.land/blog/commands-in-bubbletea/

- "Tips for Building Bubble Tea Programs": https://leg100.github.io/en/posts/building-bubbletea-programs/

- "The Bubbletea State Machine Pattern": https://zackproser.com/blog/bubbletea-state-machine



**Video Tutorials**: Available on YouTube and mentioned in official docs



---



## Quick Reference



### Installation Commands



```bash

# Bubble Tea v1 (stable)

go get github.com/charmbracelet/bubbletea



# Bubble Tea v2 (beta)

go get github.com/charmbracelet/bubbletea/v2@v2.0.0-beta.3



# Lipgloss v1

go get github.com/charmbracelet/lipgloss



# Lipgloss v2 (alpha)

go get github.com/charmbracelet/lipgloss/v2@v2.0.0-alpha.2



# Bubbles

go get github.com/charmbracelet/bubbles



# Third-party

go get github.com/evertras/bubble-table

go get github.com/erikgeiser/promptkit

go get github.com/charmbracelet/harmonica

go get github.com/NimbleMarkets/ntcharts

go get github.com/alecthomas/chroma/v2

go get github.com/lrstanley/bubblezone

```



### Minimal Application Template



```go

package main



import (

    "fmt"

    "os"



    tea "github.com/charmbracelet/bubbletea"

)



type model struct {

    // Your state here

}



func initialModel() model {

    return model{}

}



func (m model) Init() tea.Cmd {

    return nil

}



func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        switch msg.String() {

        case "ctrl+c", "q":

            return m, tea.Quit

        }

    }

    return m, nil

}



func (m model) View() string {

    return "Hello, World!\n\nPress q to quit.\n"

}



func main() {

    p := tea.NewProgram(initialModel())

    if _, err := p.Run(); err != nil {

        fmt.Printf("Error: %v", err)

        os.Exit(1)

    }

}

```



---



## Sources



### Core Framework

- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea)

- [Bubble Tea v1 Docs](https://pkg.go.dev/github.com/charmbracelet/bubbletea)

- [Bubble Tea v2 Docs](https://pkg.go.dev/github.com/charmbracelet/bubbletea/v2)

- [Bubble Tea Tutorials](https://github.com/charmbracelet/bubbletea/tree/main/tutorials/basics)

- [Commands in Bubble Tea](https://charm.land/blog/commands-in-bubbletea/)



### Styling

- [Lipgloss GitHub](https://github.com/charmbracelet/lipgloss)

- [Lipgloss v1 Docs](https://pkg.go.dev/github.com/charmbracelet/lipgloss)

- [Lipgloss v2 Docs](https://pkg.go.dev/github.com/charmbracelet/lipgloss/v2)

- [Lipgloss Table Package](https://pkg.go.dev/github.com/charmbracelet/lipgloss/table)



### Components

- [Bubbles GitHub](https://github.com/charmbracelet/bubbles)

- [Bubbles Docs](https://pkg.go.dev/github.com/charmbracelet/bubbles)

- [Viewport Package](https://pkg.go.dev/github.com/charmbracelet/bubbles/viewport)

- [List Package](https://pkg.go.dev/github.com/charmbracelet/bubbles/list)

- [Table Package](https://pkg.go.dev/github.com/charmbracelet/bubbles/table)

- [Spinner Package](https://pkg.go.dev/github.com/charmbracelet/bubbles/spinner)

- [TextInput Package](https://pkg.go.dev/github.com/charmbracelet/bubbles/textinput)

- [Key Package](https://pkg.go.dev/github.com/charmbracelet/bubbles/key)

- [Help Package](https://pkg.go.dev/github.com/charmbracelet/bubbles/help)



### Third-Party Libraries

- [bubble-table GitHub](https://github.com/Evertras/bubble-table)

- [promptkit GitHub](https://github.com/erikgeiser/promptkit)

- [harmonica GitHub](https://github.com/charmbracelet/harmonica)

- [ntcharts GitHub](https://github.com/NimbleMarkets/ntcharts)

- [chroma GitHub](https://github.com/alecthomas/chroma)

- [bubblezone GitHub](https://github.com/lrstanley/bubblezone)



### Testing

- [catwalk GitHub](https://github.com/knz/catwalk)

- [VHS GitHub](https://github.com/charmbracelet/vhs)

- [Writing Bubble Tea Tests](https://carlosbecker.com/posts/teatest/)



### Best Practices & Tutorials

- [Tips for Building Bubble Tea Programs](https://leg100.github.io/en/posts/building-bubbletea-programs/)

- [The Bubbletea State Machine Pattern](https://zackproser.com/blog/bubbletea-state-machine)

- [Building Terminal IRC Client with Bubble Tea](https://sngeth.com/go/terminal/ui/bubble-tea/2025/08/17/building-terminal-ui-with-bubble-tea/)

- [Rapidly Building Interactive CLIs](https://www.inngest.com/blog/interactive-clis-with-bubbletea)

- [Intro to Bubble Tea in Go](https://dev.to/andyhaskell/intro-to-bubble-tea-in-go-21lg)



### Community

- [Charm & Friends](https://github.com/charm-and-friends/additional-bubbles)

- [Bubble Tea Discussions](https://github.com/charmbracelet/bubbletea/discussions)

- [Bubble Tea Releases](https://github.com/charmbracelet/bubbletea/releases)
