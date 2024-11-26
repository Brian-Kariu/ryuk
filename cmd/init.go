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
package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/Brian-Kariu/ryuk/config"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ryuk",
	Long:  `Initializes ryuk application in your system`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
		defaultDbName := "default"
		envs := []string{"prod"}

		configFileInstance := newConfigFile(config.BasePath, ".ryuk.yaml")
		file, _ := os.Stat(configFileInstance.fullpath)
		if file != nil {
			log.Info("Ryuk app already initialized")
			return
		}
		configFileInstance.checkDir()
		configFileInstance.checkFile()

		config.NewWorkspaceConfig(defaultDbName, "Default ryuk workspace", envs, false)
		initGlobalDb(config.BasePath)
		log.Info("Ryuk app initalized!")
		return
	},
}

func init() {
}
