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
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/cmd/config"
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
	Args:  cobra.ExactArgs(2),
	Long: `Use this command to create a new resource. A resource
	could be a workspace, environment or variable
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "" && args[1] == "" {
			log.Fatal("Not specified key and value")
		}
		data := db.Config{Key: []byte(args[0]), Value: []byte(args[1])}
		createVar(viper.GetString("env"), data)
	},
}

func init() {
	VariablesCmd.AddCommand(createCmd)

	// myFlagSet := flags.NewCreateFlagSet(flags.Environment)
	// createCmd.Flags().AddFlagSet(myFlagSet)
	//
	// // Bind flags to Viper
	// viper.BindPFlags(createCmd.Flags())
}
