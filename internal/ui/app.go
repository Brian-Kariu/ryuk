package ui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const (
	welcomeView   = "welcome"
	workspaceView = "workspace"
)

// Styles
var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211")).Bold(true)
	infoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Align(lipgloss.Left)
	formStyle   = lipgloss.NewStyle().MarginLeft(4)
)

// Model for the app
type model struct {
	currentView string
	input       [3]string // [0]: workspace name, [1]: environment, [2]: description
	focusIndex  int
	quitting    bool
	message     string
	showPrompt  bool
	createEnv   bool
}

// Initial model
func initialModel() model {
	return model{
		currentView: welcomeView,
	}
}

// Init initializes the program
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.currentView {
		case welcomeView:
			welcome := InitialWelcomeModel()
			return welcome.Update(msg)
			// case workspaceView:
			// 	workspace := InitialWorkspaceModel()
			// 	return workspace.Update(msg)
		}
	case tea.WindowSizeMsg:
		// Handle resizing here if needed
	}
	return m, nil
}

// View renders the UI based on the current view
func (m model) View() string {
	switch m.currentView {
	case welcomeView:
		welcome := InitialWelcomeModel()
		return welcome.View()
	// case workspaceView:
	// 	workspace := InitialWorkspaceModel()
	// 	return workspace.View()
	default:
		return "Unknown view"
	}
}

func StartApp() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Info("Error running program: %v\n", err)
		os.Exit(1)
	}
}
