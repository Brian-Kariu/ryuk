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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/cmd/config"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Args:  cobra.ExactArgs(1),
	Short: "Removes a workspace",
	Long:  `Deletes the specified workspace. This is case sensitive.`,
	Run: func(cmd *cobra.Command, args []string) {
		workspaceName := args[0]
		ws, err := config.GetWorkspace(workspaceName)
		if err != nil {
			fmt.Println("Error workspaces %s doesn't exist", workspaceName)
		}

		config.DeleteWorkspace(ws.ID)
		fmt.Printf("Deleted workspace %s.", workspaceName)
	},
}

func init() {
	WorkspaceCmd.AddCommand(deleteCmd)
	viper.BindPFlags(createCmd.Flags())
}
