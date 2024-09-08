package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"golang.org/x/term"
)

var width, height, _ = term.GetSize(int(os.Stdout.Fd()))

type menu struct {
	order      SoftwarePackages
	current    int
	keys       keyMap
	help       help.Model
	inputStyle gloss.Style
	spinner    spinner.Model
	quitting   bool
}

const softwareInstructionsFile = "/Users/marley/hackin/install.fairie/software-custom.yml"

func initialModel() menu {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = gloss.NewStyle().Foreground(gloss.Color("3"))

	return menu{
		current:  0,
		keys:     keys,
		help:     help.New(),
		spinner:  s,
		quitting: false,
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

type yamlMsg YamlStructure

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m menu) Init() tea.Cmd {
	return tea.Batch(readYaml(softwareInstructionsFile), m.spinner.Tick)
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case yamlMsg:
		m.order = msg.SoftwarePackages
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m menu) View() string {
	borderStyle := gloss.NewStyle().
		BorderStyle(gloss.RoundedBorder()).
		BorderForeground(gloss.Color("5")).
		Padding(0, 1)

	mainStyle := borderStyle.
		Width(int(float64(width) * 0.65))

	sidebarStyle := borderStyle.
		Width(int(float64(width) * 0.3))

	mainContent := ""

	softwareListEnumerator := func(l list.Items, i int) string {
		if m.current == i {
			return m.spinner.View()
		} else if m.current > i {
			return ""
		}
		return ""
	}

	software := list.New().Enumerator(softwareListEnumerator)

	keys := sortMapKeys(m.order)

	for _, k := range keys {
		software.Item(m.order[k].Name)
	}

	sidebarContent := software.String()

	main := mainStyle.Render(mainContent)
	sidebar := sidebarStyle.Render(sidebarContent)

	content := gloss.JoinHorizontal(gloss.Top, main, sidebar)

	helpView := m.help.View(m.keys)

	last := ""
	if m.quitting {
		last = "\n"
	}

	page := gloss.JoinVertical(gloss.Left, content, helpView, last)

	return gloss.PlaceHorizontal(width, gloss.Center, page)
}

func main() {
	p := tea.NewProgram(initialModel())

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
