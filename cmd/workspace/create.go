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
package workspace

import (
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/Brian-Kariu/ryuk/config"
	"github.com/Brian-Kariu/ryuk/db"
)

// TODO: This should be a standalone func that can be reusable
// FIX: This might also be okay since its only used here
func createDb(dbName, description, dbConfigs string, confirm bool) {
	db.NewClient(filepath.Join(config.BasePath, dbName), dbConfigs)
	config.NewWorkspaceConfig(dbName, description, []string{}, confirm)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Args:  cobra.NoArgs,
	Short: "Create a new workspace",
	Long: `Use this command to create a new resource. A resource
	could be a workspace
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var workspaceName string
		var description string
		var confirm bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Input workspace Name").
					Placeholder("workspace").
					Value(&workspaceName),
				huh.NewInput().
					Title("Input short description").
					Placeholder("description").
					Value(&description),
				huh.NewConfirm().
					Title("Add current path as project path? (You can also input the path manually)").
					Affirmative("Yes!").
					Negative("No.").
					Value(&confirm),
			),
		)
		err := form.Run()
		dbConfigs, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatal("DB Config name is not valid")
		}
		createDb(workspaceName, description, dbConfigs, confirm)
	},
}

func init() {
	WorkspaceCmd.AddCommand(createCmd)
}
