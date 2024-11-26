package ui

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Brian-Kariu/ryuk/cmd/config"
)

type welcomeModel struct {
	message    string
	showPrompt bool
	createEnv  bool
}

// Check if .ryuk.yaml exists
func checkConfigFile() bool {
	config.InitConstants()
	fullpath := filepath.Join(config.BasePath, ".ryuk.yaml")
	_, err := os.Stat(fullpath)
	return !os.IsNotExist(err)
}

// Initial model setup
func initialWelcomeModel() welcomeModel {
	message := `
██████╗ ██╗   ██╗██╗   ██╗██╗  ██╗
██╔══██╗╚██╗ ██╔╝██║   ██║██║ ██╔╝
██████╔╝ ╚████╔╝ ██║   ██║█████╔╝ 
██╔══██╗  ╚██╔╝  ██║   ██║██╔═██╗ 
██║  ██║   ██║   ╚██████╔╝██║  ██╗
╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝

Welcome to Ryuk CLI!
`
	showPrompt := !checkConfigFile()

	return welcomeModel{
		message:    message,
		showPrompt: showPrompt,
	}
}

func (m welcomeModel) Init() tea.Cmd {
	return nil
}

func (m welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			if m.showPrompt {
				return m, func() tea.Msg {
					StartWorkspaceView()
					return nil
				}
			}
		case "enter":
			if !m.showPrompt {
				return m, tea.Quit
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m welcomeModel) View() string {
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#34D399"))

	if m.showPrompt {
		return fmt.Sprintf(
			"%s\n%s\n\n%s",
			style.Render(m.message),
			"Configuration file `.ryuk.yaml` not found!",
			"Press 'y' to create a new environment, or 'q' to quit.",
		)
	}

	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		style.Render(m.message),
		"Configuration file detected!",
		"Press 'Enter' to continue or 'q' to quit.",
	)
}

func StartWelcome() {
	p := tea.NewProgram(initialWelcomeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting welcome view: %v\n", err)
		os.Exit(1)
	}
}
