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
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"golang.org/x/term"
)

type menu struct {
	order    []string
	recipes  SoftwarePackages
	current  int
	keys     keyMap
	help     help.Model
	spinner  spinner.Model
	quitting bool
	width    int
	height   int
	sub      chan string
	output   *string
	viewport viewport.Model
	logger   *log.Logger
}

const (
	ordersFile    = "/Users/marley/hackin/install.fairie/home/.chezmoidata.yaml"
	recipesFile   = "/Users/marley/hackin/install.fairie/software.yml"
	softwareGroup = "_Full-Desktop"
)

func initialModel(logFile *os.File) menu {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	logger := log.New(logFile)
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(log.TextFormatter)

	m := menu{
		current:  0,
		sub:      make(chan string),
		output:   new(string),
		viewport: viewport.New(0, 30),
		width:    width,
		height:   height,
		spinner:  s,
		keys:     keys,
		help:     help.New(),
		logger:   logger,
		quitting: false,
	}

	return m
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

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m *menu) appendOutput(s string) {
	*m.output += "\n" + s
	m.viewport.SetContent(*m.output)
	m.viewport.GotoBottom()
}

func (m menu) Init() tea.Cmd {
	return tea.Batch(getOrders(ordersFile), getRecipes(recipesFile), m.spinner.Tick)
}

func (m menu) setDimensions() {
	m.width, m.height, _ = term.GetSize(int(os.Stdout.Fd()))
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case ordersMsg:
		m.order = msg
		if len(m.recipes) > 0 {
			cmds = append(cmds, m.installPackage(), waitForCmdResponses(m.sub))
		}

	case recipesMsg:
		m.recipes = SoftwarePackages(msg)
		if len(m.recipes) > 0 {
			cmds = append(cmds, m.installPackage(), waitForCmdResponses(m.sub))
		}

	case cmdMsg:
		m.appendOutput(string(msg))
		cmds = append(cmds, waitForCmdResponses(m.sub))

	case cmdDoneMsg:
		m.current++
		m.output = new(string)
		cmds = append(cmds, m.installPackage(), waitForCmdResponses(m.sub))

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
		m.setDimensions()
		m.viewport.Width = lipgloss.Width(m.mainView())
		m.viewport.Height = lipgloss.Height(m.mainView())

	case errMsg:
		m.logger.Error("Error: " + msg.Error())
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m menu) View() string {
	content := lipgloss.JoinHorizontal(lipgloss.Top, m.mainView(), m.sidebarView())

	top := strings.Repeat("\n", topPadding)
	last := ""
	if m.quitting {
		last = "\n"
	}

	page := lipgloss.JoinVertical(lipgloss.Left, top, content, m.helpView(), last)

	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, page)
}

func main() {
	l, err := os.Create("log.txt")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer l.Close()

	f, err := tea.LogToFile("tea-log.txt", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(
		initialModel(l),
		tea.WithAltScreen(),
	)

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
