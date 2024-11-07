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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/cmd/config"
	"github.com/Brian-Kariu/ryuk/cmd/flags"
	"github.com/Brian-Kariu/ryuk/db"
)

func createEnv(envName string) {
	client := db.NewClient(viper.GetString("workspace"), "")
	client.CreateBucket(envName)
	config.UpdateWorkspace(viper.GetString("workspace"), envName)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment",
	Long: `Use this command to create a new resource. A resource
	could be a workspace, environment or variable
	`,
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]
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

	// Bind flags to Viper
	viper.BindPFlags(createCmd.Flags())
}
