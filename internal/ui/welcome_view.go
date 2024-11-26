package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type welcomeModel struct {
	message    string
	showPrompt bool
	createEnv  bool
	Quitting   bool
}

// Initial model setup
func InitialWelcomeModel() welcomeModel {
	message := `
██████╗ ██╗   ██╗██╗   ██╗██╗  ██╗
██╔══██╗╚██╗ ██╔╝██║   ██║██║ ██╔╝
██████╔╝ ╚████╔╝ ██║   ██║█████╔╝ 
██╔══██╗  ╚██╔╝  ██║   ██║██╔═██╗ 
██║  ██║   ██║   ╚██████╔╝██║  ██╗
╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝

Welcome to Ryuk CLI!
`
	return welcomeModel{
		message: message,
	}
}

func (m welcomeModel) Init() tea.Cmd {
	return nil
}

func (m welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
		// if k == "y" {
		// 	workspace := InitialWorkspaceModel()
		// 	return workspace.Update(msg)
		// }
	}

	welcome := InitialWelcomeModel()
	return welcome.Update(msg)
}

func (m welcomeModel) View() string {
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	return m.message
}
