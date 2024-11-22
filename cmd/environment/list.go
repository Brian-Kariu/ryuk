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
	"maps"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/cmd/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get details for the current environment.",
	Long:  `Read the configuration file and displays details for the current environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.UnmarshalKey("workspaces", &config.Workspaces)
		if err != nil {
			fmt.Println("Error fetching workspaces:", err)
			return
		}
		if len(config.Workspaces) == 0 {
			fmt.Printf("No workspaces found.\n")
			return
		}

		currentWorkspace, err := config.GetWorkspace(viper.GetString("workspace"))
		if err != nil {
			fmt.Println("Error fetching current workspace:", err)
		}

		fmt.Printf("Activated Workspace: %s\n", config.CurrentWorkspace)
		fmt.Printf("Activated Environment: %s\n", config.CurrentEnv)
		fmt.Printf("Available Environments\n")
		if len(currentWorkspace.Environment) == 0 {
			fmt.Print("[]")
		}
		if len(currentWorkspace.Environment) != 0 {
			envs := slices.Collect(maps.Keys(currentWorkspace.Environment))
			fmt.Printf("Envs: %s", envs)

		}
	},
}

func init() {
	EnvironmentCmd.AddCommand(listCmd)
}
