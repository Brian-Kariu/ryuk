/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize you go-ryuk app",
	Long:  `Initiailizes the app and sets configuration settings`,
	Run: func(cmd *cobra.Command, args []string) {
		environmentUpper := strings.ToUpper(environment)
		viper.SetEnvPrefix(environment)
		Create(strings.Join([]string{environmentUpper, "ENV"}, "_"), environment)
		viper.WatchConfig()
		env := viper.GetEnvPrefix()
		fmt.Printf("%s environment has been activated!", env)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	initCmd.PersistentFlags().StringVarP(&environment, "environment", "r", "test", "The environment the user is working in i.e test, staging")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
