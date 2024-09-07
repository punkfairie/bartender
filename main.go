package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

var width, height, _ = term.GetSize(int(os.Stdout.Fd()))

type menu struct {
	order   SoftwarePackages
	current int
	done    int
}

func initialModel() menu {
	return menu{
		current: 0,
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
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m menu) View() string {
	mainStyle := gloss.NewStyle().
		Width(int(float64(width) * 0.65)).
		BorderStyle(gloss.NormalBorder()).
		BorderForeground(gloss.Color("63"))

	sidebarStyle := gloss.NewStyle().
		Width(int(float64(width) * 0.3)).
		BorderStyle(gloss.NormalBorder()).
		BorderForeground(gloss.Color("63"))

	mainContent := fmt.Sprintln(m.order)
	sidebarContent := ""

	main := mainStyle.Render(mainContent)
	sidebar := sidebarStyle.Render(sidebarContent)

	content := gloss.JoinHorizontal(gloss.Top, main, sidebar)

	return gloss.PlaceHorizontal(width, gloss.Center, content)
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)

		os.Exit(1)
	}
}
