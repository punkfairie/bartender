package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

var borderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("5")).
	Padding(0, 1)

var topPadding = 1

func (m menu) mainView() string {
	mainWidth := int(float64(m.width) * 0.65)

	mainStyle := borderStyle.
		Width(mainWidth).
		Height(m.calcInnerSidebarHeight() - 2)

	mainContent := m.viewport.View()

	return mainStyle.Render(mainContent)
}

func (m menu) calcInnerSidebarHeight() int {
	return m.height - 3 - lipgloss.Height(m.helpView()) - topPadding
}

func (m menu) sidebarView() string {
	sidebarStyle := borderStyle.
		Width(int(float64(m.width) * 0.3))

	softwareListEnumerator := func(l list.Items, i int) string {
		if m.current == i {
			return m.spinner.View()
		} else if m.current > i {
			return ""
		}
		return ""
	}

	software := list.New().Enumerator(softwareListEnumerator)

	sidebarHeight := m.calcInnerSidebarHeight()

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

	return sidebarStyle.Render(sidebarContent)
}

func (m menu) helpView() string {
	return m.help.View(m.keys)
}
