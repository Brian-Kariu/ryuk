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
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/cmd/config"
	"github.com/Brian-Kariu/ryuk/db"
)

// TODO: This should be a standalone func that can be reusable
// FIX: This might also be okay since its only used here
func createDb(dbName, dbConfigs string) {
	db.NewClient(filepath.Join(config.BasePath, dbName), dbConfigs)
	config.NewWorkspaceConfig(dbName, []string{})
	viper.Set("workspaces", config.Workspaces)
	if err := viper.WriteConfig(); err != nil {
		fmt.Printf("Error saving workspace '%s': %v\n", dbName, err)
	}
	log.Printf("%s workspace has been created.\n", dbName)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Args:  cobra.MatchAll(cobra.MaximumNArgs(1)),
	Short: "Create a new workspace",
	Long: `Use this command to create a new resource. A resource
	could be a workspace
	`,
	Run: func(cmd *cobra.Command, args []string) {
		workspaceName := args[0]
		if workspaceName == "" {
			log.Fatal("Workspace name not set")
		}
		dbConfigs, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatal("DB Config name is not valid")
		}
		createDb(workspaceName, dbConfigs)
	},
}

func init() {
	WorkspaceCmd.AddCommand(createCmd)

	// myFlagSet := flags.NewCreateFlagSet(flags.Workspace)
	// createCmd.Flags().AddFlagSet(myFlagSet)

	// createCmd.PersistentFlags().String("workspace", "", "Current Workspace")
	// createCmd.PersistentFlags().String("config", "", "custom configs for workspace")
	//
	// viper.BindPFlags(createCmd.Flags())
	// viper.BindPFlag("current_workspace", createCmd.PersistentFlags().Lookup("workspace"))
}
