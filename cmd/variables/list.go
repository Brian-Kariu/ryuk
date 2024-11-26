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
package variables

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"

	"github.com/Brian-Kariu/ryuk/config"
	"github.com/Brian-Kariu/ryuk/db"
)

type Var struct {
	key string
	val string
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environment variables",
	Long: `Use this command to create a new resource. A resource
	could be a workspace, environment or variable
	`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := db.NewClient(filepath.Join(config.BasePath, viper.GetString("workspace")), viper.GetString("env"))
		if err != nil {
			log.Error("Error creating DB!")
		}
		envVars, err := client.ListVars(viper.GetString("env"))
		if err != nil {
			log.Fatal(err)
		}
		columns := []table.Column{
			{Title: "Key", Width: 40},
			{Title: "Value", Width: 40},
		}
		vars := []Var{}
		for _, ws := range maps.Keys(envVars) {
			key := ws
			value := envVars[ws]
			vars = append(vars, Var{key: key, val: value})
		}
		individualRows := make([]table.Row, len(vars)) // Preallocate rows slice
		for i, env := range vars {
			individualRows[i] = table.Row{env.key, env.val}
		}
		t := table.New(
			table.WithColumns(columns),
			table.WithRows(individualRows),
			table.WithFocused(true),
			table.WithHeight(7),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		m := model{t}
		if _, err := tea.NewProgram(m).Run(); err != nil {
			log.Error("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	VariablesCmd.AddCommand(listCmd)
}
