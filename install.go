package main

import (
	"bufio"
	"io"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type cmdMsg string

type cmdDoneMsg struct{}

func installPackage(sub chan string) tea.Cmd {
	return func() tea.Msg {
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
				return cmdDoneMsg{}
			}

			if err != nil {
				return errMsg{err}
			}

			sub <- string(line)
		}
	}
}

func waitForCmdResponses(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return cmdMsg(<-sub)
	}
}
