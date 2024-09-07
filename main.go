package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type menu struct {
	order   []string
	current int
	done    int
}

func initialModel() menu {
	return menu{
		order: []string{"brew install yq"},
	}
}

func (m menu) Init() tea.Cmd {
	return nil
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m menu) View() string {
	s := "Installing...\n\n"

	for _, item := range m.order {
		s += fmt.Sprintf("%s\n", item)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)

		os.Exit(1)
	}
}
