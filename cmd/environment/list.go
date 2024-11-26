/*
Copyright © 2024 Brian Kariu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package environment

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/config"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	docStyle          = lipgloss.NewStyle().Margin(1, 2)
)

type envitem struct {
	title, desc string
}

func (i envitem) Title() string { return i.title }

func (i envitem) Description() string { return i.desc }

func (i envitem) FilterValue() string { return i.title }

type envItemDelegate struct{}

func (d envItemDelegate) Height() int { return 1 }

func (d envItemDelegate) Spacing() int { return 0 }

func (d envItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d envItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(envitem)
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

type envModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m envModel) Init() tea.Cmd {
	return nil
}

func (m envModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(envitem)
			if ok {
				m.choice = i.title
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m envModel) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Viewing %s env.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Exiting envs view.")
	}
	return docStyle.Render(m.list.View())
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get details for the current environment.",
	Long:  `Read the configuration file and displays details for the current environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.UnmarshalKey("workspaces", &config.Workspaces)
		if err != nil {
			log.Fatal("Error fetching workspaces:", err)
			return
		}
		if len(config.Workspaces) == 0 {
			log.Fatal("No workspaces found.\n")
			return
		}

		currentWorkspace, err := config.GetWorkspace(viper.GetString("workspace"))
		if err != nil {
			log.Fatal("Error fetching current workspace:", err)
		}

		if len(currentWorkspace.Environment) == 0 {
			log.Printf("[]")
		}

		items := []list.Item{}
		var envs []string
		if len(currentWorkspace.Environment) != 0 {
			envs = slices.Collect(maps.Keys(currentWorkspace.Environment))
		}
		for _, env := range envs {
			title := env
			items = append(items, envitem{title: title, desc: ""})
		}

		l := list.New(items, list.NewDefaultDelegate(), 14, 20)
		l.Title = fmt.Sprintf("Listing envs in %s workspace", currentWorkspace.Name)
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle
		log.Info("Envs: ", "list", items)
		m := envModel{list: l}

		if _, err := tea.NewProgram(m).Run(); err != nil {
			log.Fatal("Error rendering envs: ", err)
		}
	},
}

func init() {
	EnvironmentCmd.AddCommand(listCmd)
}
