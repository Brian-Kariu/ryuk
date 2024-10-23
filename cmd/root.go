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
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Brian-Kariu/ryuk/cmd/config"
	"github.com/Brian-Kariu/ryuk/cmd/environment"
	"github.com/Brian-Kariu/ryuk/cmd/workspace"
	"github.com/Brian-Kariu/ryuk/db"
)

var cfgFile string

type configFile struct {
	path     string
	fileName string
	fullpath string
}

func (c configFile) checkDir() {
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		err := os.MkdirAll(c.path, 0755)
		cobra.CheckErr(err)
	}
}

func (c configFile) checkFile() {
	if _, err := os.Stat(c.fullpath); os.IsNotExist(err) {
		file, err := os.Create(c.fullpath)
		cobra.CheckErr(err)
		defer file.Close()
	}
}

func newConfigFile(path, fileName string) *configFile {
	fullpath := filepath.Join(path, fileName)
	return &configFile{
		path:     path,
		fileName: fileName,
		fullpath: fullpath,
	}
}

var RootCmd = &cobra.Command{
	Use:   "ryuk",
	Short: "A fast configuration management library",
	Long:  `Ryuk is a powerful cli app that helps you manage your application configs and secrets!`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommands() {
	RootCmd.AddCommand(workspace.WorkspaceCmd, environment.EnvironmentCmd)
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ryuk/ryuk.yaml)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addSubcommands()
}

func initGlobalDb(path string) {
	dbInstance := db.NewClient(filepath.Join(path, "default"), "")
	dbInstance.CreateBucket("prod")
}

func initConfig() {
	// NOTE: cfgFile and configFile need to be aligned, could cause issues down the line
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		config.InitConstants()
		configFileInstance := newConfigFile(config.BasePath, ".ryuk.yaml")
		configFileInstance.checkDir()
		configFileInstance.checkFile()

		viper.AddConfigPath(config.BasePath)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ryuk")

	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if viper.IsSet("workspaces") == false {
		defaultDbName := "default"
		envs := []string{"prod"}
		config.NewWorkspaceConfig(defaultDbName, envs)
		initGlobalDb(config.BasePath)

		viper.WriteConfig()

	}

	err := viper.UnmarshalKey("workspaces", &config.Workspaces)
	if err != nil {
		fmt.Println("Error initializing workspaces:", err)
	}
}
