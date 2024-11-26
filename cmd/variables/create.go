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
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/config"
	"github.com/Brian-Kariu/ryuk/db"
)

type Config struct {
	key   []byte
	value []byte
}

func createVar(bucket string, data db.Config) {
	client, err := db.NewClient(filepath.Join(config.BasePath, viper.GetString("workspace")), bucket)
	if err != nil {
		fmt.Printf("Error creating DB!")
	}
	err = client.AddKey(bucket, data)
	if err != nil {
		log.Fatalf("Failed to add key: %v", err)
	}
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
		var envValue string
		var confirm bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Input var name").
					Placeholder("Value").
					Value(&envName),
				huh.NewInput().
					Title("Input var value").
					Placeholder("Value").
					Value(&envValue),
				huh.NewConfirm().
					Title("Are you sure?").
					Affirmative("Yes!").
					Negative("No.").
					Value(&confirm),
			),
		)
		err := form.Run()
		if err != nil {
			log.Fatalf("Workspace name not set!")
		}
		if !viper.IsSet("workspace") {
			log.Fatal("Workspace flag not set!")
		}
		if !viper.IsSet("env") {
			log.Fatal("Env flag not set!")
		}
		log.Info("Env Name: ", "key", envName)
		log.Info("Env Value: ", "value", envValue)
		data := db.Config{Key: []byte(envName), Value: []byte(envValue)}
		createVar(viper.GetString("env"), data)
	},
}

func init() {
	VariablesCmd.AddCommand(createCmd)

	createCmd.MarkPersistentFlagRequired("workspace")
	createCmd.MarkPersistentFlagRequired("env")
}
