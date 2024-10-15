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
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/cmd/flags"
	"github.com/Brian-Kariu/ryuk/db"
)

func createDb(dbName, dbConfigs string) {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	basePath := filepath.Join(home, ".ryuk/")
	db.NewClient(filepath.Join(basePath, dbName), dbConfigs)
	workspace := map[string]string{
		"name":        dbName,
		"path":        filepath.Join(basePath, dbName+".db"),
		"environment": "prod",
	}
	viper.Set("workspaces", workspace)
	viper.SafeWriteConfig()

	log.Printf("%s workspace has been created.\n", dbName)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new workspace",
	Long: `Use this command to create a new resource. A resource
	could be a workspace
	`,
	Run: func(cmd *cobra.Command, args []string) {
		dbName, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal("DB name is not valid")
		}
		dbConfigs, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatal("DB Config name is not valid")
		}
		createDb(dbName, dbConfigs)
	},
}

func init() {
	WorkspaceCmd.AddCommand(createCmd)

	myFlagSet := flags.NewCreateFlagSet(flags.Workspace)
	createCmd.Flags().AddFlagSet(myFlagSet)

	createCmd.PersistentFlags().String("config", "", "custom configs for workspace")

	viper.BindPFlags(createCmd.Flags())
}
