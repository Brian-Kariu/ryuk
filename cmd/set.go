/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var envPair string

func appendToFile(filename string, data string) error {
	// Open the file with os.O_APPEND and os.O_WRONLY
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data to the file
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	fmt.Println("Data appended successfully")
	return nil
}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set an environment variable",
	Long:  `This command takes a single key/value pair and creates an environment variable.`,
	Run: func(cmd *cobra.Command, args []string) {
		envVars := strings.Join(args, "=")
		fileName := viper.GetString("env")
		appendToFile(fileName, envVars)
		fmt.Printf("Env variable has been set for file, %s", fileName)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// setCmd.Flags().StringVarP(&envPair, "env", "e", "", "Add the key/value pair of the environment variable i.e 'Foo Something'")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
