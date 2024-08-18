/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var envPair string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set an environment variable",
	Long:  `This command takes a single key/value pair and creates an environment variable.`,
	Run: func(cmd *cobra.Command, args []string) {
		envMap := strings.Split(envPair, " ")
		os.Setenv(strings.Join([]string{environment, envMap[0]}, "_"), envMap[1])
		fmt.Println("set called")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&envPair, "env", "e", "", "Add the key/value pair of the environment variable i.e 'Foo Something'")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
