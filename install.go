package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type cmdMsg string

type cmdDoneMsg struct{}

func (m menu) installPackage() tea.Cmd {
	return func() tea.Msg {
		pkg := m.order[m.current]

		m.appendOutput(fmt.Sprintf("Installing %s", pkg))

		cmd := exec.Command("./test.sh")
		out, err := cmd.StdoutPipe()
		if err != nil {
			return errMsg{err}
		}

		if err := cmd.Start(); err != nil {
			return errMsg{err}
		}

		buf := bufio.NewReader(out)
		for {
			line, _, err := buf.ReadLine()

			if err == io.EOF {
				m.appendOutput(fmt.Sprintf("Finished installing %s!", pkg))
				time.Sleep(3 * time.Second)
				return cmdDoneMsg{}
			}

			if err != nil {
				return errMsg{err}
			}

			m.sub <- string(line)
		}
	}
}

func waitForCmdResponses(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return cmdMsg(<-sub)
	}
}
