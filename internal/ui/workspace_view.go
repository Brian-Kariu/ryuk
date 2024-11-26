package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Brian-Kariu/ryuk/config"
	"github.com/Brian-Kariu/ryuk/internal/utils"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	titleStyle          = lipgloss.NewStyle().MarginLeft(2)
	focusedButton       = focusedStyle.Render("[ Submit ]")
	blurredButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	itemStyle           = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle     = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	quitTextStyle       = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	mainStyle           = lipgloss.NewStyle().MarginLeft(2)
)

const (
	read   = "read"
	create = "create"
	remove = "delete"
)

type item string

const listHeight = 14

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int { return 1 }

func (d itemDelegate) Spacing() int { return 0 }

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// WorkspaceModel defines the state of the form
type workspaceModel struct {
	currentView     string
	inputs          []textinput.Model
	workspaceList   list.Model
	workspaceChoice string
	focusIndex      int
	submitted       bool
	cursorMode      cursor.Mode
	altscreenActive bool
	viewMode        bool
	Quitting        bool
	err             error
}

// Initialize the form inputs
func InitialWorkspaceModel(w []config.WorkspaceConfig) workspaceModel {
	spaces := []list.Item{}
	for _, ws := range w {
		spaces = append(spaces, item(ws.Name))
	}
	const defaultWidth = 20
	l := list.New(spaces, itemDelegate{}, defaultWidth, listHeight)

	workspace := workspaceModel{
		inputs:        make([]textinput.Model, 3),
		currentView:   read,
		workspaceList: l,
	}
	l.Title = "Workspaces List"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	var t textinput.Model
	for i := range workspace.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Workspace Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Environment (prod, staging, local)"
			t.PromptStyle = focusedStyle
		case 2:
			t.Placeholder = "Short Description"
			t.PromptStyle = focusedStyle
		}

		workspace.inputs[i] = t
	}

	return workspace
}

// Bubble Tea Init function
func (m workspaceModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m workspaceModel) createWorkspaceUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "a":
			m.altscreenActive = !m.altscreenActive
			cmd := tea.EnterAltScreen
			if !m.altscreenActive {
				cmd = tea.ExitAltScreen
			}
			return m, cmd
		case "ctrl+n":
			m.submitted = true
			var cmd tea.Cmd
			return m, cmd

		case "ctrl+s":
			m.currentView = read
			var cmd tea.Cmd
			return m, cmd

		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				args := []string{"env", "create", "-w", m.inputs[0].View()}
				result := utils.ExecuteCommand(args)
				m.submitted = true
				return m, result
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m workspaceModel) viewWorkspaceUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.workspaceList.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "a":
			m.altscreenActive = !m.altscreenActive
			cmd := tea.EnterAltScreen
			if !m.altscreenActive {
				cmd = tea.ExitAltScreen
			}
			return m, cmd
			// Change cursor mode
		case "ctrl+s":
			m.currentView = create
			var cmd tea.Cmd
			return m, cmd

		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)
		case "enter":
			i, ok := m.workspaceList.SelectedItem().(item)
			if ok {
				m.workspaceChoice = string(i)
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.workspaceList, cmd = m.workspaceList.Update(msg)
	return m, cmd
}

// Bubble Tea Update function
func (m workspaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.currentView {
		case create:
			return m.createWorkspaceUpdate(msg)
		case read:
			return m.viewWorkspaceUpdate(msg)
		}
	}
	return m, nil
}

func (m *workspaceModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Bubble Tea View function
func (m workspaceModel) createView() string {
	var b strings.Builder

	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}
	if m.submitted == true {
		b.WriteString(helpStyle.Render("Workspace successfully created!"))
		b.WriteString(helpStyle.Render("Press ctrl+s to view workspaces or a to add another workspace\n"))
		return b.String()
	}

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func (m workspaceModel) readView() string {
	if m.workspaceChoice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.workspaceChoice))
	}
	if m.Quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.workspaceList.View()
}

// View renders the UI based on the current view
func (m workspaceModel) View() string {
	var s string
	if m.Quitting {
		return "\n Exiting... \n"
	}
	if m.submitted == true {
		var b strings.Builder
		b.WriteString(helpStyle.Render("Workspace successfully created!"))
		b.WriteString(helpStyle.Render("Press ctrl+s to view workspaces or a to add another workspace"))
		s = b.String()
	}
	switch m.currentView {
	case create:
		s = m.createView()
	case read:
		s = m.readView()
	default:
		s = "Unknown view"
	}
	return mainStyle.Render("\n" + s + "\n\n")
}

func RenderWorkspace(w []config.WorkspaceConfig) {
	m := InitialWorkspaceModel(w)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
