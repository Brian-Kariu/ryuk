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
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/config"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Removes a workspace",
	Long:  `Deletes the specified workspace. This is case sensitive.`,
	Run: func(cmd *cobra.Command, args []string) {
		var selected string
		var opt []huh.Option[string]
		for _, ws := range config.Workspaces {
			name := ws.Name
			id := ws.ID
			opt = append(opt, huh.NewOption(name, id))
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select workspace to delete").
					Value(&selected).
					Height(8).
					Options(opt...),
			),
		)
		err := form.Run()
		if err != nil {
			log.Fatal("Error with selection: %v", err)
		}
		config.DeleteWorkspace(selected)
		log.Info("Deleted workspace %v.", "selected", selected)
	},
}

func init() {
	WorkspaceCmd.AddCommand(deleteCmd)
	viper.BindPFlags(createCmd.Flags())
}
