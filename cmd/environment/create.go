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
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/cmd/flags"
	"github.com/Brian-Kariu/ryuk/config"
	"github.com/Brian-Kariu/ryuk/db"
)

func createEnv(envName string) {
	client, err := db.NewClient(filepath.Join(config.BasePath, viper.GetString("workspace")), viper.GetString("env"))
	if err != nil {
		log.Fatal("Error creating DB!")
	}
	client.CreateBucket(envName)
	config.UpdateWorkspace(viper.GetString("workspace"), envName)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment",
	Args:  cobra.NoArgs,
	Long: `Use this command to create a new resource. A resource
	could be a workspace, environment or variable
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var envName string
		input := huh.NewInput().
			Title("Input env Name").
			Prompt("?").
			Value(&envName)
		err := input.Run()
		if err != nil {
			log.Fatalf("Environment name not set!")
		}
		if !viper.IsSet("workspace") {
			log.Fatal("Workspace config not set.")
		}

		if envName == "" {
			log.Fatal("Environment name not set")
		}
		createEnv(envName)
	},
}

func init() {
	EnvironmentCmd.AddCommand(createCmd)

	myFlagSet := flags.NewCreateFlagSet(flags.Environment)
	createCmd.Flags().AddFlagSet(myFlagSet)

	createCmd.MarkPersistentFlagRequired("workspace")

	// Bind flags to Viper
	viper.BindPFlags(createCmd.Flags())
}
