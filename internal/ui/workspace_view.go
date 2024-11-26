package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func createWorkspace(name, env, description string) error {
	// cmd := exec.Command(".ryuk", "workspace", "create", name, "-e", env)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	//
	// log.Info("Running command: %s\n", "command", cmd.String())
	// return cmd.Run()
	log.Info("Workspace created")
	return nil
}

// WorkspaceModel defines the state of the form
type workspaceModel struct {
	workspaceName textinput.Model
	env           textinput.Model
	description   textinput.Model
	focusIndex    int
	submitted     bool
	err           error
}

// Initialize the form inputs
func initialWorkspaceModel() workspaceModel {
	w := workspaceModel{}

	// Workspace Name Input
	w.workspaceName = textinput.New()
	w.workspaceName.Placeholder = "Workspace Name"
	w.workspaceName.Focus()

	w.workspaceName.PromptStyle = focusedStyle
	w.workspaceName.TextStyle = focusedStyle

	// Environment Input
	w.env = textinput.New()
	w.env.Placeholder = "Environment (prod, staging)"
	w.env.PromptStyle = focusedStyle

	// Description Input
	w.description = textinput.New()
	w.description.Placeholder = "Short Description"
	w.description.PromptStyle = focusedStyle

	w.focusIndex = 0

	return w
}

// Bubble Tea Init function
func (m workspaceModel) Init() tea.Cmd {
	return textinput.Blink
}

// Bubble Tea Update function
func (m workspaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.focusIndex == 2 { // On the last field
				name := m.workspaceName.Value()
				env := m.env.Value()
				description := m.description.Value()

				// Validate input
				if name == "" || env == "" || description == "" {
					fmt.Println("All fields are required.")
					return m, nil
				}

				if env != "prod" && env != "staging" {
					fmt.Println("Environment must be 'prod' or 'staging'.")
					return m, nil
				}

				// Execute command
				err := createWorkspace(name, env, description)
				if err != nil {
					fmt.Printf("Error creating workspace: %v\n", err)
				} else {
					fmt.Println("Workspace created successfully!")
				}
				return m, tea.Quit
			}
		case "tab":
			m.focusIndex = (m.focusIndex + 1) % 3
			m.updateFocus()
		}
	}

	// Update the focused input field
	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		m.workspaceName, cmd = m.workspaceName.Update(msg)
	case 1:
		m.env, cmd = m.env.Update(msg)
	case 2:
		m.description, cmd = m.description.Update(msg)
	}

	return m, cmd
} // Bubble Tea View function
func (m workspaceModel) View() string {
	if m.submitted {
		// Simulate triggering CLI command here
		workspaceName := m.workspaceName.Value()
		env := m.env.Value()
		description := m.description.Value()

		return fmt.Sprintf(
			"Workspace Created!\n\nName: %s\nEnvironment: %s\nDescription: %s\n",
			workspaceName, env, description,
		)
	}

	var b strings.Builder

	// Form Title
	b.WriteString("Create New Workspace\n")
	b.WriteString("Use Tab to navigate, Enter to submit, Ctrl+C to quit.\n\n")

	// Workspace Name
	b.WriteString(m.workspaceName.View())
	b.WriteString("\n\n")

	// Environment
	b.WriteString(m.env.View())
	b.WriteString("\n\n")

	// Description
	b.WriteString(m.description.View())

	return b.String()
}

// Helper to update focus on form fields
func (m *workspaceModel) updateFocus() {
	m.workspaceName.Blur()
	m.env.Blur()
	m.description.Blur()

	switch m.focusIndex {
	case 0:
		m.workspaceName.Focus()
	case 1:
		m.env.Focus()
	case 2:
		m.description.Focus()
	}
}

// StartWorkspaceView launches the form
func StartWorkspaceView() {
	p := tea.NewProgram(initialWorkspaceModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting workspace view: %v\n", err)
	}
}
