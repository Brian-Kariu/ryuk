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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/config"
)

// variablesCmd represents the environment command
var VariablesCmd = &cobra.Command{
	Use:   "var",
	Short: "Manage the variables in your projects.",
	Long:  `Manage the various variables for your projects.`,
}

func init() {
	VariablesCmd.PersistentFlags().StringVarP(&config.CurrentWorkspace, "workspace", "w", "default", "Workspace currently in use.")
	viper.BindPFlag("workspace", VariablesCmd.PersistentFlags().Lookup("workspace"))

	VariablesCmd.MarkPersistentFlagRequired("workspace")
	VariablesCmd.PersistentFlags().StringVarP(&config.CurrentEnv, "env", "e", "", "Env currently in use.")
	viper.BindPFlag("env", VariablesCmd.PersistentFlags().Lookup("env"))
	VariablesCmd.MarkPersistentFlagRequired("env")
}
