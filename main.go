package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

var width, height, _ = term.GetSize(int(os.Stdout.Fd()))

type menu struct {
	order      SoftwarePackages
	current    int
	done       int
	keys       keyMap
	help       help.Model
	inputStyle gloss.Style
}

func initialModel() menu {
	return menu{
		current: 0,
		keys:    keys,
		help:    help.New(),
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

type YamlStructure struct {
	InstallerPreferences InstallerPreferences `yaml:"installerPreferences"`
	SoftwarePackages     SoftwarePackages     `yaml:"softwarePackages"`
}

type InstallerPreferences struct {
	Apt     []string `yaml:"apt"`
	Darwin  []string `yaml:"darwin"`
	Dnf     []string `yaml:"dnf"`
	Freebsd []string `yaml:"freebsd"`
	Pacman  []string `yaml:"pacman"`
	Ubuntu  []string `yaml:"ubuntu"`
	Windows []string `yaml:"windows"`
	Zypper  []string `yaml:"zypper"`
}

type SoftwarePackages map[string]SoftwareDef

type SoftwareDef struct {
	Bin        *string   `yaml:"_bin"`
	Desc       string    `yaml:"_desc"`
	Deps       *[]string `yaml:"_deps"`
	Docs       *string   `yaml:"_docs"`
	Github     *string   `yaml:"_github"`
	Home       *string   `yaml:"_home"`
	Name       string    `yaml:"_name"`
	Apk        *string   `yaml:"apk"`
	Appimage   *string   `yaml:"appimage"`
	Basher     *string   `yaml:"basher"`
	Binary     *osNames  `yaml:"binary"`
	Bpkg       *string   `yaml:"bpkg"`
	Brew       *string   `yaml:"brew"`
	Cargo      *string   `yaml:"cargo"`
	Cask       *string   `yaml:"cask"`
	Crew       *string   `yaml:"crew"`
	Choco      *string   `yaml:"choco"`
	Dnf        *string   `yaml:"dnf"`
	Flatpak    *string   `yaml:"flatpak"`
	Gem        *string   `yaml:"gem"`
	Go         *string   `yaml:"go"`
	Krew       *string   `yaml:"krew"`
	Nix        *string   `yaml:"nix"`
	Npm        *string   `yaml:"npm"`
	Pacman     *string   `yaml:"pacman"`
	Pipx       *string   `yaml:"pipx"`
	PkgFreebsd *string   `yaml:"pkg-freebsd"`
	PkgTermux  *string   `yaml:"pkg-termux"`
	Port       *string   `yaml:"port"`
	Scoop      *string   `yaml:"scoop"`
	Script     *string   `yaml:"string"`
	Snap       *string   `yaml:"snap"`
	Whalebrew  *string   `yaml:"whalebrew"`
	Winget     *string   `yaml:"winget"`
	Xbps       *string   `yaml:"xbps"`
	Yay        *string   `yaml:"yay"`
	Zypper     *string   `yaml:"zypper"`
}

type osNames struct {
	Darwin  *string `yaml:"darwin"`
	Linux   *string `yaml:"linux"`
	Windows *string `yaml:"windows"`
}

type yamlMsg YamlStructure

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func readYaml() tea.Msg {
	fileData, fileErr := os.ReadFile("/Users/marley/hackin/install.fairie/software-custom.yml")
	if fileErr != nil {
		return errMsg{fileErr}
	}

	var parsedYaml YamlStructure

	yamlErr := yaml.Unmarshal(fileData, &parsedYaml)
	if yamlErr != nil {
		return errMsg{yamlErr}
	}

	return yamlMsg(parsedYaml)
}

func (m menu) Init() tea.Cmd {
	return readYaml
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case yamlMsg:
		m.order = msg.SoftwarePackages
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
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
	sidebarContent := ""

	for _, item := range m.order {
		sidebarContent += fmt.Sprintf("%s\n", item.Name)
	}

	main := mainStyle.Render(mainContent)
	sidebar := sidebarStyle.Render(sidebarContent)

	content := gloss.JoinHorizontal(gloss.Top, main, sidebar)

	helpView := m.help.View(m.keys)

	page := gloss.JoinVertical(gloss.Left, content, helpView)

	return gloss.PlaceHorizontal(width, gloss.Center, page)
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)

		os.Exit(1)
	}
}
