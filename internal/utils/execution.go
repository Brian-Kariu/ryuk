package utils

import (
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type ExecuteFinishedMsg struct{ err error }

func ExecuteCommand(args []string) tea.Cmd {
	log.Info("Running command: %s\n", "command", strings.Join(args, " "))
	c := exec.Command("go run main.go", args...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return ExecuteFinishedMsg{err}
	})
}
