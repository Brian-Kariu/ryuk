package ui

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Brian-Kariu/ryuk/cmd"
	"github.com/Brian-Kariu/ryuk/cmd/config"
)

const (
	welcomeView   = "welcome"
	workspaceView = "workspace"
)

// Styles
var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211")).Bold(true)
	infoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	formStyle   = lipgloss.NewStyle().MarginLeft(4)
)

// Model for the app
type model struct {
	currentView string
	input       [3]string // [0]: workspace name, [1]: environment, [2]: description
	focusIndex  int
	quitting    bool
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
			return updateWelcome(msg, m)
		case workspaceView:
			return updateWorkspace(msg, m)
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
		return viewWelcome(m)
	case workspaceView:
		return viewWorkspace(m)
	default:
		return "Unknown view"
	}
}

// --- Welcome View Logic ---
func updateWelcome(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.currentView = workspaceView
		return m, nil
	case "n", "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func viewWelcome(m model) string {
	message := `
██████╗ ██╗   ██╗██╗   ██╗██╗  ██╗
██╔══██╗╚██╗ ██╔╝██║   ██║██║ ██╔╝
██████╔╝ ╚████╔╝ ██║   ██║█████╔╝ 
██╔══██╗  ╚██╔╝  ██║   ██║██╔═██╗ 
██║  ██║   ██║   ╚██████╔╝██║  ██╗
╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝

Welcome to Ryuk CLI!
`
	b := headerStyle.Render(message)
	b += "\n"
	b += infoStyle.Render("Checking for .ryuk.yaml...\n")

	config.InitConstants()
	fullpath := filepath.Join(config.BasePath, ".ryuk.yaml")
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		b += infoStyle.Render("No config found. Do you want to create a new workspace? (y/n)\n")
	} else {
		b += infoStyle.Render(".ryuk.yaml found! Press any key to continue.\n")
	}
	return fmt.Sprintf(b)
}

// Helper to execute a Cobra command programmatically
func ExecuteCommand(args ...string) error {
	root := cmd.RootCmd // Access the exported root command
	root.SetArgs(args)
	return root.Execute()
}

// --- Workspace View Logic ---
func updateWorkspace(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.focusIndex = (m.focusIndex + 1) % len(m.input)
	case "enter":
		if m.focusIndex == len(m.input)-1 {
			args := []string{
				"workspace", "create",
				m.input[0],
				"-e", m.input[1],
			}
			err := ExecuteCommand(args...)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return m, nil
			}

			m.quitting = true
			return m, tea.Quit
		}
	case "backspace":
		if len(m.input[m.focusIndex]) > 0 {
			m.input[m.focusIndex] = m.input[m.focusIndex][:len(m.input[m.focusIndex])-1]
		}
	default:
		m.input[m.focusIndex] += msg.String()
	}
	return m, nil
}

func viewWorkspace(m model) string {
	form := fmt.Sprintf(
		"Workspace Name: [%s]\nEnvironment (prod/staging): [%s]\nDescription: [%s]",
		formInputStyle(m.input[0], m.focusIndex == 0),
		formInputStyle(m.input[1], m.focusIndex == 1),
		formInputStyle(m.input[2], m.focusIndex == 2),
	)
	return formStyle.Render("Create a New Workspace\n\n" + form + "\n\n[tab: next field] [enter: submit] [backspace: delete]")
}

func formInputStyle(input string, focused bool) string {
	if focused {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(input)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(input)
}

// --- Main Function ---
func StartApp() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
