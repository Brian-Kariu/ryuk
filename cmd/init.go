/*
Copyright © 2024 Brian-Kariu

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
	"fmt"
	"slices"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize you go-ryuk app",
	Long:  `Initiailizes the app and sets configuration settings`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("environment", environment)
		viper.Set("config", viper.ConfigFileUsed())
		viper.WriteConfigAs(viper.ConfigFileUsed())

		envFiles := []string{}
		validEnvFiles := 0
		scanFiles, err := Scan("", viper.GetString("environment"))
		if err != nil {
			return err
		}
		if len(scanFiles) != 0 && slices.Contains(scanFiles, viper.GetString("environment")) {
			validEnvs, err := Validate(scanFiles)
			if err != nil {
				return err
			}
			validEnvFiles = len(validEnvs)
			envFiles = scanFiles
		} else {
			envFileName := viper.GetString("environment") + ".env"
			envFiles = append(envFiles, envFileName)
		}

		environmentUpper := strings.ToUpper(viper.GetString("environment"))
		if validEnvFiles == 0 {
			if err := Create(strings.Join([]string{environmentUpper, "ENV"}, "_"), viper.GetString("environment"), envFiles[0]); err != nil {
				return err
			}
			if err := Create(SentinelKey, "TRUE", envFiles[0]); err != nil {
				return err
			}
		}
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
		})
		fmt.Printf("%s environment has been activated!", viper.GetString("environment"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
