package tui

import (
    "fmt"
    "time"

    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"

    "tache/internal/discovery"
    "tache/internal/types"
)

type item types.Session

func (i item) Title() string       { return fmt.Sprintf("%s", i.SuffixOrNone()) }
func (i item) Description() string { return fmt.Sprintf("%s  pid=%d  %s", i.Prefix, i.PID, i.Socket) }
func (i item) FilterValue() string { return fmt.Sprintf("%s %s %s %s", i.Suffix, i.Prefix, i.Command, i.Socket) }

func (i item) SuffixOrNone() string {
    if i.Suffix == "" {
        return "[none]"
    }
    return i.Suffix
}

type model struct {
    list   list.Model
    status string
}

var (
    titleStyle = lipgloss.NewStyle().Bold(true)
)

func initialModel() model {
    l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
    l.Title = titleStyle.Render("Taché — dtach sessions")
    l.SetShowHelp(true)
    l.SetFilteringEnabled(true)
    l.SetShowStatusBar(false)
    l.SetShowPagination(true)
    l.SetShowTitle(true)
    return model{list: l}
}

func (m model) Init() tea.Cmd { return loadSessions() }

func loadSessions() tea.Cmd {
    return func() tea.Msg {
        time.Sleep(10 * time.Millisecond)
        sessions, err := discovery.Discover(nil)
        if err != nil {
            return errMsg{err}
        }
        items := make([]list.Item, 0, len(sessions))
        for _, s := range sessions {
            items = append(items, item(s))
        }
        return itemsMsg(items)
    }
}

type itemsMsg []list.Item
type errMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.list.SetSize(msg.Width, msg.Height-2)
    case itemsMsg:
        m.list.SetItems(msg)
        m.status = fmt.Sprintf("Loaded %d session(s)", len(msg))
    case errMsg:
        m.status = fmt.Sprintf("Error: %v", msg.err)
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        }
    }
    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func (m model) View() string {
    if m.status != "" {
        return m.list.View() + "\n" + m.status + "\n"
    }
    return m.list.View()
}

// Run starts the TUI program.
func Run() error {
    p := tea.NewProgram(initialModel(), tea.WithAltScreen())
    _, err := p.Run()
    return err
}

