package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/config"
)

const useHighPerformanceRenderer = false

type (
	errMsg struct{ error }
)

var cmd tea.Cmd

type Entry struct {
	ready              bool
	viewport           viewport.Model
	error              string
	quitting           bool
	workspaceSelection *huh.Form
}

func NewEntry() Entry {
	err := viper.UnmarshalKey("workspaces", &config.Workspaces)
	if err != nil {
		log.Fatal("Error rendering workspaces: ", err)
		os.Exit(1)
	}
	if len(config.Workspaces) == 0 {
		log.Fatal("No workspaces found.", err)
		os.Exit(1)
	}
	items := []string{}
	for _, ws := range config.Workspaces {
		title := ws.Name
		items = append(items, title)
	}

	m := Entry{workspaceSelection: huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("workspace").
				Title("Choose a workspace").
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(items...)
				}, &items),
		),
	)}
	return m
}

func (m Entry) Init() tea.Cmd {
	return m.workspaceSelection.Init()
}

func (m Entry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if !m.ready {
		m.ready = true
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := DocStyle.GetMargin()
		m.viewport = viewport.New(WindowSize.Width-left-right, WindowSize.Height-top-bottom-6)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Create):
			// TODO: remove m.quitting after bug in Bubble Tea (#431) is fixed
			m.quitting = true
			return m, nil
		case key.Matches(msg, Keymap.Back):
			return InitWorkspace()
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}
	form, cmd := m.workspaceSelection.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.workspaceSelection = f
	}

	if m.workspaceSelection.State == huh.StateCompleted {
		return InitWorkspace()
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Entry) View() string {
	if m.quitting {
		return ""
	}
	if !m.ready {
		return "\n  Initializing..."
	}
	message := `
██████╗ ██╗   ██╗██╗   ██╗██╗  ██╗
██╔══██╗╚██╗ ██╔╝██║   ██║██║ ██╔╝
██████╔╝ ╚████╔╝ ██║   ██║█████╔╝ 
██╔══██╗  ╚██╔╝  ██║   ██║██╔═██╗ 
██║  ██║   ██║   ╚██████╔╝██║  ██╗
╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝

Welcome to Ryuk CLI!
`
	if m.workspaceSelection.State == huh.StateCompleted {
		workspace := m.workspaceSelection.GetString("workspace")
		return fmt.Sprintf("You selected: %s", workspace)
	}
	formatted := fmt.Sprintf("%s\n%s", message, m.workspaceSelection.View())
	return fmt.Sprintf(formatted)
}

func (m Entry) helpView() string {
	// TODO: use the keymaps to populate the help string
	return HelpStyle("\n ↑/↓: navigate  • esc: back • c: create entry • d: delete entry • q: quit\n")
}

func InitEntry() {
	m := NewEntry()
	top, right, bottom, left := DocStyle.GetMargin()
	m.viewport = viewport.New(WindowSize.Width-left-right, WindowSize.Height-top-bottom-1)
	m.viewport.Style = lipgloss.NewStyle().Align(lipgloss.Bottom)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal("Error rendering workspaces: ", err)
		os.Exit(1)
	}
}
