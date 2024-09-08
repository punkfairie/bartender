package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"golang.org/x/term"
)

type menu struct {
	order      []string
	current    int
	keys       keyMap
	help       help.Model
	inputStyle gloss.Style
	spinner    spinner.Model
	quitting   bool
	viewport   viewport.Model
	width      int
	height     int
}

const (
	softwareInstructionsFile = "/Users/marley/hackin/install.fairie/home/.chezmoidata.yaml"
	softwareGroup            = "_Full-Desktop"
)

func initialModel() menu {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = gloss.NewStyle().Foreground(gloss.Color("3"))
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	return menu{
		current:  150,
		keys:     keys,
		help:     help.New(),
		spinner:  s,
		quitting: false,
		width:    width,
		height:   height,
	}
}

type keyMap struct {
	Quit key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	keys := k.ShortHelp()

	return [][]key.Binding{
		keys,
	}
}

type softwareListMsg []string

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m menu) Init() tea.Cmd {
	return tea.Batch(getSoftwareList(softwareInstructionsFile), m.spinner.Tick)
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case softwareListMsg:
		m.order = msg

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.width, m.height, _ = term.GetSize(int(os.Stdout.Fd()))
	}

	return m, tea.Batch(cmds...)
}

func (m menu) View() string {
	borderStyle := gloss.NewStyle().
		BorderStyle(gloss.RoundedBorder()).
		BorderForeground(gloss.Color("5")).
		Padding(0, 1)

	mainStyle := borderStyle.
		Width(int(float64(m.width) * 0.65))

	sidebarStyle := borderStyle.
		Width(int(float64(m.width) * 0.3))

	topPadding := 1

	mainContent := ""

	helpView := m.help.View(m.keys)

	softwareListEnumerator := func(l list.Items, i int) string {
		if m.current == i {
			return m.spinner.View()
		} else if m.current > i {
			return ""
		}
		return ""
	}

	software := list.New().Enumerator(softwareListEnumerator)

	sidebarHeight := m.height - 3 - gloss.Height(helpView) - topPadding

	if len(m.order) > 0 {
		start := max(m.current-10, 0)
		end := min(start+sidebarHeight, len(m.order))

		if (end - start) < sidebarHeight {
			start = (len(m.order) - sidebarHeight)
		}

		for _, item := range m.order[start:end] {
			software.Item(item)
		}
	}

	sidebarContent := software.String()

	main := mainStyle.Render(mainContent)
	sidebar := sidebarStyle.Render(sidebarContent)

	content := gloss.JoinHorizontal(gloss.Top, main, sidebar)

	top := strings.Repeat("\n", topPadding)
	last := ""
	if m.quitting {
		last = "\n"
	}

	page := gloss.JoinVertical(gloss.Left, top, content, helpView, last)

	return gloss.PlaceHorizontal(m.width, gloss.Center, page)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)

		os.Exit(1)
	}
}

func sortMapKeys(m SoftwarePackages) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
