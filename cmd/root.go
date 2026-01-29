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
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     = true
	environment = "test"
)

type config struct {
	CfgFile     string `mapstructure:"config"`
	Environment string `mapstructure:"environment"`
}

var (
	configInstance config
	rootCmd        = &cobra.Command{
		Use:              "go-ryuk",
		Short:            "A simple environment variables management tool",
		Long:             `Ryuk (pronounced reeeuuk) aims to streamline the development experience of working with configs for your project. Ryuk handles creating,storing and sharing configs of your app.`,
		PersistentPreRun: preExecute,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func preExecute(cmd *cobra.Command, args []string) {
	err := viper.Unmarshal(&configInstance)
	if err != nil {
		log.Fatalf("Unable to read Viper options into configuration: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVarP(&cfgFile, "config", "c", true, "The config file for your project. Ensure you run this in the base directory of your project (default is $CWD/test.local)")
	viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", environment, "The environment the user is working in i.e test, staging")
	viper.BindPFlag("environment", rootCmd.Flags().Lookup("environment"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// TODO: Push this to its own utility file
func createFileIfNotExists(filename string) error {
	// Use os.OpenFile to open or create the file
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("File already exists")
			return nil
		}
		return err
	}
	defer file.Close()

	fmt.Println("File created successfully")
	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.Getwd()
	cobra.CheckErr(err)

	if cfgFile {
		// Use config file from the flag.
		configFile := filepath.Join(home, ".ryuk.yaml")
		createFileIfNotExists(configFile)

		viper.SetConfigFile(configFile)
	} else {
		fmt.Println("Ryuk not configured")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
