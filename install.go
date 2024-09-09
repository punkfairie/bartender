package main

import (
	"bufio"
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

		m.logger.Infof("Installing %s...", pkg)

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
				m.logger.Infof("Finished installing %s!", pkg)
				m.logger.Infof("Output: %s\n\n", *m.output)
				time.Sleep(1 * time.Second)
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
